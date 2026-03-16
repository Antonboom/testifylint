package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

const (
	wrongTNilReport   = "do not pass nil as testing.T"
	wrongTFreshReport = "do not pass freshly created testing.T"
	wrongTScopeReport = "assertion object %s was created with outer scope's testing.T"
)

// WrongT detects situations where testify assertions are called with a wrong testing.T.
//
// Case 1 (nil or freshly-created testing.T passed to package-level assertions):
//
//	func TestFoo(t *testing.T) {
//		assert.Equal(nil, expected, actual)  ❌
//
//		u := &testing.T{}
//		assert.Equal(u, expected, actual)    ❌
//	}
//
// Case 2 (Assertions object created with outer t used inside t.Run subtest):
//
//	func TestFoo(t *testing.T) {
//		a := assert.New(t)
//		t.Run("sub", func(t *testing.T) {
//			a.Equal(expected, actual)  ❌ a was created with outer t
//		})
//	}
type WrongT struct{}

// NewWrongT constructs WrongT checker.
func NewWrongT() *WrongT    { return new(WrongT) }
func (WrongT) Name() string { return "wrong-t" }

func (checker WrongT) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	// Phase 1: collect variables assigned from freshly-created *testing.T (for Case 1 indirect).
	freshTVars := collectFreshTestingTVars(pass, insp)

	// Phase 2: collect variables assigned from assert.New / require.New (for Case 2).
	assertVars := collectAssertionVars(pass, insp)

	// Phase 3: check package-level testify calls for nil / fresh testing.T (Case 1).
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		ce := node.(*ast.CallExpr)
		call := NewCallMeta(pass, ce)
		if call == nil || !call.IsPkg {
			return
		}

		sig := call.Fn.Signature
		if !sigExpectsTestingTFirstParam(sig) {
			return
		}
		if len(call.ArgsRaw) == 0 {
			return
		}
		tArg := call.ArgsRaw[0]

		if isNil(tArg) {
			d := newDiagnostic(checker.Name(), call, wrongTNilReport)
			diagnostics = append(diagnostics, *d)
			return
		}

		if isFreshTestingTExpr(pass, tArg) {
			d := newDiagnostic(checker.Name(), call, wrongTFreshReport)
			diagnostics = append(diagnostics, *d)
			return
		}

		if tArgIdent, ok := tArg.(*ast.Ident); ok {
			obj := pass.TypesInfo.ObjectOf(tArgIdent)
			if obj != nil && freshTVars[obj] {
				d := newDiagnostic(checker.Name(), call, wrongTFreshReport)
				diagnostics = append(diagnostics, *d)
			}
		}
	})

	// Phase 4: check t.Run callbacks for assertion object used from outer scope (Case 2).
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		ce := node.(*ast.CallExpr)
		if !isSubTestRun(pass, ce) {
			return
		}

		// Find function literal callback with a testing.T parameter.
		var callback *ast.FuncLit
		for _, arg := range ce.Args {
			if fl, ok := arg.(*ast.FuncLit); ok && hasTestingTParam(pass, fl.Type) {
				callback = fl
				break
			}
		}
		if callback == nil {
			return
		}

		callbackBodyStart := callback.Body.Lbrace

		// Obtain the types.Object for the callback's own *testing.T parameter so that
		// we can verify rebindings actually use the subtest's t, not some other TestingT.
		callbackTParam := testingTParamObj(pass, callback)

		// Collect the earliest position where each outer-scope assertion object is
		// rebound inside the callback via a plain `=` assignment to assert.New / require.New,
		// where the New call's first argument is the callback's own *testing.T parameter.
		// E.g.:  a = assert.New(t)  — t is the callback's t, not an outer one.
		// Calls on `a` that occur AFTER such a rebinding are suppressed to avoid false positives.
		innerRebindings := collectInnerRebindings(pass, callback, callbackTParam, callbackBodyStart, assertVars)

		// Walk callback body looking for calls on assertion objects from outer scope.
		// Don't recurse into nested function literals to avoid double-reporting –
		// the inner t.Run callbacks are processed separately by this same loop.
		ast.Inspect(callback.Body, func(n ast.Node) bool {
			// Stop recursion at nested function literals.
			if _, isFuncLit := n.(*ast.FuncLit); isFuncLit {
				return false
			}
			callCE, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			se, ok := callCE.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := se.X.(*ast.Ident)
			if !ok {
				return true
			}
			obj := pass.TypesInfo.ObjectOf(id)
			if obj == nil {
				return true
			}
			if assertVars[obj] && obj.Pos() < callbackBodyStart {
				// Suppress if the object was rebound before this call.
				if rebindPos, rebound := innerRebindings[obj]; rebound && rebindPos < callCE.Pos() {
					return true
				}
				d := newDiagnostic(checker.Name(), callCE,
					fmt.Sprintf(wrongTScopeReport, id.Name))
				diagnostics = append(diagnostics, *d)
			}
			return true
		})
	})

	return diagnostics
}

