package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type FloatCompare struct{}

func NewFloatCompare() FloatCompare {
	return FloatCompare{}
}

func (FloatCompare) Name() string {
	return "float-compare"
}

func (checker FloatCompare) Check(pass *analysis.Pass, call CallMeta) {
	invalid := func() bool {
		switch call.Fn.Name {
		case "Equal", "Equalf", "NotEqual", "NotEqualf",
			"Greater", "Greaterf", "GreaterOrEqual", "GreaterOrEqualf",
			"Less", "Lessf", "LessOrEqual", "LessOrEqualf":
			return len(call.Args) >= 2 && isFloat(pass, call.Args[0]) && isFloat(pass, call.Args[1])

		case "True", "Truef", "False", "Falsef":
			return len(call.Args) >= 1 && isFloatCompare(pass, call.Args[0])
		}
		return false
	}()

	if invalid {
		r.ReportUseFunction(pass, checker.Name(), call, "InDelta", nil)
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
