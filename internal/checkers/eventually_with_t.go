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
	eventuallyWithTPkgWrongTReport = "do not use outer testing.T as the first argument inside " +
		"EventuallyWithT callback, use the provided %s (*assert.CollectT)"
	eventuallyWithTRequireReport = "do not use require in EventuallyWithT callback, use %s instead"
	eventuallyWithTObjReport     = "assertion object %s was created with outer testing.T, use " +
		"assert.New(%s) inside EventuallyWithT callback"

	// assertNewFnName is the name of the constructor used to create assertion objects
	// from both the assert and require packages (assert.New / require.New).
	assertNewFnName = "New"
)

// EventuallyWithT detects situations inside assert.EventuallyWithT and
// require.EventuallyWithT callbacks where the outer testing.T is used instead of
// the provided *assert.CollectT, or where require package assertions are used
// instead of assert.
//
// Package-level assertion calls:
//
//	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
//		assert.Equal(t, expected, actual)    // ❌ use collect
//		require.NoError(t, err)              // ❌ use assert.NoError(collect, err)
//	}, timeout, interval)
//
// Assertion-object calls:
//
//	a := assert.New(t)
//	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
//		a.Equal(expected, actual)            // ❌ a was created with outer t
//	}, timeout, interval)
type EventuallyWithT struct{}

// NewEventuallyWithT constructs EventuallyWithT checker.
func NewEventuallyWithT() EventuallyWithT { return EventuallyWithT{} }
func (EventuallyWithT) Name() string      { return "eventually-with-t" }

func (checker EventuallyWithT) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	// Collect all assertion object variables (assigned from assert.New or require.New).
	assertVars := collectEventuallyAssertionVars(pass, insp)

	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		ce := node.(*ast.CallExpr)
		call := NewCallMeta(pass, ce)
		if call == nil || !call.IsPkg {
			return
		}
		if call.Fn.NameFTrimmed != "EventuallyWithT" {
			return
		}
		// EventuallyWithT(t, condition, waitFor, tick, ...)
		// After stripping the TestingT arg, Args[0] is the condition callback.
		if len(call.Args) < 1 {
			return
		}

		callback, ok := call.Args[0].(*ast.FuncLit)
		if !ok {
			return
		}

		// Find the *CollectT parameter in the callback.
		collectParam := findCollectTParam(pass, callback)
		if collectParam == nil {
			return
		}

		collectName := collectParam.Name()
		callbackBodyStart := callback.Body.Lbrace

		// Walk the callback body looking for wrong assertions.
		// All nested function literals are intentionally excluded: assertions inside
		// nested goroutines or t.Run closures have their own scoping rules and may
		// legitimately capture a different TestingT.
		ast.Inspect(callback.Body, func(n ast.Node) bool {
			if _, isFuncLit := n.(*ast.FuncLit); isFuncLit {
				return false
			}

			callCE, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			innerCall := NewCallMeta(pass, callCE)
			if innerCall == nil {
				return true
			}

			if innerCall.IsPkg {
				diagnostics = append(diagnostics, checker.checkPkgCall(pass, innerCall, collectParam, collectName)...)
			} else {
				if d := checker.checkObjCall(pass, callCE, assertVars, callbackBodyStart, collectName); d != nil {
					diagnostics = append(diagnostics, *d)
				}
			}
			return true
		})
	})

	return diagnostics
}

