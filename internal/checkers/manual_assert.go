package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// ManualAssert detects hand-rolled `if cond { t.Fatal/Errorf(...) }` patterns
// and suggests the equivalent testify assertion.
//
//	if err != nil { t.Fatalf("foo: %v", err) }   -> require.NoError(t, err)
//	if got != want { t.Errorf("got %v want %v", got, want) } -> assert.Equal(t, want, got)
//	if s == "" { t.Fatal("empty") }              -> require.NotEmpty(t, s)
//
// `t.Fatal*` is rewritten to `require.*`, `t.Error*` is rewritten to `assert.*`.
// By default the original failure message is forwarded as `msgAndArgs`; pass
// `manual-assert.drop-message=true` to drop it instead.
type ManualAssert struct {
	dropMessage bool
}

// NewManualAssert constructs ManualAssert checker.
func NewManualAssert() *ManualAssert {
	return &ManualAssert{}
}

func (*ManualAssert) Name() string { return "manual-assert" }

// SetDropMessage controls whether the original `t.Fatal/Errorf` message is
// dropped (default false — the message is preserved as `msgAndArgs`).
func (checker *ManualAssert) SetDropMessage(v bool) *ManualAssert {
	checker.dropMessage = v
	return checker
}

// Check if the if statement with a single t.Error or t.Fatal can be converted to testify.
func (checker *ManualAssert) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	insp.Preorder([]ast.Node{(*ast.IfStmt)(nil)}, func(node ast.Node) {
		ifStmt := node.(*ast.IfStmt)
		if ifStmt.Else != nil {
			return
		}
		if len(ifStmt.Body.List) != 1 {
			return
		}

		exprStmt, ok := ifStmt.Body.List[0].(*ast.ExprStmt)
		if !ok {
			return
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return
		}

		tInfo, ok := matchTestingFailCall(pass, call)
		if !ok {
			return
		}

		fn, args := checker.matchPattern(pass, ifStmt.Cond)
		if fn == "" {
			return
		}

		// Forward the original `t.Fatal/Errorf` message as `msgAndArgs` unless dropped.
		if !checker.dropMessage && len(tInfo.msgArgs) > 0 {
			args = append(args, tInfo.msgArgs...)
		}

		pkg := "assert"
		if tInfo.useFatal {
			pkg = "require"
		}

		// Inline `if x := f(); cond { ... }` would require lifting `x := f()`
		// out of the if header — skip for now to avoid scope-leak fixes.
		if ifStmt.Init != nil {
			diagnostics = append(diagnostics, *newDiagnostic(checker.Name(), ifStmt,
				fmt.Sprintf("replace with %s.%s", pkg, fn)))
			return
		}

		recv := analysisutil.NodeString(pass.Fset, tInfo.recv)
		newCall := buildCall(pass, pkg, fn, recv, args)

		diagnostics = append(diagnostics, *newDiagnostic(checker.Name(), ifStmt,
			fmt.Sprintf("replace with %s.%s", pkg, fn),
			analysis.SuggestedFix{
				Message: fmt.Sprintf("Replace with %s.%s", pkg, fn),
				TextEdits: []analysis.TextEdit{{
					Pos:     ifStmt.Pos(),
					End:     ifStmt.End(),
					NewText: []byte(newCall),
				}},
			},
		))
	})
	return diagnostics
}

type testFailCall struct {
	recv     ast.Expr   // the `t` in `t.Fatal(...)`
	useFatal bool       // Fatal/Fatalf -> require; Error/Errorf -> assert
	msgArgs  []ast.Expr // arguments to t.Fatal/Errorf
}

// matchTestingFailCall returns metadata when `call` is `t.Fatal/Fatalf/Error/Errorf`
// on something implementing testing.T.
func matchTestingFailCall(pass *analysis.Pass, call *ast.CallExpr) (testFailCall, bool) {
	se, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return testFailCall{}, false
	}
	if !implementsTestingT(pass, se.X) {
		return testFailCall{}, false
	}
	switch se.Sel.Name {
	case "Fatal", "Fatalf":
		return testFailCall{recv: se.X, useFatal: true, msgArgs: call.Args}, true
	case "Error", "Errorf":
		return testFailCall{recv: se.X, msgArgs: call.Args}, true
	}
	return testFailCall{}, false
}

