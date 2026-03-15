package checkers

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// FailNow detects situations like
//
// assert.Fail(t, "msg")
// assert.Fail(t, "msg", args...)
// assert.Failf(t, "failure", "format %s", arg)
// assert.FailNow(t, "msg")
// assert.FailNow(t, "msg", args...)
// assert.FailNowf(t, "failure", "format %s", arg)
//
// and requires
//
//	t.Error("msg") / t.Errorf("format %s", arg)
//	t.Fatal("msg") / t.Fatalf("format %s", arg)
type FailNow struct{}

// NewFailNow constructs FailNow checker.
func NewFailNow() FailNow    { return FailNow{} }
func (FailNow) Name() string { return "fail-now" }

func (checker FailNow) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		ce := node.(*ast.CallExpr)
		call := NewCallMeta(pass, ce)
		if call == nil {
			return true
		}

		var msg string
		var proposedFn string
		switch call.Fn.NameFTrimmed {
		case "Fail":
			msg = "use t.Error or t.Errorf instead"
			proposedFn = "Error"
		case "FailNow":
			msg = "use t.Fatal or t.Fatalf instead"
			proposedFn = "Fatal"
		default:
			return true
		}

		// Skip non-pkg non-suite method calls (e.g., assertObj.Fail(...)) where there is
		// no accessible testing.T variable to suggest in the message.
		if !call.IsPkg && !implementsTestifySuite(pass, call.Selector.X) {
			return true
		}

		// Only provide the autofix when the call result is definitely discarded,
		// to avoid producing uncompilable code when the bool result is used in an expression.
		var fixes []analysis.SuggestedFix
		if isCallResultDiscarded(stack) {
			fixes = checker.fix(pass, call, proposedFn)
		}

		diagnostics = append(diagnostics, *newDiagnostic(checker.Name(), call, msg, fixes...))
		return true
	})
	return diagnostics
}

// isCallResultDiscarded reports whether the call result is discarded
// (i.e., the call is an expression statement, go statement, or defer statement).
func isCallResultDiscarded(stack []ast.Node) bool {
	if len(stack) < 2 {
		return false
	}
	switch stack[len(stack)-2].(type) {
	case *ast.ExprStmt, *ast.GoStmt, *ast.DeferStmt:
		return true
	}
	return false
}

// fix builds a SuggestedFix replacing the testify call with a standard testing.T method call.
//
// Argument mapping:
//
// Fmt variants (Failf/FailNowf): drop failureMessage, keep format + args
//
//	assert.Failf(t, "failure", "fmt %s", arg) → t.Errorf("fmt %s", arg)
//
// Non-fmt, 1 arg (failureMessage only):
//
//	assert.Fail(t, "msg") → t.Error("msg")
//
// Non-fmt, 2 args (failureMessage + one msgAndArgs element):
//
//	assert.Fail(t, "failure", "msg") → t.Error("msg")
//
// Non-fmt, 3+ args (failureMessage + format + args):
//
//	assert.Fail(t, "failure", "fmt %s", arg) → t.Errorf("fmt %s", arg)
func (checker FailNow) fix(pass *analysis.Pass, call *CallMeta, proposedFn string) []analysis.SuggestedFix {
	callerVar, callerExpr, ok := checker.callerVar(pass, call)
	if !ok {
		return nil
	}

	var newArgs []ast.Expr
	fn := proposedFn

	if call.Fn.IsFmt {
		// Failf(t, failureMessage, format, args...) → callerVar.Errorf(format, args...)
		if len(call.Args) < 2 {
			return nil
		}
		fn += "f"
		newArgs = call.Args[1:]
	} else {
		switch len(call.Args) {
		case 0:
			return nil
		case 1:
			// Fail(t, failureMessage) → callerVar.Error(failureMessage)
			newArgs = call.Args
		case 2:
			// Fail(t, failureMessage, msg) → callerVar.Error(msg)
			newArgs = call.Args[1:]
		default:
			// Fail(t, failureMessage, format, args...) → callerVar.Errorf(format, args...)
			fn += "f"
			newArgs = call.Args[1:]
		}
	}

	// Guard: only emit the fix if the target method exists on the caller's type.
	if !typeHasMethod(pass, callerExpr, fn) {
		return nil
	}

	newText := []byte(fmt.Sprintf("%s.%s(%s)", callerVar, fn, formatAsCallArgs(pass, newArgs...)))
	return []analysis.SuggestedFix{{
		Message: fmt.Sprintf("Replace `%s` with `%s.%s`", call.Fn.Name, callerVar, fn),
		TextEdits: []analysis.TextEdit{{
			Pos:     call.Pos(),
			End:     call.End(),
			NewText: newText,
		}},
	}}
}

// callerVar returns the string form and AST expression for the testing.T variable:
//   - for package-level calls (assert.Fail(t, ...)): returns the t argument and its expr.
//   - for suite method calls (s.Fail(...)): returns "s.T()" with nil expr (s.T() is *testing.T).
func (checker FailNow) callerVar(pass *analysis.Pass, call *CallMeta) (callerStr string, callerExpr ast.Expr, ok bool) {
	if call.IsPkg {
		if len(call.ArgsRaw) == 0 {
			return "", nil, false
		}
		arg := call.ArgsRaw[0]
		return analysisutil.NodeString(pass.Fset, arg), arg, true
	}

	// Suite method call: use receiver.T() to access *testing.T.
	if implementsTestifySuite(pass, call.Selector.X) {
		return analysisutil.NodeString(pass.Fset, call.Selector.X) + ".T()", nil, true
	}
	return "", nil, false
}

// typeHasMethod reports whether the type of expr has an accessible method named methodName.
// If expr is nil (used for synthetic callers like "s.T()"), returns true assuming *testing.T.
func typeHasMethod(pass *analysis.Pass, expr ast.Expr, methodName string) bool {
	if expr == nil {
		// Caller is a derived expression (e.g., "s.T()"). *testing.T has all needed methods.
		return true
	}
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	obj, _, _ := types.LookupFieldOrMethod(t, true, nil, methodName)
	return obj != nil
}
