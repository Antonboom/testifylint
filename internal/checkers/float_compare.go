package checkers

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

func isFloat(expr ast.Expr, pass *analysis.Pass) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	bt, ok := t.Underlying().(*types.Basic)
	return ok && (bt.Info()&types.IsFloat == 1)
}
