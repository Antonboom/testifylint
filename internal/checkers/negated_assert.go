package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

// NegatedAssertReport is the diagnostic message emitted by NegatedAssert.
const NegatedAssertReport = "use require instead"

// NegatedAssert detects patterns where a negated assertion in an if condition guards
// a return or continue statement, and provides a fix that replaces the whole if block
// with the corresponding require call.
//
// For example:
//
//	if !assert.NoError(t, err) {
//	    return
//	}
//
// is replaced with:
//
//	require.NoError(t, err)
//
// This applies to any assertion, not only error-related ones:
//
//	if !assert.Equal(t, expected, actual) {
//	    return
//	}
//
// becomes:
//
//	require.Equal(t, expected, actual)
//
// Compound || conditions are also handled:
//
//	if !assert.NoError(t, err) || !assert.Equal(t, expected, actual) {
//	    return
//	}
//
// becomes:
//
//	require.NoError(t, err)
//	require.Equal(t, expected, actual)
//
// NegatedAssert skips patterns in goroutines, HTTP handlers, and test cleanup
// functions, as well as else-if chains and if blocks with complex bodies.
type NegatedAssert struct{}

// NewNegatedAssert constructs NegatedAssert checker.
func NewNegatedAssert() NegatedAssert { return NegatedAssert{} }
func (NegatedAssert) Name() string    { return "negated-assert" }

func (checker NegatedAssert) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	insp.WithStack([]ast.Node{(*ast.IfStmt)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		ifStmt := node.(*ast.IfStmt)

		// Pattern requires: single return/continue body, no else clause, no init.
		// Skip if there is an init statement (e.g. `if err := foo(); !assert.NoError(...)`):
		// replacing the whole if-statement would drop the init and break compilation.
		if ifStmt.Init != nil || !isNegatedAssertSimpleBody(ifStmt) || ifStmt.Else != nil {
			return true
		}

		// Skip else-if branches: fixing them would leave a dangling "else" in the caller.
		if len(stack) >= 2 {
			if _, parentIsIf := stack[len(stack)-2].(*ast.IfStmt); parentIsIf {
				return true
			}
		}

		// Skip goroutines, HTTP handlers, and test cleanup functions.
		fID := findSurroundingFunc(pass, stack)
		if fID == nil {
			return true
		}
		if fID.meta.isGoroutine || fID.meta.isHTTPHandler || fID.meta.isTestCleanup {
			return true
		}

		// Collect all top-level negated calls from the condition.
		condCalls, allNegatedOr := collectNegatedAssertOrCalls(ifStmt.Cond)
		if !allNegatedOr || len(condCalls) == 0 {
			return true
		}

		// All calls must be testify assert package calls (not object/suite calls).
		assertCalls := make([]*CallMeta, 0, len(condCalls))
		for _, ce := range condCalls {
			cm := NewCallMeta(pass, ce)
			if cm == nil || !cm.IsAssert || !cm.IsPkg() {
				return true
			}
			assertCalls = append(assertCalls, cm)
		}

		// Sort by source position for stable, deterministic output.
		sort.Slice(assertCalls, func(i, j int) bool {
			return assertCalls[i].Call.Pos() < assertCalls[j].Call.Pos()
		})

		firstCall := assertCalls[0]
		if fix, ok := buildNegatedAssertFix(pass, ifStmt, assertCalls); ok {
			diagnostics = append(diagnostics,
				*newDiagnostic(checker.Name(), firstCall, NegatedAssertReport, fix))
		} else {
			diagnostics = append(diagnostics,
				*newDiagnostic(checker.Name(), firstCall, NegatedAssertReport))
		}

		return true
	})

	return diagnostics
}

// isNegatedAssertSimpleBody reports whether the body of ifStmt is exactly one
// return or continue statement. Other statement types (e.g. assignments, expression
// statements) are intentionally excluded: they imply side effects that make the
// if-block semantically non-trivial to replace with a bare require call.
func isNegatedAssertSimpleBody(ifStmt *ast.IfStmt) bool {
	if len(ifStmt.Body.List) != 1 {
		return false
	}
	switch s := ifStmt.Body.List[0].(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return s.Tok == token.CONTINUE
	}
	return false
}

