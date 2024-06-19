package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isZero(e ast.Expr) bool { return isIntNumber(e, 0) }

func isOne(e ast.Expr) bool { return isIntNumber(e, 1) }

func isSignedZero(e ast.Expr) bool {
	return isIntNumber(e, 0) || isTypedSignedIntNumber(e, 0)
}

func isSignedNotZero(pass *analysis.Pass, e ast.Expr) bool {
	return !isUnsigned(pass, e) && !isSignedZero(e)
}

func isIntNumber(e ast.Expr, v int) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == fmt.Sprintf("%d", v)
}

func isTypedSignedIntNumber(e ast.Expr, v int) bool {
	ce, ok := e.(*ast.CallExpr)
	if !ok || len(ce.Args) != 1 {
		return false
	}

	fn, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	switch fn.Name {
	case "int", "int8", "int16", "int32", "int64":
		return isIntNumber(ce.Args[0], v)
	}

	return false
}

func isBasicLit(e ast.Expr) bool {
	_, ok := e.(*ast.BasicLit)
	return ok
}

func isIntBasicLit(e ast.Expr) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT
}

func isUntypedConst(pass *analysis.Pass, e ast.Expr) bool {
	return isUnderlying(pass, e, types.IsUntyped)
}

func isTypedConst(pass *analysis.Pass, e ast.Expr) bool {
	tt, ok := pass.TypesInfo.Types[e]
	return ok && tt.IsValue() && tt.Value != nil
}

func isFloat(pass *analysis.Pass, e ast.Expr) bool {
	return isUnderlying(pass, e, types.IsFloat)
}

func isUnsigned(pass *analysis.Pass, e ast.Expr) bool {
	return isUnderlying(pass, e, types.IsUnsigned)
}

func isUnderlying(pass *analysis.Pass, e ast.Expr, flag types.BasicInfo) bool {
	t := pass.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}

	bt, ok := t.Underlying().(*types.Basic)
	return ok && (bt.Info()&flag > 0)
}

func isPointer(pass *analysis.Pass, e ast.Expr) bool {
	_, ok := pass.TypesInfo.TypeOf(e).(*types.Pointer)
	return ok
}
