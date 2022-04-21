package checkers

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

func Empty(pass *analysis.Pass, fn FnMeta) {
	if invalid := checkEmpty(pass, fn); invalid {
		r.ReportUseFunction(pass, fn, "Empty")
	}

	if invalid := checkNotEmpty(pass, fn); invalid {
		r.ReportUseFunction(pass, fn, "NotEmpty")
	}
}

func checkEmpty(pass *analysis.Pass, fn FnMeta) bool {
	switch fn.Name {
	case "Len", "Lenf":
		return len(fn.Args) >= 3 && isZero(fn.Args[2])

	case "Equal", "Equalf":
		return len(fn.Args) >= 3 &&
			(isLenCall(fn.Args[1]) && isZero(fn.Args[2]) || isZero(fn.Args[1]) && isLenCall(fn.Args[2]))

	case "True", "Truef":
		return len(fn.Args) >= 2 && isZeroLenCheck(fn.Args[1])
	}
	return false
}

func checkNotEmpty(pass *analysis.Pass, fn FnMeta) bool {
	return false
}

func isZero(e ast.Expr) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == "0"
}

func isZeroLenCheck(e ast.Expr) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return (be.Op == token.EQL) &&
		(isLenCall(be.X) && isZero(be.Y) || isZero(be.X) && isLenCall(be.Y))
}