// matchPattern inspects the if condition and returns (pkg, fn, leadingArgs) for
// the equivalent 'testify' assertion. `leadingArgs` does NOT include the leading
// receiver (`t`); buildCall prepends it.
func (checker *ManualAssert) matchPattern(pass *analysis.Pass, cond ast.Expr) (string, []ast.Expr) {
	switch c := cond.(type) {
	case *ast.BinaryExpr:
		return matchBinaryPattern(pass, c)

	case *ast.UnaryExpr:
		// `if !X { fail }` means we expect X to be true.
		if c.Op == token.NOT {
			if fn, args := matchTruthyPattern(pass, c.X, true); fn != "" {
				return fn, args
			}
		}

	case *ast.CallExpr:
		// `if errors.Is(err, target) { fail }` means we expect !errors.Is.
		if isErrorsIsCall(pass, c) && len(c.Args) == 2 {
			return "NotErrorIs", []ast.Expr{c.Args[0], c.Args[1]}
		}
		if isPkgFnCall(pass, c, "strings", "Contains") && len(c.Args) == 2 {
			return "NotContains", []ast.Expr{c.Args[0], c.Args[1]}
		}

	case *ast.Ident:
		// `if cond { fail }` means cond should have been false.
		if hasBoolType(pass, c) {
			return "False", []ast.Expr{c}
		}

	case *ast.SelectorExpr:
		if hasBoolType(pass, c) {
			return "False", []ast.Expr{c}
		}
	}
	return "", nil
}

// matchTruthyPattern handles `!X` where the negation flips the expected fix.
// `negated=true` means we are inside an outer `!`.
func matchTruthyPattern(pass *analysis.Pass, inner ast.Expr, negated bool) (string, []ast.Expr) {
	switch c := inner.(type) {
	case *ast.CallExpr:
		if isErrorsIsCall(pass, c) && len(c.Args) == 2 {
			if negated {
				return "ErrorIs", []ast.Expr{c.Args[0], c.Args[1]}
			}
			return "NotErrorIs", []ast.Expr{c.Args[0], c.Args[1]}
		}
		if isPkgFnCall(pass, c, "strings", "Contains") && len(c.Args) == 2 {
			if negated {
				return "Contains", []ast.Expr{c.Args[0], c.Args[1]}
			}
			return "NotContains", []ast.Expr{c.Args[0], c.Args[1]}
		}
	case *ast.Ident, *ast.SelectorExpr:
		if hasBoolType(pass, inner) {
			if negated {
				return "True", []ast.Expr{inner}
			}
			return "False", []ast.Expr{inner}
		}
	}
	return "", nil
}

func matchBinaryPattern(pass *analysis.Pass, be *ast.BinaryExpr) (string, []ast.Expr) {
	if be.Op == token.NEQ || be.Op == token.EQL {
		return matchEqualityPattern(pass, be)
	}
	return "", nil
}

func matchEqualityPattern(pass *analysis.Pass, be *ast.BinaryExpr) (string, []ast.Expr) {
	eq := be.Op == token.EQL // `if a == b { fail }` => fail when equal => assert NotEqual
	x, y := be.X, be.Y

	// Normalize so that literal/nil/"" sits on the right when one side is a value.
	if isNil(x) || isEmptyStringLit(x) || isZero(x) {
		x, y = y, x
	}

	// x == nil / x != nil
	if isNil(y) {
		if isError(pass, x) {
			if eq {
				return "Error", []ast.Expr{x}
			}
			return "NoError", []ast.Expr{x}
		}
		if eq {
			return "NotNil", []ast.Expr{x}
		}
		return "Nil", []ast.Expr{x}
	}

	// s == "" / s != ""
	if isEmptyStringLit(y) {
		if eq {
			return "NotEmpty", []ast.Expr{x}
		}
		return "Empty", []ast.Expr{x}
	}

	// len(x) == 0 / != 0 / != N
	if inner, ok := isBuiltinLenCall(pass, x); ok {
		if isZero(y) {
			if eq {
				return "NotEmpty", []ast.Expr{inner}
			}
			return "Empty", []ast.Expr{inner}
		}
		if !eq { // len(x) != N
			return "Len", []ast.Expr{inner, y}
		}
	}

	// Generic equality.
	// `if a != b { fail }` => assert.Equal(t, b, a) — testify wants (expected, actual);
	// we cannot reliably know which is which, but convention is `got != want`, so we
	// pass y as the "expected" (matches the most common t.Errorf("got %v want %v", a, b) pattern).
	if eq {
		return "NotEqual", []ast.Expr{y, x}
	}
	return "Equal", []ast.Expr{y, x}
}

func buildCall(pass *analysis.Pass, pkg, fn, recv string, args []ast.Expr) string {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, recv)
	for _, a := range args {
		parts = append(parts, analysisutil.NodeString(pass.Fset, a))
	}
	return pkg + "." + fn + "(" + strings.Join(parts, ", ") + ")"
}
