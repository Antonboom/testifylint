package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

const requireErrorReport = "for error assertions use require"

// RequireError detects error assertions like
//
// assert.Error(t, err) // s.Error(err), s.Assert().Error(err)
// assert.ErrorIs(t, err, io.EOF)
// assert.ErrorAs(t, err, &target)
// assert.EqualError(t, err, "end of file")
// assert.ErrorContains(t, err, "end of file")
// assert.NoError(t, err)
// assert.NotErrorIs(t, err, io.EOF)
//
// and requires
//
// require.Error(t, err) // s.Require().Error(err), s.Require().Error(err)
// require.ErrorIs(t, err, io.EOF)
// require.ErrorAs(t, err, &target)
// ...
//
// RequireError ignores:
// - non-negated assertions in the `if` condition;
// - assertions in the bool expression;
// - the entire `if-else[-if]` block, if there is an assertion in any `if` condition;
// - the last assertion in the block, if there are no methods/functions calls after it;
// - assertions in an explicit goroutine (including `http.Handler` and inline `sync.WaitGroup.Go` calls);
// - assertions in an explicit testing cleanup function or suite teardown methods;
// - sequence of NoError assertions.
//
// Indirect goroutine callbacks, such as `go callback()` or `wg.Go(callback)`, are not supported.
//
// RequireError also reports and provides a fix for negated error assertions in an
// if condition when the if body consists solely of a return or continue statement
// and there is no else clause, e.g.:
//
//	if !assert.NoError(t, err) {
//	   return
//	}
//
// requires
//
// require.NoError(t, err)
//
// For the same pattern applied to non-error assertions, see [NegatedAssert].
type RequireError struct {
	fnPattern *regexp.Regexp
}

// NewRequireError constructs RequireError checker.
func NewRequireError() *RequireError { return new(RequireError) }
func (RequireError) Name() string    { return "require-error" }

func (checker *RequireError) SetFnPattern(p *regexp.Regexp) *RequireError {
	if p != nil {
		checker.fnPattern = p
	}
	return checker
}

