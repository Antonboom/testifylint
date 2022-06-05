package checker

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type BoolCompare struct{}

func NewBoolCompare() BoolCompare {
	return BoolCompare{}
}

func (BoolCompare) Name() string  { return "bool-compare" }
func (BoolCompare) Priority() int { return 1 }

func (checker BoolCompare) Check(pass *analysis.Pass, call CallMeta) {
	const (
		needSimplifyMsg  = "need to simplify the check"
		simplifyCheckMsg = "Simplify the check"
	)

	switch call.Fn.Name {
	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return
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

			r.ReportUseFunction(pass, checker.Name(), call, "True",
				newFixViaFnReplacement(call, "True", analysis.TextEdit{
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

			r.ReportUseFunction(pass, checker.Name(), call, "False",
				newFixViaFnReplacement(call, "False", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)
		}

	case "NotEqual", "NotEqualf":
		if len(call.Args) < 2 {
			return
		}

		if len(call.Args) < 2 {
			return
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

			r.ReportUseFunction(pass, checker.Name(), call, "False",
				newFixViaFnReplacement(call, "False", analysis.TextEdit{
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

			r.ReportUseFunction(pass, checker.Name(), call, "True",
				newFixViaFnReplacement(call, "True", analysis.TextEdit{
					Pos:     arg1.Pos(),
					End:     arg2.End(),
					NewText: []byte(types.ExprString(survivingArg)),
				}),
			)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return
		}

		expr := call.Args[0]

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.NEQ)

			if survivingArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
				r.Report(pass, checker.Name(), call, needSimplifyMsg,
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
				r.ReportUseFunction(pass, checker.Name(), call, "False",
					newFixViaFnReplacement(call, "False", analysis.TextEdit{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(types.ExprString(survivingArg)),
					}),
				)
			}
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return
		}

		expr := call.Args[0]

		{
			arg1, ok1 := isComparisonWithTrue(pass, expr, token.EQL)
			arg2, ok2 := isComparisonWithFalse(pass, expr, token.NEQ)

			if survivingArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
				r.Report(pass, checker.Name(), call, needSimplifyMsg,
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
				r.ReportUseFunction(pass, checker.Name(), call, "True",
					newFixViaFnReplacement(call, "True", analysis.TextEdit{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(types.ExprString(survivingArg)),
					}),
				)
			}
		}
	}
}

var (
	trueObj  = types.Universe.Lookup("true")
	falseObj = types.Universe.Lookup("false")
)

func isUntypedTrue(pass *analysis.Pass, e ast.Expr) bool {
	return isObj(pass, e, trueObj)
}

func isUntypedFalse(pass *analysis.Pass, e ast.Expr) bool {
	return isObj(pass, e, falseObj)
}

func isObj(pass *analysis.Pass, e ast.Expr, expected types.Object) bool {
	if expected == nil {
		panic("expect obj must be defined")
	}

	id, ok := e.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.ObjectOf(id)
	return obj == expected
}

func isComparisonWithTrue(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) { // TODO: common code
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}
	if be.Op != op {
		return nil, false
	}

	t1, t2 := isUntypedTrue(pass, be.X), isUntypedTrue(pass, be.Y)
	if xor(t1, t2) {
		if t1 {
			return be.Y, true
		}
		return be.X, true
	}
	return nil, false
}

func isComparisonWithFalse(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}
	if be.Op != op {
		return nil, false
	}

	f1, f2 := isUntypedFalse(pass, be.X), isUntypedFalse(pass, be.Y)
	if xor(f1, f2) {
		if f1 {
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
