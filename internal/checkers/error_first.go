package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// ErrorFirst detects situations where a function returning an error is called
// and an assertion is made on the result before asserting on the error (or
// where the error is discarded with _). For example:
//
//	res, err := myfunc()
//	assert.NotNil(t, res) // want "assert error before making other assertions"
//
// or:
//
//	res, _ := myfunc()
//	assert.NotNil(t, res) // want "error return value was discarded; assert the error before asserting the result"
//
// Valid usage:
//
//	res, err := myfunc()
//	require.NoError(t, err)
//	assert.NotNil(t, res)
type ErrorFirst struct{}

// NewErrorFirst constructs ErrorFirst checker.
func NewErrorFirst() *ErrorFirst { return new(ErrorFirst) }
func (ErrorFirst) Name() string  { return "error-first" }

// errorFirstAssignInfo stores information about a multi-return function
// assignment that includes at least one error-typed return value.
type errorFirstAssignInfo struct {
	pos         token.Pos
	errObj      *types.Var
	parentBlock *ast.BlockStmt
}

// funcAssigns tracks result and error variable → assignInfo mappings for a function.
type funcAssigns struct {
	// resultToAssign maps a result (non-error) variable to the assignInfo for
	// the multi-return call it came from.
	resultToAssign map[*types.Var]*errorFirstAssignInfo
	// errToAssign maps an error variable to the assignInfo for the multi-return
	// call it came from. Used to find canonical error assertions.
	errToAssign map[*types.Var]*errorFirstAssignInfo
}