func (checker RequireError) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	callsByFunc := make(map[funcID][]*callMeta)

	// Stage 1. Collect meta information about any calls inside functions.

	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}
		if len(stack) < 3 {
			return true
		}

		fID := findSurroundingFunc(pass, stack)
		if fID == nil {
			return true
		}

		_, prevIsIfStmt := stack[len(stack)-2].(*ast.IfStmt)
		_, prevIsAssignStmt := stack[len(stack)-2].(*ast.AssignStmt)
		_, prevPrevIsIfStmt := stack[len(stack)-3].(*ast.IfStmt)
		inIfCond := prevIsIfStmt || (prevPrevIsIfStmt && prevIsAssignStmt)

		// Detect !assert.xxx() in if condition, treating ParenExpr as transparent.
		// Handles patterns like: if !assert.xxx() {}, if (!assert.xxx()) {},
		// if !assert.xxx() || !assert.yyy() {}, etc.
		negatedInIfCond := callIsNegatedInIfCond(stack)
		if negatedInIfCond {
			inIfCond = true
		}

		_, inBoolExpr := stack[len(stack)-2].(*ast.BinaryExpr)

		// Also detect !assert.xxx() inside a BinaryExpr in an if condition.
		// Stack pattern: [..., IfStmt, BinaryExpr, ..., CallExpr]
		if !inBoolExpr && negatedInIfCond {
			// If the call is negated-in-if-cond and there's a BinaryExpr on the path,
			// mark as inBoolExpr too.
			for j := len(stack) - 2; j >= 0; j-- {
				if _, ok := stack[j].(*ast.IfStmt); ok {
					break
				}
				if _, ok := stack[j].(*ast.BinaryExpr); ok {
					inBoolExpr = true
					break
				}
			}
		}

		callExpr := node.(*ast.CallExpr)
		testifyCall := NewCallMeta(pass, callExpr)

		call := &callMeta{
			call:            callExpr,
			testifyCall:     testifyCall,
			rootIf:          findRootIf(stack),
			parentIf:        findNearestNode[*ast.IfStmt](stack),
			parentBlock:     findNearestNode[*ast.BlockStmt](stack),
			inIfCond:        inIfCond,
			inBoolExpr:      inBoolExpr,
			inNoErrorSeq:    false, // Will be filled in below.
			negatedInIfCond: negatedInIfCond,
		}

		callsByFunc[*fID] = append(callsByFunc[*fID], call)
		return testifyCall == nil // Do not support asserts in asserts.
	})

	// Stage 2. Analyze calls and block context.

	var diagnostics []analysis.Diagnostic

	callsByBlock := map[*ast.BlockStmt][]*callMeta{}
	for _, calls := range callsByFunc {
		for _, c := range calls {
			if b := c.parentBlock; b != nil {
				callsByBlock[b] = append(callsByBlock[b], c)
			}
		}
	}

	markCallsInNoErrorSequence(callsByBlock)

	// Stage 2a. Identify fixable negated-if patterns for error assertions.
	// A negated-if pattern is fixable when:
	//   - All assertions in the if condition are negated error assertions (via ||)
	//   - The if body is a single return or continue statement
	//   - There is no else clause
	fixableIfs := make(map[*ast.IfStmt]*fixableIfInfo)

	for _, calls := range callsByFunc {
		negatedIfGroups := make(map[*ast.IfStmt][]*callMeta)
		for _, c := range calls {
			if !c.negatedInIfCond || c.parentIf == nil || c.testifyCall == nil {
				continue
			}
			// Skip if the if is part of an else-if chain: when rootIf != parentIf,
			// parentIf is directly nested inside another if's Else field. Replacing
			// it in isolation would leave a dangling "else" in the parent.
			if c.rootIf != c.parentIf {
				continue
			}
			if !c.testifyCall.IsAssert {
				continue
			}
			// Only handle error assertions — the general case is covered by [NegatedAssert].
			switch c.testifyCall.Fn.NameFTrimmed {
			case "Error", "ErrorIs", "ErrorAs", "EqualError", "ErrorContains", "NoError", "NotErrorIs":
				negatedIfGroups[c.parentIf] = append(negatedIfGroups[c.parentIf], c)
			}
		}

		for ifStmt, group := range negatedIfGroups {
			// Skip if-statements with an init clause: replacing them would drop
			// the init statement and break compilation due to scope changes.
			if ifStmt.Init != nil || !requireErrorSimpleBody(ifStmt) || ifStmt.Else != nil {
				continue
			}
			// Check that ALL top-level calls in the if condition are negated
			// error assertions (i.e. our group covers the entire condition).
			condCalls, allNegatedOr := collectRequireErrorNegatedOrCalls(ifStmt.Cond)
			if !allNegatedOr || len(condCalls) != len(group) {
				continue
			}
			groupSet := make(map[*ast.CallExpr]struct{}, len(group))
			for _, c := range group {
				groupSet[c.call] = struct{}{}
			}
			allMatch := true
			for _, cc := range condCalls {
				if _, ok := groupSet[cc]; !ok {
					allMatch = false
					break
				}
			}
			if !allMatch {
				continue
			}
			// Sort by source position to keep a stable replacement order.
			sorted := make([]*callMeta, len(group))
			copy(sorted, group)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].call.Pos() < sorted[j].call.Pos()
			})
			fixableIfs[ifStmt] = &fixableIfInfo{calls: sorted}
		}
	}

	// Tracks if statements that have already been reported (compound || case).
	reportedIfs := make(map[*ast.IfStmt]bool)

	for funcInfo, calls := range callsByFunc {
		for i, c := range calls {
			if m := funcInfo.meta; m.isTestCleanup || m.isGoroutine || m.isHTTPHandler {
				continue
			}

			if c.testifyCall == nil {
				continue
			}
			if c.testifyCall.Obj != nil && c.testifyCall.Obj.Name() != "Assertions" {
				continue
			}
			if !c.testifyCall.IsAssert {
				continue
			}
			switch c.testifyCall.Fn.NameFTrimmed {
			default:
				continue
			case "Error", "ErrorIs", "ErrorAs", "EqualError", "ErrorContains", "NoError", "NotErrorIs":
			}

			// Handle negated-if error patterns before the normal skip logic.
			// This ensures require-error catches these even when it is the only
			// checker enabled (the general case is also caught by [NegatedAssert]).
			if c.negatedInIfCond && c.parentIf != nil {
				diagnostics = checker.handleNegatedIfCall(pass, c, diagnostics, fixableIfs, reportedIfs)
				continue
			}

			if needToSkipBasedOnContext(c, i, calls, callsByBlock) {
				continue
			}
			if p := checker.fnPattern; p != nil && !p.MatchString(c.testifyCall.Fn.Name) {
				continue
			}

			diagnostics = append(diagnostics,
				*newDiagnostic(checker.Name(), c.testifyCall, requireErrorReport))
		}
	}

	return diagnostics
}

type fixableIfInfo struct {
	calls []*callMeta // Ordered error-assertion calls in the if condition.
}

