package checkers

import (
	util "github.com/Antonboom/testifylint/internal/analysisutil"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
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
		t1, t2 := util.IsUntypedTrue(pass, arg1), util.IsUntypedTrue(pass, arg2)
		f1, f2 := util.IsUntypedFalse(pass, arg1), util.IsUntypedFalse(pass, arg2)

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
		t1, t2 := util.IsUntypedTrue(pass, arg1), util.IsUntypedTrue(pass, arg2)
		f1, f2 := util.IsUntypedFalse(pass, arg1), util.IsUntypedFalse(pass, arg2)

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
			arg1, ok1 := util.IsComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := util.IsComparisonWithFalse(pass, expr, token.NEQ)

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
			arg1, ok1 := util.IsComparisonWithTrue(pass, expr, token.NEQ)
			arg2, ok2 := util.IsComparisonWithFalse(pass, expr, token.EQL)
			arg3, ok3 := util.IsNegation(expr)

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
			arg1, ok1 := util.IsComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := util.IsComparisonWithFalse(pass, expr, token.NEQ)

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
			arg1, ok1 := util.IsComparisonWithTrue(pass, expr, token.NEQ)
			arg2, ok2 := util.IsComparisonWithFalse(pass, expr, token.EQL)
			arg3, ok3 := util.IsNegation(expr)

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
