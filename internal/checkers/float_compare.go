package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// FloatCompare detects situation like
//
//	assert.Equal(t, 42.42, a)
//
// and requires
//
//	assert.InEpsilon(t, 42.42, a, 0.0001) // Or assert.InDelta
type FloatCompare struct{}

// NewFloatCompare constructs FloatCompare checker.
func NewFloatCompare() FloatCompare { return FloatCompare{} }
func (FloatCompare) Name() string   { return "float-compare" }

func (checker FloatCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	invalid := func() bool {
		switch call.Fn.Name {
		case "Equal", "Equalf", "NotEqual", "NotEqualf",
			"Greater", "Greaterf", "GreaterOrEqual", "GreaterOrEqualf",
			"Less", "Lessf", "LessOrEqual", "LessOrEqualf":
			return len(call.Args) > 1 && isFloat(pass, call.Args[0]) && isFloat(pass, call.Args[1])

		case "True", "Truef":
			return len(call.Args) > 0 && isFloatCompare(pass, call.Args[0])

		case "False", "Falsef":
			return len(call.Args) > 0 && isFloatCompare(pass, call.Args[0])
		}
		return false
	}()

	if invalid {
		return newUseFunctionDiagnostic(checker.Name(), call, "InEpsilon (or InDelta)", nil)
	}
	return nil
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