// handleNegatedIfCall processes a call that is in a negated if condition.
// If the if statement is fixable, it is reported once with a SuggestedFix.
// Otherwise the call is skipped silently.
func (checker RequireError) handleNegatedIfCall(
	pass *analysis.Pass,
	c *callMeta,
	diagnostics []analysis.Diagnostic,
	fixableIfs map[*ast.IfStmt]*fixableIfInfo,
	reportedIfs map[*ast.IfStmt]bool,
) []analysis.Diagnostic {
	info := fixableIfs[c.parentIf]
	if info == nil {
		// Non-fixable negated-if (e.g. mixed && conditions or non-error assertions): skip.
		return diagnostics
	}
	if reportedIfs[c.parentIf] {
		return diagnostics // Already reported for this if statement.
	}
	if p := checker.fnPattern; p != nil && !p.MatchString(c.testifyCall.Fn.Name) {
		return diagnostics
	}
	reportedIfs[c.parentIf] = true
	if fix, ok := buildRequireErrorNegatedIfFix(pass, c.parentIf, info.calls); ok {
		return append(diagnostics, *newDiagnostic(checker.Name(), c.testifyCall, requireErrorReport, fix))
	}
	return append(diagnostics, *newDiagnostic(checker.Name(), c.testifyCall, requireErrorReport))
}

// requireErrorSimpleBody reports whether the body of ifStmt consists of exactly
// one return or continue statement. Other statement types (e.g. assignments,
// expression statements) are intentionally excluded: they imply side effects that
// make the if-block semantically non-trivial to replace with a bare require call.
func requireErrorSimpleBody(ifStmt *ast.IfStmt) bool {
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

// collectRequireErrorNegatedOrCalls collects all top-level CallExpr nodes from
// a condition that consists exclusively of negated calls (UnaryExpr with !) joined
// by logical-or (||). ParenExpr nodes are treated as transparent. Returns the calls
// in left-to-right order and true iff the entire expression matches that pattern.
func collectRequireErrorNegatedOrCalls(expr ast.Expr) ([]*ast.CallExpr, bool) {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return collectRequireErrorNegatedOrCalls(e.X)
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
		left, leftOk := collectRequireErrorNegatedOrCalls(e.X)
		right, rightOk := collectRequireErrorNegatedOrCalls(e.Y)
		if !leftOk || !rightOk {
			return nil, false
		}
		return append(left, right...), true
	default:
		return nil, false
	}
}

