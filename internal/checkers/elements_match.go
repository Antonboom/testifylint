package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// ElementsMatch detects situations like
//
//	slices.Sort(a)
//	slices.Sort(b)
//	assert.True(t, slices.Equal(a, b))
//
// and requires
//
//	assert.ElementsMatch(t, a, b)
type ElementsMatch struct{}

// NewElementsMatch constructs ElementsMatch checker.
func NewElementsMatch() ElementsMatch { return ElementsMatch{} }
func (ElementsMatch) Name() string    { return "elements-match" }

func (checker ElementsMatch) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	insp.Preorder([]ast.Node{(*ast.BlockStmt)(nil)}, func(node ast.Node) {
		block := node.(*ast.BlockStmt)
		stmts := block.List

		for i := 0; i+2 < len(stmts); i++ {
			d := checker.checkPattern(pass, stmts[i], stmts[i+1], stmts[i+2])
			if d != nil {
				diagnostics = append(diagnostics, *d)
			}
		}
	})

	return diagnostics
}

// checkPattern checks if three consecutive statements form the pattern:
// slices.Sort(x), slices.Sort(y), assert.True(t, slices.Equal(x, y)).
func (checker ElementsMatch) checkPattern(
	pass *analysis.Pass,
	s0, s1, s2 ast.Stmt,
) *analysis.Diagnostic {
	x, ok := extractSlicesSortCallArg(pass, s0)
	if !ok {
		return nil
	}
	y, ok := extractSlicesSortCallArg(pass, s1)
	if !ok {
		return nil
	}

	exprStmt, ok := s2.(*ast.ExprStmt)
	if !ok {
		return nil
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}

	testifyCall := NewCallMeta(pass, callExpr)
	if testifyCall == nil {
		return nil
	}
	if testifyCall.Fn.NameFTrimmed != "True" {
		return nil
	}
	if len(testifyCall.Args) < 1 {
		return nil
	}

	eqX, eqY, ok := extractSlicesEqualArgs(pass, testifyCall.Args[0])
	if !ok {
		return nil
	}

	xStr := analysisutil.NodeString(pass.Fset, x)
	yStr := analysisutil.NodeString(pass.Fset, y)
	eqXStr := analysisutil.NodeString(pass.Fset, eqX)
	eqYStr := analysisutil.NodeString(pass.Fset, eqY)

	// The slices being sorted must be the same as those being compared.
	if (xStr != eqXStr || yStr != eqYStr) && (xStr != eqYStr || yStr != eqXStr) {
		return nil
	}

	proposedFn := "ElementsMatch"
	if testifyCall.Fn.IsFmt {
		proposedFn += "f"
	}
	msg := fmt.Sprintf("use %s.%s", testifyCall.SelectorXStr, proposedFn)

	// Replace only the first arg (slices.Equal(x, y) → x, y), preserving any trailing msgAndArgs.
	argReplaceEdit := analysis.TextEdit{
		Pos:     testifyCall.Args[0].Pos(),
		End:     testifyCall.Args[0].End(),
		NewText: formatAsCallArgs(pass, eqX, eqY),
	}
	fix := newSuggestedFuncReplacement(testifyCall, "ElementsMatch", argReplaceEdit)
	// Prepend two separate sort-removal edits (one per statement) so that
	// any comments or blank lines between the sort calls and the assertion are preserved.
	file := pass.Fset.File(s0.Pos())
	s0LineStart := file.LineStart(file.Line(s0.Pos()))
	s1LineStart := file.LineStart(file.Line(s1.Pos()))
	s1LineEnd := s2.Pos() // default: everything up to assertion
	if nextLine := file.Line(s1.End()) + 1; nextLine <= file.LineCount() {
		s1LineEnd = file.LineStart(nextLine)
	}
	fix.TextEdits = append([]analysis.TextEdit{
		{Pos: s0LineStart, End: s1LineStart, NewText: []byte{}},
		{Pos: s1LineStart, End: s1LineEnd, NewText: []byte{}},
	}, fix.TextEdits...)

	return newDiagnostic(checker.Name(), testifyCall, msg, fix)
}

// extractSlicesSortCallArg returns the first argument of a slices.Sort call,
// if the statement is an expression statement containing such a call.
func extractSlicesSortCallArg(pass *analysis.Pass, stmt ast.Stmt) (ast.Expr, bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	if !isSlicesSortCall(pass, callExpr) {
		return nil, false
	}
	if len(callExpr.Args) < 1 {
		return nil, false
	}
	return callExpr.Args[0], true
}

// extractSlicesEqualArgs returns the two arguments of a slices.Equal call,
// if the expression is such a call.
func extractSlicesEqualArgs(pass *analysis.Pass, expr ast.Expr) (x, y ast.Expr, found bool) {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, nil, false
	}
	if !isSlicesEqualCall(pass, callExpr) {
		return nil, nil, false
	}
	if len(callExpr.Args) < 2 {
		return nil, nil, false
	}
	return callExpr.Args[0], callExpr.Args[1], true
}

func isSlicesSortCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "slices", "Sort") ||
		isPkgFnCall(pass, ce, "golang.org/x/exp/slices", "Sort")
}

func isSlicesEqualCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "slices", "Equal") ||
		isPkgFnCall(pass, ce, "golang.org/x/exp/slices", "Equal")
}