// collectInnerRebindings collects the earliest position where each outer-scope assertion object
// is rebound inside the callback via a plain `=` assignment to assert.New / require.New,
// where the New call's first argument is the callback's own *testing.T parameter.
func collectInnerRebindings(
	pass *analysis.Pass,
	callback *ast.FuncLit,
	callbackTParam types.Object,
	callbackBodyStart token.Pos,
	assertVars map[types.Object]bool,
) map[types.Object]token.Pos {
	result := make(map[types.Object]token.Pos)
	if callbackTParam == nil {
		return result
	}
	ast.Inspect(callback.Body, func(n ast.Node) bool {
		if _, isFuncLit := n.(*ast.FuncLit); isFuncLit {
			return false
		}
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN {
			return true
		}
		for i, rhs := range as.Rhs {
			callExpr, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}
			cm := NewCallMeta(pass, callExpr)
			if cm == nil || !cm.IsPkg || cm.Fn.Name != "New" {
				continue
			}
			// Only suppress when assert.New is called with the callback's own t.
			if len(cm.ArgsRaw) == 0 {
				continue
			}
			argIdent, ok := cm.ArgsRaw[0].(*ast.Ident)
			if !ok {
				continue
			}
			if pass.TypesInfo.ObjectOf(argIdent) != callbackTParam {
				continue
			}
			if i >= len(as.Lhs) {
				continue
			}
			lhsID, ok := as.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}
			obj := pass.TypesInfo.ObjectOf(lhsID)
			if obj == nil || !assertVars[obj] || obj.Pos() >= callbackBodyStart {
				continue
			}
			// Record the earliest rebinding position for this object.
			if existing, seen := result[obj]; !seen || as.Pos() < existing {
				result[obj] = as.Pos()
			}
		}
		return true
	})
	return result
}

