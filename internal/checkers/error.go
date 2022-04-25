package checkers

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

func Error(pass *analysis.Pass, fn FnMeta) {
	if len(fn.Args) < 2 {
		return
	}

	switch fn.Name {
	case "NotNil", "NotNilf":
		if isError(pass, fn.Args[1]) {
			r.ReportUseFunction(pass, fn, "Error")
		}

	case "Nil", "Nilf":
		if isError(pass, fn.Args[1]) {
			r.ReportUseFunction(pass, fn, "NoError")
		}
	}
}

var errIface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isError(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	_, ok := t.Underlying().(*types.Interface)
	return ok && types.Implements(t, errIface)
}