// buildRequireErrorNegatedIfFix constructs a SuggestedFix that replaces the
// entire ifStmt with one require.XXX call per error assertion in calls (in source
// order). The require import is added to the file if it is not already present.
func buildRequireErrorNegatedIfFix(pass *analysis.Pass, ifStmt *ast.IfStmt, calls []*callMeta) (analysis.SuggestedFix, bool) {
	if len(calls) == 0 {
		return analysis.SuggestedFix{}, false
	}

	qualName, importEdit, ok := addImportFix(pass.Files, calls[0].call.Pos(), testify.RequirePkgPath)
	if !ok {
		return analysis.SuggestedFix{}, false
	}

	indent := requireErrorLineIndent(pass, ifStmt.Pos())
	requireCalls := make([]string, 0, len(calls))
	for _, c := range calls {
		callText := analysisutil.NodeString(pass.Fset, c.testifyCall.Call)
		var newCallText string
		if qualName == "" {
			// dot-import: call the function directly without a qualifier;
			// strip the old "assert." prefix entirely.
			newCallText = callText[len(c.testifyCall.SelectorXStr)+1:]
		} else {
			newCallText = qualName + callText[len(c.testifyCall.SelectorXStr):]
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
		msg = "Replace with " + calls[0].testifyCall.Fn.Name
		if len(calls) > 1 {
			msg = "Replace with require calls"
		}
	} else {
		msg = fmt.Sprintf("Replace with %s.%s", qualName, calls[0].testifyCall.Fn.Name)
		if len(calls) > 1 {
			msg = "Replace with " + qualName + " calls"
		}
	}

	return analysis.SuggestedFix{
		Message:   msg,
		TextEdits: textEdits,
	}, true
}

// requireErrorLineIndent returns the leading whitespace of the source line
// containing pos.
func requireErrorLineIndent(pass *analysis.Pass, pos token.Pos) string {
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

func needToSkipBasedOnContext(
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	callsByBlock map[*ast.BlockStmt][]*callMeta,
) bool {
	if currCall.inIfCond || currCall.inBoolExpr || currCall.inNoErrorSeq {
		return true
	}

	if currCall.rootIf != nil {
		for _, rootCall := range otherCalls {
			if (rootCall.rootIf == currCall.rootIf) && rootCall.inIfCond {
				// Skip assertions in the entire if-else[-if] block, if some of "if condition" contains assertion.
				return true
			}
		}
	}

	block := currCall.parentBlock
	blockCalls := callsByBlock[block]
	isLastCallInBlock := blockCalls[len(blockCalls)-1] == currCall

	noCallsAfter := true

	_, blockEndWithReturn := block.List[len(block.List)-1].(*ast.ReturnStmt)
	if !blockEndWithReturn {
		for i := currCallIndex + 1; i < len(otherCalls); i++ {
			nextCall := otherCalls[i]
			nextCallInElseBlock := false

			if pIf := currCall.parentIf; pIf != nil && pIf.Else != nil {
				ast.Inspect(pIf.Else, func(n ast.Node) bool {
					if n == nextCall.call {
						nextCallInElseBlock = true
						return false
					}
					return true
				})
			}

			if !nextCallInElseBlock {
				noCallsAfter = false
				break
			}
		}
	}

	// Skip assertion if this is the last operation in the test.
	return isLastCallInBlock && noCallsAfter
}

func findRootIf(stack []ast.Node) *ast.IfStmt {
	nearestIf, i := findNearestNodeWithIdx[*ast.IfStmt](stack)
	for ; i > 0; i-- {
		parent, ok := stack[i-1].(*ast.IfStmt)
		if !ok {
			break
		}
		nearestIf = parent
	}
	return nearestIf
}

func markCallsInNoErrorSequence(callsByBlock map[*ast.BlockStmt][]*callMeta) {
	for _, calls := range callsByBlock {
		for i, c := range calls {
			if c.testifyCall == nil {
				continue
			}

			var prevIsNoError bool
			if i > 0 {
				if prev := calls[i-1].testifyCall; prev != nil {
					prevIsNoError = isNoErrorAssertion(prev.Fn.Name)
				}
			}

			var nextIsNoError bool
			if i < len(calls)-1 {
				if next := calls[i+1].testifyCall; next != nil {
					nextIsNoError = isNoErrorAssertion(next.Fn.Name)
				}
			}

			if isNoErrorAssertion(c.testifyCall.Fn.Name) && (prevIsNoError || nextIsNoError) {
				calls[i].inNoErrorSeq = true
			}
		}
	}
}

type callMeta struct {
	call            *ast.CallExpr
	testifyCall     *CallMeta
	rootIf          *ast.IfStmt // The root `if` in if-else[-if] chain.
	parentIf        *ast.IfStmt // The nearest `if`, can be equal with rootIf.
	parentBlock     *ast.BlockStmt
	inIfCond        bool // True for code like `if assert.ErrorAs(t, err, &target) {`.
	inBoolExpr      bool // True for code like `assert.Error(t, err) && assert.ErrorContains(t, err, "value")`
	inNoErrorSeq    bool // True for sequence of `assert.NoError` assertions.
	negatedInIfCond bool // True for code like `if !assert.NoError(t, err) {`.
}

func isNoErrorAssertion(fnName string) bool {
	return (fnName == "NoError") || (fnName == "NoErrorf")
}

// callIsNegatedInIfCond reports whether the CallExpr at the end of stack is a
// directly-negated assertion in an if condition. ParenExpr nodes are treated as
// transparent, and any BinaryExpr nodes above the UnaryExpr are treated as
// transparent when verifying that the nearest enclosing expression is an IfStmt.
//
// NOTE: both || and && are accepted here so that all negated assertions in
// if conditions are marked as inIfCond. The narrower fixable-pattern check
// (only ||) is enforced separately in Stage 2a via collectRequireErrorNegatedOrCalls.
//
// Handles patterns such as:
//
//	if !assert.xxx() {}
//	if (!assert.xxx()) {}
//	if !assert.xxx() || !assert.yyy() {}
//	if (!assert.xxx() || !assert.yyy()) {}
//	if !assert.xxx() && !assert.yyy() {}
func callIsNegatedInIfCond(stack []ast.Node) bool {
	n := len(stack)
	if n < 3 {
		return false
	}
	// Start at the parent of the CallExpr and walk upward treating ParenExpr as transparent.
	i := n - 2
	for i >= 0 {
		if _, ok := stack[i].(*ast.ParenExpr); !ok {
			break
		}
		i--
	}
	if i < 0 {
		return false
	}
	// Must be a ! UnaryExpr directly enclosing the call (through optional parens).
	unary, ok := stack[i].(*ast.UnaryExpr)
	if !ok || unary.Op != token.NOT {
		return false
	}
	i--
	// Continue upward through ParenExpr and any BinaryExpr until we hit an IfStmt.
	for i >= 0 {
		switch stack[i].(type) {
		case *ast.ParenExpr, *ast.BinaryExpr:
			i--
		case *ast.IfStmt:
			return true
		default:
			return false
		}
	}
	return false
}