// collectFreshTestingTVars collects all variables that are assigned from freshly created *testing.T.
// E.g.: u := &testing.T{} or u := new(testing.T).
// Variables that are subsequently reassigned (via =) to a non-fresh expression are excluded,
// to avoid false positives for patterns like: u := &testing.T{}; u = t; assert.Equal(u, ...).
func collectFreshTestingTVars(pass *analysis.Pass, insp *inspector.Inspector) map[types.Object]bool {
	result := make(map[types.Object]bool)

	insp.Preorder([]ast.Node{(*ast.AssignStmt)(nil)}, func(node ast.Node) {
		as := node.(*ast.AssignStmt)
		switch as.Tok {
		case token.DEFINE:
			for i, rhs := range as.Rhs {
				if !isFreshTestingTExpr(pass, rhs) {
					continue
				}
				if i >= len(as.Lhs) {
					continue
				}
				id, ok := as.Lhs[i].(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.ObjectOf(id)
				if obj != nil {
					result[obj] = true
				}
			}
		case token.ASSIGN:
			// If a variable from result is reassigned to a non-fresh expression, remove it
			// to avoid false positives when the variable is later used with a valid TestingT.
			// For multi-return assignments (a, b = f()), where LHS and RHS counts differ,
			// we cannot determine per-element freshness, so conservatively remove all
			// LHS variables that are in the fresh set.
			if len(as.Lhs) != len(as.Rhs) {
				for _, lhs := range as.Lhs {
					id, ok := lhs.(*ast.Ident)
					if !ok {
						continue
					}
					obj := pass.TypesInfo.ObjectOf(id)
					if obj != nil {
						delete(result, obj)
					}
				}
				return
			}
			for i, lhs := range as.Lhs {
				id, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.ObjectOf(id)
				if obj == nil || !result[obj] {
					continue
				}
				if !isFreshTestingTExpr(pass, as.Rhs[i]) {
					delete(result, obj)
				}
			}
		default:
		}
	})

	return result
}

// collectAssertionVars collects all variables assigned from assert.New(t) or require.New(t).
func collectAssertionVars(pass *analysis.Pass, insp *inspector.Inspector) map[types.Object]bool {
	result := make(map[types.Object]bool)

	insp.Preorder([]ast.Node{(*ast.AssignStmt)(nil)}, func(node ast.Node) {
		as := node.(*ast.AssignStmt)
		if as.Tok != token.DEFINE {
			return
		}
		for i, rhs := range as.Rhs {
			ce, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}
			cm := NewCallMeta(pass, ce)
			if cm == nil || !cm.IsPkg || cm.Fn.Name != "New" {
				continue
			}
			if i >= len(as.Lhs) {
				continue
			}
			id, ok := as.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}
			obj := pass.TypesInfo.ObjectOf(id)
			if obj != nil {
				result[obj] = true
			}
		}
	})

	return result
}

// sigExpectsTestingTFirstParam returns true if the function signature's first parameter is a testify TestingT.
func sigExpectsTestingTFirstParam(sig *types.Signature) bool {
	if sig.Params().Len() == 0 {
		return false
	}
	return isTestifyTestingTType(sig.Params().At(0).Type())
}

// isTestifyTestingTType returns true if t is the TestingT interface from testify/assert or testify/require.
func isTestifyTestingTType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Name() != "TestingT" {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	return analysisutil.IsPkg(pkg, testify.AssertPkgName, testify.AssertPkgPath) ||
		analysisutil.IsPkg(pkg, testify.RequirePkgName, testify.RequirePkgPath)
}

// isFreshTestingTExpr returns true if expr is a freshly created *testing.T:
//   - &testing.T{}
//   - new(testing.T)
func isFreshTestingTExpr(pass *analysis.Pass, expr ast.Expr) bool {
	// Check for &testing.T{}.
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		if cl, ok := unary.X.(*ast.CompositeLit); ok {
			return isTestingTTypeExpr(pass, cl.Type)
		}
	}
	// Check for new(testing.T).
	if ce, ok := expr.(*ast.CallExpr); ok {
		if isIdentWithName("new", ce.Fun) && len(ce.Args) == 1 {
			return isTestingTTypeExpr(pass, ce.Args[0])
		}
	}
	return false
}

// isTestingTTypeExpr returns true if expr is the type expression testing.T.
func isTestingTTypeExpr(pass *analysis.Pass, expr ast.Expr) bool {
	if expr == nil {
		return false
	}
	se, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if se.Sel == nil || se.Sel.Name != "T" {
		return false
	}
	obj := pass.TypesInfo.ObjectOf(se.Sel)
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	return pkg != nil && pkg.Path() == "testing"
}

// testingTParamObj returns the types.Object for the first *testing.T parameter of a function literal,
// or nil if no such parameter exists.
func testingTParamObj(pass *analysis.Pass, fl *ast.FuncLit) types.Object {
	if fl.Type == nil || fl.Type.Params == nil {
		return nil
	}
	for _, field := range fl.Type.Params.List {
		if !implementsTestingT(pass, field.Type) {
			continue
		}
		for _, name := range field.Names {
			if obj := pass.TypesInfo.ObjectOf(name); obj != nil {
				return obj
			}
		}
	}
	return nil
}