func (checker EventuallyWithT) checkPkgCall(
	pass *analysis.Pass,
	call *CallMeta,
	collectParam types.Object,
	collectName string,
) []analysis.Diagnostic {
	if len(call.ArgsRaw) == 0 {
		return nil
	}
	tArg := call.ArgsRaw[0]
	if !implementsTestingT(pass, tArg) {
		return nil
	}
	// If the first arg IS the collect param, or an alias with the same type, the call is correct.
	if isObjIdent(pass, tArg, collectParam) || isSameType(pass, tArg, collectParam.Type()) {
		return nil
	}

	if call.IsAssert {
		// Safe autofix: replace the wrong t argument with the collect param name.
		d := newDiagnostic(checker.Name(), call,
			fmt.Sprintf(eventuallyWithTPkgWrongTReport, collectName),
			analysis.SuggestedFix{
				Message: fmt.Sprintf("Replace argument with `%s`", collectName),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     tArg.Pos(),
						End:     tArg.End(),
						NewText: []byte(collectName),
					},
				},
			})
		return []analysis.Diagnostic{*d}
	}

	// require call: report that assert should be used with collect.
	d := newDiagnostic(checker.Name(), call,
		fmt.Sprintf(eventuallyWithTRequireReport, suggestedAssertCall(call.Fn.Name, collectName)))
	return []analysis.Diagnostic{*d}
}

func (checker EventuallyWithT) checkObjCall(
	pass *analysis.Pass,
	callCE *ast.CallExpr,
	assertVars map[types.Object]bool,
	callbackBodyStart token.Pos,
	collectName string,
) *analysis.Diagnostic {
	se, ok := callCE.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	id, ok := se.X.(*ast.Ident)
	if !ok {
		return nil
	}
	obj := pass.TypesInfo.ObjectOf(id)
	if obj == nil {
		return nil
	}
	// Only flag assertion objects created OUTSIDE the EventuallyWithT callback.
	// The position comparison is valid because variable declarations in Go always
	// appear textually before the EventuallyWithT call that contains the callback.
	// Objects created inside the callback use the callback's own collect param,
	// so we let the package-level check handle any wrong usages there.
	if assertVars[obj] && obj.Pos() < callbackBodyStart {
		return newDiagnostic(checker.Name(), callCE,
			fmt.Sprintf(eventuallyWithTObjReport, id.Name, collectName))
	}
	return nil
}

// collectEventuallyAssertionVars collects all variables assigned via assert.New or require.New
// using short variable declaration (:=).
func collectEventuallyAssertionVars(pass *analysis.Pass, insp *inspector.Inspector) map[types.Object]bool {
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
			if cm == nil || !cm.IsPkg || cm.Fn.Name != assertNewFnName {
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

// findCollectTParam returns the types.Object for the *assert.CollectT parameter
// of a function literal, or nil if no such parameter exists.
func findCollectTParam(pass *analysis.Pass, fl *ast.FuncLit) types.Object {
	if fl.Type == nil || fl.Type.Params == nil {
		return nil
	}
	for _, field := range fl.Type.Params.List {
		if isCollectTExpr(pass, field.Type) {
			for _, name := range field.Names {
				if obj := pass.TypesInfo.ObjectOf(name); obj != nil {
					return obj
				}
			}
		}
	}
	return nil
}

// isCollectTExpr returns true if expr represents *assert.CollectT.
func isCollectTExpr(pass *analysis.Pass, expr ast.Expr) bool {
	starExpr, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	t := pass.TypesInfo.TypeOf(starExpr.X)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil || obj.Name() != "CollectT" {
		return false
	}
	pkg := obj.Pkg()
	return pkg != nil && analysisutil.IsPkg(pkg, testify.AssertPkgName, testify.AssertPkgPath)
}

// isObjIdent returns true if expr is an identifier referring to the given types.Object.
func isObjIdent(pass *analysis.Pass, expr ast.Expr, obj types.Object) bool {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return pass.TypesInfo.ObjectOf(id) == obj
}

func isSameType(pass *analysis.Pass, expr ast.Expr, targetType types.Type) bool {
	exprType := pass.TypesInfo.TypeOf(expr)
	if exprType == nil || targetType == nil {
		return false
	}

	return types.Identical(exprType, targetType)
}

func suggestedAssertCall(fnName, collectName string) string {
	if fnName == assertNewFnName {
		return fmt.Sprintf("assert.%s(%s)", fnName, collectName)
	}
	return fmt.Sprintf("assert.%s(%s, ...)", fnName, collectName)
}