func (checker ErrorFirst) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	callsByFunc := make(map[funcID][]*callMeta)
	assignsByFunc := make(map[funcID]*funcAssigns)

	// Stage 1a: collect assignment information for multi-return calls that
	// include at least one error-typed return value.
	insp.WithStack([]ast.Node{(*ast.AssignStmt)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		fID := findSurroundingFunc(pass, stack)
		if fID == nil {
			return true
		}

		as := node.(*ast.AssignStmt)
		if fa := assignsByFunc[*fID]; fa != nil {
			errorFirstForgetAssignedVars(pass, fa, as.Lhs)
		}

		// RHS must be a single call expression (multi-return function).
		if len(as.Rhs) != 1 {
			return true
		}
		callExpr, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		// Get the return type of the call.
		tv, ok := pass.TypesInfo.Types[callExpr]
		if !ok {
			return true
		}
		tuple, ok := tv.Type.(*types.Tuple)
		if !ok {
			// Single return value — not a multi-return call.
			return true
		}

		// Require LHS length to match tuple length.
		if len(as.Lhs) != tuple.Len() {
			return true
		}

		// Find the first error-typed position in the return values.
		errIdx := -1
		for i := range tuple.Len() {
			t := tuple.At(i).Type()
			if types.Implements(t, errorIface) || types.Implements(types.NewPointer(t), errorIface) {
				errIdx = i
				break
			}
		}
		if errIdx < 0 {
			// No error return — nothing to track.
			return true
		}

		parentBlock := findNearestNode[*ast.BlockStmt](stack)

		fa, exists := assignsByFunc[*fID]
		if !exists {
			fa = &funcAssigns{
				resultToAssign: make(map[*types.Var]*errorFirstAssignInfo),
				errToAssign:    make(map[*types.Var]*errorFirstAssignInfo),
			}
			assignsByFunc[*fID] = fa
		}

		// Build the errorFirstAssignInfo for the error variable at errIdx.
		// If the error return was discarded (_), skip tracking — this is
		// already caught by errcheck.
		if isIdentWithName("_", as.Lhs[errIdx]) {
			return true
		}

		info := &errorFirstAssignInfo{
			pos:         as.Pos(),
			parentBlock: parentBlock,
		}

		if id, ok := as.Lhs[errIdx].(*ast.Ident); ok {
			if obj, ok := pass.TypesInfo.ObjectOf(id).(*types.Var); ok {
				info.errObj = obj
				// On reassignment, overwrite the old entry for this error var.
				delete(fa.errToAssign, obj)
				fa.errToAssign[obj] = info
			}
		}

		// Register all non-error LHS variables as result variables pointing
		// at this assignInfo.
		for i := range tuple.Len() {
			t := tuple.At(i).Type()
			if types.Implements(t, errorIface) || types.Implements(types.NewPointer(t), errorIface) {
				continue
			}
			lhsExpr := as.Lhs[i]
			if isIdentWithName("_", lhsExpr) {
				continue
			}
			id, ok := lhsExpr.(*ast.Ident)
			if !ok {
				continue
			}
			obj, ok := pass.TypesInfo.ObjectOf(id).(*types.Var)
			if !ok {
				continue
			}
			// On reassignment, overwrite the old entry for this result var.
			delete(fa.resultToAssign, obj)
			fa.resultToAssign[obj] = info
		}

		return true
	})

	// Stage 1b: collect call metadata (same pattern as RequireError).
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

		_, inBoolExpr := stack[len(stack)-2].(*ast.BinaryExpr)

		callExpr := node.(*ast.CallExpr)
		testifyCall := NewCallMeta(pass, callExpr)

		call := &callMeta{
			call:        callExpr,
			testifyCall: testifyCall,
			rootIf:      findRootIf(stack),
			parentIf:    findNearestNode[*ast.IfStmt](stack),
			parentBlock: findNearestNode[*ast.BlockStmt](stack),
			inIfCond:    inIfCond,
			inBoolExpr:  inBoolExpr,
		}

		callsByFunc[*fID] = append(callsByFunc[*fID], call)
		return testifyCall == nil // Do not recurse into testify call arguments.
	})

	// Stage 2: for each function, flag non-error testify assertions that use a
	// tracked result variable whose associated error has not been canonically
	// asserted before the assertion.
	var diagnostics []analysis.Diagnostic

	for funcInfo, calls := range callsByFunc {
		if m := funcInfo.meta; m.isTestCleanup || m.isGoroutine || m.isHTTPHandler {
			continue
		}

		fa := assignsByFunc[funcInfo]
		if fa == nil {
			continue
		}
		if len(fa.resultToAssign) == 0 {
			continue
		}
		reportedAssigns := make(map[*errorFirstAssignInfo]struct{})

		for _, c := range calls {
			if c.testifyCall == nil {
				continue
			}
			if c.inIfCond || c.inBoolExpr {
				// Assertions in if-conditions are skipped as flagging targets.
				continue
			}

			// Skip canonical error assertions — they are the "good" assertions
			// that satisfy the requirement.
			if isErrorFirstErrorAssertion(c.testifyCall.Fn.NameFTrimmed) {
				continue
			}

			// Check all arguments for tracked result variables.
			for _, arg := range c.testifyCall.Args {
				id, ok := arg.(*ast.Ident)
				if !ok {
					continue
				}
				obj, ok := pass.TypesInfo.ObjectOf(id).(*types.Var)
				if !ok {
					continue
				}

				info, exists := fa.resultToAssign[obj]
				if !exists {
					continue
				}
				if _, reported := reportedAssigns[info]; reported {
					continue
				}

				// The result variable must have been assigned before this call
				// and in an enclosing or equal block.
				if !errorFirstAssignIsActiveForCall(info, c) {
					continue
				}

				if !errorFirstIsErrChecked(pass, info, c, calls) {
					diagnostics = append(diagnostics, *newDiagnostic(
						checker.Name(), c.testifyCall,
						"assert error before making other assertions",
					))
					reportedAssigns[info] = struct{}{}
				}
				// Only report once per assignment origin.
				break
			}
		}
	}

	return diagnostics
}

// errorFirstForgetAssignedVars removes variables written by an assignment from
// the tracking maps, so stale origins don't survive unrelated reassignments.
func errorFirstForgetAssignedVars(pass *analysis.Pass, fa *funcAssigns, lhs []ast.Expr) {
	for _, lhsExpr := range lhs {
		if isIdentWithName("_", lhsExpr) {
			continue
		}
		id, ok := lhsExpr.(*ast.Ident)
		if !ok {
			continue
		}
		obj, ok := pass.TypesInfo.ObjectOf(id).(*types.Var)
		if !ok {
			continue
		}
		delete(fa.resultToAssign, obj)
		delete(fa.errToAssign, obj)
	}
}