// collectNegatedAssertOrCalls recursively collects all top-level CallExpr nodes from
// a condition that consists exclusively of negated calls (UnaryExpr with !) joined
// by logical-or (||). ParenExpr nodes are treated as transparent. Returns the calls
// in left-to-right order and true if the entire expression matches that pattern;
// returns nil and false otherwise.
func collectNegatedAssertOrCalls(expr ast.Expr) ([]*ast.CallExpr, bool) {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return collectNegatedAssertOrCalls(e.X)
	case *ast.UnaryExpr:
		if e.Op == token.NOT {
			x := unwrapParen(e.X)
			if ce, ok := x.(*ast.CallExpr); ok {
				return []*ast.CallExpr{ce}, true
			}
		}
		return nil, false
	case *ast.BinaryExpr:
		if e.Op != token.LOR {
			return nil, false
		}
		left, leftOk := collectNegatedAssertOrCalls(e.X)
		right, rightOk := collectNegatedAssertOrCalls(e.Y)
		if !leftOk || !rightOk {
			return nil, false
		}
		return append(left, right...), true
	default:
		return nil, false
	}
}

// unwrapParen peels off any number of parentheses and returns the innermost expression.
func unwrapParen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

// buildNegatedAssertFix constructs a SuggestedFix that replaces the entire ifStmt
// with one require.XXX call per assertion in calls (in source order).
// The require import is added to the file if it is not already present.
func buildNegatedAssertFix(pass *analysis.Pass, ifStmt *ast.IfStmt, calls []*CallMeta) (analysis.SuggestedFix, bool) {
	if len(calls) == 0 {
		return analysis.SuggestedFix{}, false
	}

	qualName, importEdit, ok := addImportFix(pass.Files, calls[0].Call.Pos(), testify.RequirePkgPath)
	if !ok {
		return analysis.SuggestedFix{}, false
	}

	indent := negatedAssertLineIndent(pass, ifStmt.Pos())
	requireCalls := make([]string, 0, len(calls))
	for _, c := range calls {
		callText := analysisutil.NodeString(pass.Fset, c.Call)
		var newCallText string
		if qualName == "" {
			// dot-import: call the function directly without a qualifier;
			// strip the old "assert." prefix entirely.
			newCallText = callText[len(c.SelectorXStr)+1:]
		} else {
			newCallText = qualName + callText[len(c.SelectorXStr):]
		}
		requireCalls = append(requireCalls, newCallText)
	}
	newText := strings.Join(requireCalls, "\n"+indent)

	textEdits := []analysis.TextEdit{
		{
			Pos:     ifStmt.Pos(),
			End:     ifStmt.End(),
			NewText: []byte(newText),
		},
	}
	if importEdit != nil {
		textEdits = append(textEdits, *importEdit)
	}

	var msg string
	if qualName == "" {
		msg = "Replace with " + calls[0].Fn.Name
		if len(calls) > 1 {
			msg = "Replace with require calls"
		}
	} else {
		msg = fmt.Sprintf("Replace with %s.%s", qualName, calls[0].Fn.Name)
		if len(calls) > 1 {
			msg = "Replace with " + qualName + " calls"
		}
	}

	return analysis.SuggestedFix{
		Message:   msg,
		TextEdits: textEdits,
	}, true
}

// negatedAssertLineIndent returns the leading whitespace of the source line containing pos.
func negatedAssertLineIndent(pass *analysis.Pass, pos token.Pos) string {
	tokenFile := pass.Fset.File(pos)
	if tokenFile == nil {
		return "\t"
	}

	content, err := pass.ReadFile(tokenFile.Name())
	if err != nil {
		return "\t"
	}

	offset := tokenFile.Offset(pos)

	lineStart := offset
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	end := lineStart
	for end < offset && (content[end] == ' ' || content[end] == '\t') {
		end++
	}

	return string(content[lineStart:end])
}
