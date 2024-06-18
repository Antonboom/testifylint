package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isZero(e ast.Expr) bool { return isIntNumber(e, 0) }

func isNotZero(e ast.Expr) bool { return !isZero(e) }

func isOne(e ast.Expr) bool { return isIntNumber(e, 1) }

func isIntNumber(e ast.Expr, v int) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == fmt.Sprintf("%d", v)
}

func isBasicLit(e ast.Expr) bool {
	_, ok := e.(*ast.BasicLit)
	return ok
}

func isIntBasicLit(e ast.Expr) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT
}

func isUntypedConst(p *analysis.Pass, e ast.Expr) bool {
	t := p.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}

	b, ok := t.(*types.Basic)
	return ok && b.Info()&types.IsUntyped > 0
}

func isTypedConst(p *analysis.Pass, e ast.Expr) bool {
	tt, ok := p.TypesInfo.Types[e]
	return ok && tt.IsValue() && tt.Value != nil
}

func isFloat(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	bt, ok := t.Underlying().(*types.Basic)
	return ok && (bt.Info()&types.IsFloat > 0)
}

func isPointer(pass *analysis.Pass, expr ast.Expr) bool {
	_, ok := pass.TypesInfo.TypeOf(expr).(*types.Pointer)
	return ok
}

func isEmptyString(expr ast.Expr) bool {
	bl, ok := expr.(*ast.BasicLit)
	return ok && bl.Kind == token.STRING && bl.Value == `""`
}

func isStringConversion(e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok || len(ce.Args) != 1 {
		return nil, false
	}

	fn, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return nil, false
	}

	if fn.Name != "string" {
		return nil, false
	}
	return ce.Args[0], true
}
