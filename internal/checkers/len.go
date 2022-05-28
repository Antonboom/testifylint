package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Len struct{}

func NewLen() Len {
	return Len{}
}

func (Len) Name() string {
	return "len"
}

func (checker Len) Check(pass *analysis.Pass, call CallMeta) {
	switch call.Fn.Name {
	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, targetVal, ok := xorLenCall(pass, a, b); ok {
			r.ReportUseFunction(pass, checker.Name(), call, "Len",
				newFixViaFnReplacement(call, "Len", analysis.TextEdit{
					Pos:     a.Pos(),
					End:     b.End(),
					NewText: []byte(types.ExprString(lenArg) + ", " + types.ExprString(targetVal)),
				}),
			)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return
		}
		expr := call.Args[0]

		if lenArg, targetVal, ok := isLenEquality(pass, expr); ok {
			r.ReportUseFunction(pass, checker.Name(), call, "Len",
				newFixViaFnReplacement(call, "Len", analysis.TextEdit{
					Pos:     expr.Pos(),
					End:     expr.End(),
					NewText: []byte(types.ExprString(lenArg) + ", " + types.ExprString(targetVal)),
				}),
			)
		}
	}
}

func xorLenCall(pass *analysis.Pass, a, b ast.Expr) (lenArg ast.Expr, targetVal ast.Expr, ok bool) {
	arg1, ok1 := isLenCall(pass, a)
	arg2, ok2 := isLenCall(pass, b)

	if xor(ok1, ok2) {
		if ok1 {
			return arg1, b, true
		}
		if ok2 {
			return arg2, a, true
		}
	}
	return nil, nil, false
}

var lenObj = types.Universe.Lookup("len")

func isLenCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	if isObj(pass, ce.Fun, lenObj) && len(ce.Args) == 1 {
		return ce.Args[0], true
	}
	return nil, false
}

func isLenEquality(pass *analysis.Pass, e ast.Expr) (ast.Expr, ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, nil, false
	}

	if be.Op != token.EQL {
		return nil, nil, false
	}
	return xorLenCall(pass, be.X, be.Y)
}
