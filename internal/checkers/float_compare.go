package checkers

/*
import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

func FloatCompare(pass *analysis.Pass, fn CallMeta) {
	invalid := func() bool {
		switch fn.Name {
		case "Equal", "Equalf", "NotEqual", "NotEqualf",
			"Greater", "Greaterf", "GreaterOrEqual", "GreaterOrEqualf",
			"Less", "Lessf", "LessOrEqual", "LessOrEqualf":
			return len(fn.Args) >= 3 && isFloat(pass, fn.Args[1]) && isFloat(pass, fn.Args[2])

		case "True", "Truef", "False", "Falsef":
			return len(fn.Args) >= 2 && isFloatCompare(pass, fn.Args[1])
		}
		return false
	}

	if invalid() {
		r.ReportUseFunction(pass, fn, "InDelta")
	}
}

func isFloat(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	bt, ok := t.Underlying().(*types.Basic)
	return ok && (bt.Info()&types.IsFloat > 0)
}

func isFloatCompare(p *analysis.Pass, e ast.Expr) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return isFloat(p, be.X) && isFloat(p, be.Y)
}
*/
