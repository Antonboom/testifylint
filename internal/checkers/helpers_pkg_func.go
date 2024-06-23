package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func isRegexpMustCompileCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "regexp", "MustCompile")
}

func isStringsContainsCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "strings", "Contains")
}

func isPkgFnCall(pass *analysis.Pass, ce *ast.CallExpr, pkg, fn string) bool {
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	errorsIsObj := analysisutil.ObjectOf(pass.Pkg, pkg, fn)
	if errorsIsObj == nil {
		return false
	}

	return analysisutil.IsObj(pass.TypesInfo, se.Sel, errorsIsObj)
}
