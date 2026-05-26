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

// isWaitGroupGoCall returns true if ce is a call to (*sync.WaitGroup).Go,
// introduced in Go 1.25. Such calls run the callback in a new goroutine,
// so any assertions inside the callback must follow the same rules as
// assertions inside an explicit `go func() {...}()` statement.
func isWaitGroupGoCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || se.Sel == nil || se.Sel.Name != "Go" {
		return false
	}

	sel, ok := pass.TypesInfo.Selections[se]
	if !ok {
		return false
	}

	fn, ok := sel.Obj().(*types.Func)
	if !ok {
		return false
	}

	recvType := sel.Recv()
	if recvType == nil {
		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			return false
		}
		recv := sig.Recv()
		if recv == nil {
			return false
		}
		recvType = recv.Type()
	}

	t := recvType
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "sync" && obj.Name() == "WaitGroup"
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
