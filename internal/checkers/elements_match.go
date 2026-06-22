package checkers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"

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

		for i := 0; i < len(stmts); {
			// checkPattern returns the first index after a consumed pattern,
			// or i+1 when no pattern matches.
			d, nextIndex := checker.checkPattern(pass, stmts, i)
			if d != nil {
				diagnostics = append(diagnostics, *d)
			}
			i = nextIndex
		}
	})

	return diagnostics
}

// checkPattern checks if consecutive statements form one of supported patterns:
//   - slices.Sort(x), slices.Sort(y), assert.True(t, slices.Equal(x, y))
//   - [optional len precheck], sort/slices sort calls and element-wise comparison loop.
func (checker ElementsMatch) checkPattern(
	pass *analysis.Pass,
	stmts []ast.Stmt,
	start int,
) (d *analysis.Diagnostic, next int) {
	if d, next := checker.checkSortAndLoopPattern(pass, stmts, start); d != nil {
		return d, next
	}
	if start+2 < len(stmts) {
		if d := checker.checkSortAndEqualPattern(pass, stmts[start], stmts[start+1], stmts[start+2]); d != nil {
			return d, start + 3
		}
	}
	return nil, start + 1
}

// checkSortAndEqualPattern checks if three consecutive statements form the pattern:
// slices.Sort(x), slices.Sort(y), assert.True(t, slices.Equal(x, y)).
func (checker ElementsMatch) checkSortAndEqualPattern(
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

// checkSortAndLoopPattern checks if statements form one of the patterns:
// [assert/require.Equal(t, len(x), len(y))]
// sort/slices.Sort(x)
// sort/slices.Sort(y)
//
//	for i := range y {
//	    assert/require.Equal(t, x[i], y[i])
//	}
func (checker ElementsMatch) checkSortAndLoopPattern(
	pass *analysis.Pass,
	stmts []ast.Stmt,
	start int,
) (d *analysis.Diagnostic, next int) {
	if start+2 >= len(stmts) {
		return nil, start + 1
	}

	idx := start
	var (
		lenX, lenY ast.Expr
		hasLen     bool
	)
	if x, y, ok := extractLenEqualCallArgs(pass, stmts[idx]); ok {
		hasLen = true
		lenX = x
		lenY = y
		idx++
	}
	if idx+2 >= len(stmts) {
		return nil, start + 1
	}

	sortX, ok := extractSortCallArg(pass, stmts[idx])
	if !ok {
		return nil, start + 1
	}
	sortY, ok := extractSortCallArg(pass, stmts[idx+1])
	if !ok {
		return nil, start + 1
	}
	loop, ok := stmts[idx+2].(*ast.RangeStmt)
	if !ok {
		return nil, start + 1
	}

	loopCall, loopX, loopY, ok := extractLoopEqualCall(pass, loop)
	if !ok {
		return nil, start + 1
	}

	sortXStr := analysisutil.NodeString(pass.Fset, sortX)
	sortYStr := analysisutil.NodeString(pass.Fset, sortY)
	loopXStr := analysisutil.NodeString(pass.Fset, loopX)
	loopYStr := analysisutil.NodeString(pass.Fset, loopY)
	if !matchesEitherOrder(sortXStr, sortYStr, loopXStr, loopYStr) {
		return nil, start + 1
	}
	if hasLen {
		lenXStr := analysisutil.NodeString(pass.Fset, lenX)
		lenYStr := analysisutil.NodeString(pass.Fset, lenY)
		if !matchesEitherOrder(sortXStr, sortYStr, lenXStr, lenYStr) {
			return nil, start + 1
		}
	}

	proposedFn := "ElementsMatch"
	if loopCall.Fn.IsFmt {
		proposedFn += "f"
	}
	msg := fmt.Sprintf("use %s.%s", loopCall.SelectorXStr, proposedFn)

	// Keep the fix only for structurally-safe pattern without extra length assertion.
	if hasLen {
		return newDiagnostic(checker.Name(), loopCall, msg), idx + 3
	}

	file := pass.Fset.File(stmts[idx].Pos())
	if file == nil {
		return newDiagnostic(checker.Name(), loopCall, msg), idx + 3
	}
	sortXLineStart := file.LineStart(file.Line(stmts[idx].Pos()))
	sortYLineStart := file.LineStart(file.Line(stmts[idx+1].Pos()))
	deleteEndPos := stmts[idx+1].End()
	nextLineAfterSort := file.Line(stmts[idx+1].End()) + 1
	if nextLineAfterSort <= file.LineCount() {
		deleteEndPos = file.LineStart(nextLineAfterSort)
	}

	replaceLoopEdit := analysis.TextEdit{
		Pos:     loop.Pos(),
		End:     loop.End(),
		NewText: buildElementsMatchCallText(pass, loopCall, loopX, loopY),
	}
	fix := analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace `%s` with `%s`", loopCall.Fn.Name, proposedFn),
		TextEdits: []analysis.TextEdit{
			{Pos: sortXLineStart, End: sortYLineStart, NewText: []byte{}},
			{Pos: sortYLineStart, End: deleteEndPos, NewText: []byte{}},
			replaceLoopEdit,
		},
	}
	return newDiagnostic(checker.Name(), loopCall, msg, fix), idx + 3
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

// extractSortCallArg returns the first argument of a sort/slices sort call.
func extractSortCallArg(pass *analysis.Pass, stmt ast.Stmt) (ast.Expr, bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	if !isSortCall(pass, callExpr) {
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

func extractLenEqualCallArgs(pass *analysis.Pass, stmt ast.Stmt) (x, y ast.Expr, found bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, nil, false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, nil, false
	}

	testifyCall := NewCallMeta(pass, callExpr)
	if testifyCall == nil {
		return nil, nil, false
	}
	if testifyCall.Fn.NameFTrimmed != "Equal" || len(testifyCall.Args) < 2 {
		return nil, nil, false
	}

	x, ok = extractLenArg(pass, testifyCall.Args[0])
	if !ok {
		return nil, nil, false
	}
	y, ok = extractLenArg(pass, testifyCall.Args[1])
	if !ok {
		return nil, nil, false
	}
	return x, y, true
}

func extractLenArg(pass *analysis.Pass, expr ast.Expr) (ast.Expr, bool) {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	id, ok := ce.Fun.(*ast.Ident)
	if !ok || id.Name != "len" {
		return nil, false
	}
	// If len is shadowed by a package-level identifier, reject this match.
	if obj := pass.TypesInfo.ObjectOf(id); obj != nil && obj.Pkg() != nil {
		return nil, false
	}
	if len(ce.Args) != 1 {
		return nil, false
	}
	return ce.Args[0], true
}

func extractLoopEqualCall(pass *analysis.Pass, loop *ast.RangeStmt) (call *CallMeta, x, y ast.Expr, found bool) {
	idx, ok := loop.Key.(*ast.Ident)
	if !ok || idx.Name == "_" {
		return nil, nil, nil, false
	}
	if loop.Value != nil {
		return nil, nil, nil, false
	}
	if loop.Tok != token.DEFINE && loop.Tok != token.ASSIGN {
		return nil, nil, nil, false
	}
	if len(loop.Body.List) != 1 {
		return nil, nil, nil, false
	}

	exprStmt, ok := loop.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return nil, nil, nil, false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, nil, nil, false
	}

	testifyCall := NewCallMeta(pass, callExpr)
	if testifyCall == nil {
		return nil, nil, nil, false
	}
	if testifyCall.Fn.NameFTrimmed != "Equal" || len(testifyCall.Args) < 2 {
		return nil, nil, nil, false
	}

	x, ok = extractIndexedBy(testifyCall.Args[0], idx.Name)
	if !ok {
		return nil, nil, nil, false
	}
	y, ok = extractIndexedBy(testifyCall.Args[1], idx.Name)
	if !ok {
		return nil, nil, nil, false
	}

	loopRangeStr := analysisutil.NodeString(pass.Fset, loop.X)
	xStr := analysisutil.NodeString(pass.Fset, x)
	yStr := analysisutil.NodeString(pass.Fset, y)
	if loopRangeStr != xStr && loopRangeStr != yStr {
		return nil, nil, nil, false
	}

	return testifyCall, x, y, true
}

func extractIndexedBy(expr ast.Expr, idxName string) (ast.Expr, bool) {
	ie, ok := expr.(*ast.IndexExpr)
	if !ok {
		return nil, false
	}
	id, ok := ie.Index.(*ast.Ident)
	if !ok || id.Name != idxName {
		return nil, false
	}
	return ie.X, true
}

func matchesEitherOrder(a1, a2, b1, b2 string) bool {
	return (a1 == b1 && a2 == b2) || (a1 == b2 && a2 == b1)
}

func buildElementsMatchCallText(pass *analysis.Pass, call *CallMeta, x, y ast.Expr) []byte {
	var args []ast.Expr
	if call.IsPkg && len(call.ArgsRaw) != 0 {
		args = append(args, call.ArgsRaw[0])
	}
	args = append(args, x, y)
	if len(call.Args) > 2 {
		args = append(args, call.Args[2:]...)
	}

	fnName := "ElementsMatch"
	if call.Fn.IsFmt {
		fnName += "f"
	}

	var buf bytes.Buffer
	buf.WriteString(call.SelectorXStr)
	buf.WriteString(".")
	buf.WriteString(fnName)
	buf.WriteString("(")
	buf.Write(formatAsCallArgs(pass, args...))
	buf.WriteString(")")
	return buf.Bytes()
}

func isSortCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isSlicesSortCall(pass, ce) || isSortPackageSliceCall(pass, ce)
}

func isSlicesSortCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "slices", "Sort") ||
		isPkgFnCall(pass, ce, "golang.org/x/exp/slices", "Sort")
}

func isSortPackageSliceCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "sort", "Slice")
}

func isSlicesEqualCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "slices", "Equal") ||
		isPkgFnCall(pass, ce, "golang.org/x/exp/slices", "Equal")
}
