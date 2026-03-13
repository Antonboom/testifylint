package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isSubTestRun(pass *analysis.Pass, ce *ast.CallExpr) bool {
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || se.Sel == nil {
		return false
	}
	return (implementsTestingT(pass, se.X) || implementsTestifySuite(pass, se.X)) && se.Sel.Name == "Run"
}

func isTestingFuncOrMethod(pass *analysis.Pass, fd *ast.FuncDecl) bool {
	return hasTestingTParam(pass, fd.Type) || isSuiteMethod(pass, fd)
}

func isTestingAnonymousFunc(pass *analysis.Pass, ft *ast.FuncType) bool {
	return hasTestingTParam(pass, ft)
}

func hasTestingTParam(pass *analysis.Pass, ft *ast.FuncType) bool {
	if ft == nil || ft.Params == nil {
		return false
	}

	for _, param := range ft.Params.List {
		if implementsTestingT(pass, param.Type) {
			return true
		}
	}
	return false
}

// hasStdTestingTParam reports whether the function type has a parameter of type *testing.T
// from the standard library. Unlike hasTestingTParam, this does not require testify assert
// or require packages to be imported.
func hasStdTestingTParam(pass *analysis.Pass, ft *ast.FuncType) bool {
	if ft == nil || ft.Params == nil {
		return false
	}
	for _, param := range ft.Params.List {
		if isStdTestingTType(pass, param.Type) {
			return true
		}
	}
	return false
}

// isStdTestingTType reports whether the expression has type *testing.T from the standard library.
func isStdTestingTType(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "testing" && obj.Name() == "T"
}
