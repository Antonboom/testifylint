package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// BoolCompare checks situation like
//
//	assert.Equal(t, false, result)
//
// and requires e.g.
//
//	assert.False(t, result)
type BoolCompare struct{} //

// NewBoolCompare constructs BoolCompare checker.
func NewBoolCompare() BoolCompare { return BoolCompare{} }
func (BoolCompare) Name() string  { return "bool-compare" }

func (checker BoolCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	const (
		needSimplifyMsg  = "need to simplify the check"
		simplifyCheckMsg = "Simplify the check"
	)

	switch call.Fn.Name {
	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return nil
		}

		arg1, arg2 := call.Args[0], call.Args[1]
		t1, t2 := isUntypedTrue(pass, arg1), isUntypedTrue(pass, arg2)
		f1, f2 := isUntypedFalse(pass, arg1), isUntypedFalse(pass, arg2)

		switch {
		case xor(t1, t2):
			survivingArg := arg1
			if t1 {
				survivingArg = arg2
			}

			return newUseFunctionDiagnostic(checker.Name(), call, "True",
				newSuggestedFuncReplacement(call, "True", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)

		case xor(f1, f2):
			survivingArg := arg1
			if f1 {
				survivingArg = arg2
			}

			return newUseFunctionDiagnostic(checker.Name(), call, "False",
				newSuggestedFuncReplacement(call, "False", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)
		}

	case "NotEqual", "NotEqualf":
		if len(call.Args) < 2 {
			return nil
		}

		arg1, arg2 := call.Args[0], call.Args[1]
		t1, t2 := isUntypedTrue(pass, arg1), isUntypedTrue(pass, arg2)
		f1, f2 := isUntypedFalse(pass, arg1), isUntypedFalse(pass, arg2)

		switch {
		case xor(t1, t2):
			survivingArg := arg1
			if t1 {
				survivingArg = arg2
			}

			return newUseFunctionDiagnostic(checker.Name(), call, "False",
				newSuggestedFuncReplacement(call, "False", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)

		case xor(f1, f2):
			survivingArg := arg1
			if f1 {
				survivingArg = arg2
			}

			return newUseFunctionDiagnostic(checker.Name(), call, "True",
				newSuggestedFuncReplacement(call, "True", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return nil
		}

		expr := call.Args[0]

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.NEQ)

			if survivingArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
				return newDiagnostic(checker.Name(), call, needSimplifyMsg,
					&analysis.SuggestedFix{
						Message: simplifyCheckMsg,
						TextEdits: []analysis.TextEdit{{
							Pos:     expr.Pos(),
							End:     expr.End(),
							NewText: []byte(types.ExprString(survivingArg)),
						}},
					},
				)
			}
		}

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.NEQ)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.EQL)
			arg3, ok3 := isNegation(expr)

			if survivingArg, ok := anyVal([]bool{ok1, ok2, ok3}, arg1, arg2, arg3); ok {
				return newUseFunctionDiagnostic(checker.Name(), call, "False",
					newSuggestedFuncReplacement(call, "False", analysis.TextEdit{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(types.ExprString(survivingArg)),
					}),
				)
			}
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return nil
		}

		expr := call.Args[0]

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.NEQ)

			if survivingArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
				return newDiagnostic(checker.Name(), call, needSimplifyMsg,
					&analysis.SuggestedFix{
						Message: simplifyCheckMsg,
						TextEdits: []analysis.TextEdit{{
							Pos:     expr.Pos(),
							End:     expr.End(),
							NewText: []byte(types.ExprString(survivingArg)),
						}},
					},
				)
			}
		}

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.NEQ)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.EQL)
			arg3, ok3 := isNegation(expr)

			if survivingArg, ok := anyVal([]bool{ok1, ok2, ok3}, arg1, arg2, arg3); ok {
				return newUseFunctionDiagnostic(checker.Name(), call, "True",
					newSuggestedFuncReplacement(call, "True", analysis.TextEdit{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(types.ExprString(survivingArg)),
					}),
				)
			}
		}
	}
	return nil
}

var (
	falseObj = types.Universe.Lookup("false")
	trueObj  = types.Universe.Lookup("true")
)

func isUntypedTrue(pass *analysis.Pass, e ast.Expr) bool {
	return analysisutil.IsObj(pass.TypesInfo, e, trueObj)
}

func isUntypedFalse(pass *analysis.Pass, e ast.Expr) bool {
	return analysisutil.IsObj(pass.TypesInfo, e, falseObj)
}

func isComparisonWithTrue(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) {
	return isComparisonWith(pass, e, isUntypedTrue, op)
}

func isComparisonWithFalse(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) {
	return isComparisonWith(pass, e, isUntypedFalse, op)
}

func isComparisonWith(
	pass *analysis.Pass,
	e ast.Expr,
	predicate func(pass *analysis.Pass, e ast.Expr) bool,
	op token.Token,
) (ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}
	if be.Op != op {
		return nil, false
	}

	t1, t2 := predicate(pass, be.X), predicate(pass, be.Y)
	if xor(t1, t2) {
		if t1 {
			return be.Y, true
		}
		return be.X, true
	}
	return nil, false
}

func isNegation(e ast.Expr) (ast.Expr, bool) {
	ue, ok := e.(*ast.UnaryExpr)
	if !ok {
		return nil, false
	}
	return ue.X, ue.Op == token.NOT
}

func xor(a, b bool) bool {
	return a != b
}

func anyVal[T any](bools []bool, vals ...T) (T, bool) {
	if len(bools) != len(vals) {
		panic("inconsistent usage of valOr")
	}

	for i, b := range bools {
		if b {
			return vals[i], true
		}
	}

	var _default T
	return _default, false
}
