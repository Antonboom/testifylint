package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Empty struct{}

func NewEmpty() Empty {
	return Empty{}
}

func (Empty) Name() string {
	return "empty"
}

func (checker Empty) Check(pass *analysis.Pass, call CallMeta) {
	checker.checkEmpty(pass, call)
	checker.checkNotEmpty(pass, call)
}

func (checker Empty) checkEmpty(pass *analysis.Pass, call CallMeta) {
	reportUseEmpty := func(replaceStart, replaceEnd token.Pos, replaceWith ast.Expr) {
		r.ReportUseFunction(pass, checker.Name(), call, "Empty",
			newFixViaFnReplacement(call, "Empty", analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: []byte(types.ExprString(replaceWith)),
			}),
		)
	}

	switch call.Fn.Name {
	case "Len", "Lenf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if isZero(b) {
			reportUseEmpty(a.Pos(), b.End(), a)
		}

	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			reportUseEmpty(a.Pos(), b.End(), lenArg)
		}

	case "Less", "Lessf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isLenCall(pass, a); ok && isOne(b) {
			reportUseEmpty(a.Pos(), b.End(), lenArg)
		}

	case "Greater", "Greaterf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isLenCall(pass, b); ok && isOne(a) {
			reportUseEmpty(a.Pos(), b.End(), lenArg)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) == 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.EQL

		// 0 == len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.EQL

		// len(%s) < 1
		arg3, ok3 := isLenCall(pass, a)
		ok3 = ok3 && isOne(b) && op == token.LSS

		// 1 > len(%s)
		arg4, ok4 := isLenCall(pass, b)
		ok4 = ok4 && isOne(a) && op == token.GTR

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			reportUseEmpty(a.Pos(), b.End(), lenArg)
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) != 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.NEQ

		// 0 != len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.NEQ

		// len(%s) >= 1
		arg3, ok3 := isLenCall(pass, a)
		ok3 = ok3 && isOne(b) && op == token.GEQ

		// 1 <= len(%s)
		arg4, ok4 := isLenCall(pass, b)
		ok4 = ok4 && isOne(a) && op == token.LEQ

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			reportUseEmpty(a.Pos(), b.End(), lenArg)
		}
	}
}

func (checker Empty) checkNotEmpty(pass *analysis.Pass, call CallMeta) {
	reportUseNotEmpty := func(replaceStart, replaceEnd token.Pos, replaceWith ast.Expr) {
		r.ReportUseFunction(pass, checker.Name(), call, "NotEmpty",
			newFixViaFnReplacement(call, "NotEmpty", analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: []byte(types.ExprString(replaceWith)),
			}),
		)
	}

	switch call.Fn.Name {
	case "NotEqual", "NotEqualf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			reportUseNotEmpty(a.Pos(), b.End(), lenArg)
		}

	case "Greater", "Greaterf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isLenCall(pass, a); ok && isZero(b) {
			reportUseNotEmpty(a.Pos(), b.End(), lenArg)
		}

	case "Less", "Lessf":
		if len(call.Args) < 2 {
			return
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isLenCall(pass, b); ok && isZero(a) {
			reportUseNotEmpty(a.Pos(), b.End(), lenArg)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) != 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.NEQ

		// 0 != len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.NEQ

		// len(%s) > 0
		arg3, ok3 := isLenCall(pass, a)
		ok3 = ok3 && isZero(b) && op == token.GTR

		// 0 < len(%s)
		arg4, ok4 := isLenCall(pass, b)
		ok4 = ok4 && isZero(a) && op == token.LSS

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			reportUseNotEmpty(a.Pos(), b.End(), lenArg)
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) == 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.EQL

		// 0 == len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.EQL

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			reportUseNotEmpty(a.Pos(), b.End(), lenArg)
		}
	}
}

func isZero(e ast.Expr) bool {
	return isIntNumber(e, "0")
}

func isOne(e ast.Expr) bool {
	return isIntNumber(e, "1")
}

func isIntNumber(e ast.Expr, v string) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == v
}

func isLenCallAndZero(pass *analysis.Pass, a, b ast.Expr) (ast.Expr, bool) {
	lenArg, ok := isLenCall(pass, a)
	return lenArg, ok && isZero(b)
}