// isErrorFirstErrorAssertion returns true if the function name (with f suffix
// stripped) is one of the canonical error assertion functions.
func isErrorFirstErrorAssertion(nameFTrimmed string) bool {
	switch nameFTrimmed {
	case "Error", "ErrorIs", "ErrorAs", "EqualError", "ErrorContains", "NoError", "NotErrorIs":
		return true
	}
	return false
}

// errorFirstAssignIsActiveForCall returns true if the assignment info applies
// to the given call: the assignment happened before the call and the
// assignment's block is an ancestor of (or equal to) the call's block.
func errorFirstAssignIsActiveForCall(info *errorFirstAssignInfo, c *callMeta) bool {
	if info.parentBlock == nil || c.parentBlock == nil {
		return false
	}
	if info.pos >= c.call.Pos() {
		return false
	}
	// The call must be inside the assignment's block (including nested blocks).
	if c.call.Pos() < info.parentBlock.Pos() || c.call.Pos() > info.parentBlock.End() {
		return false
	}
	return true
}

// errorFirstIsErrChecked returns true if there is any testify assertion for
// info.errObj that appears between the assignment and the call c.
func errorFirstIsErrChecked(
	pass *analysis.Pass,
	info *errorFirstAssignInfo,
	c *callMeta,
	allCalls []*callMeta,
) bool {
	if info.errObj == nil {
		return false
	}

	for _, other := range allCalls {
		if other.testifyCall == nil {
			continue
		}
		// The error assertion must appear after the assignment.
		if other.call.Pos() <= info.pos {
			continue
		}
		// The error variable must be in non-message assertion arguments.
		if !errorFirstErrArgMatchesAssertionArgs(pass, other.testifyCall, info.errObj) {
			continue
		}

		// Case 1: error assertion is in the same block as the call, and appears before it.
		if other.parentBlock == c.parentBlock && other.call.Pos() < c.call.Pos() {
			return true
		}

		// Case 2: error assertion is in an if-condition, and the call is inside
		// the body of that if statement.
		if other.inIfCond {
			if otherIf := other.parentIf; otherIf != nil {
				if nodeContains(otherIf.Body, c.call) {
					return true
				}
			}
		}

		// Case 3: error assertion is in an ancestor block and appears before the call.
		if other.parentBlock != c.parentBlock &&
			other.call.Pos() < c.call.Pos() &&
			nodeContains(other.parentBlock, c.call) {
			return true
		}
	}
	return false
}

// errorFirstErrArgMatchesAssertionArgs returns true if any non-message argument
// of the testify call refers to the given error variable object.
func errorFirstErrArgMatchesAssertionArgs(pass *analysis.Pass, cm *CallMeta, errObj *types.Var) bool {
	for _, arg := range cm.Args[:errorFirstNonMessageArgLimit(pass, cm)] {
		id, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}
		if pass.TypesInfo.ObjectOf(id) == errObj {
			return true
		}
	}
	return false
}

func errorFirstNonMessageArgLimit(pass *analysis.Pass, cm *CallMeta) int {
	sig := cm.Fn.Signature
	if sig == nil {
		return len(cm.Args)
	}

	paramsLen := sig.Params().Len()
	if len(cm.ArgsRaw) > 0 && implementsTestingT(pass, cm.ArgsRaw[0]) {
		paramsLen--
	}
	if paramsLen < 0 {
		return 0
	}

	limit := len(cm.Args)
	if sig.Variadic() {
		// Variadic parameter is msgAndArgs; base assertion args are fixed params.
		limit = paramsLen - 1
		// Formatted assertions have explicit message arg before msgAndArgs.
		if cm.Fn.IsFmt {
			limit--
		}
	}
	if limit < 0 {
		return 0
	}
	if limit > len(cm.Args) {
		return len(cm.Args)
	}
	return limit
}

// nodeContains reports whether node n contains target anywhere in its subtree.
func nodeContains(n ast.Node, target ast.Node) bool {
	if n == nil || target == nil {
		return false
	}
	targetPos := target.Pos()
	return targetPos >= n.Pos() && targetPos < n.End()
}
