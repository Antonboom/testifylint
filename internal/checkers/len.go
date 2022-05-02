package checkers

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

// todo нашёл issue, пока реализовывал этот линтер

func Len(pass *analysis.Pass, fn FnMeta) {
	invalid := func() bool {
		switch fn.Name {
		case "Equal", "Equalf":
			return len(fn.Args) >= 3 && xor(isLenCall(fn.Args[1]), isLenCall(fn.Args[2]))

		case "True", "Truef":
			return len(fn.Args) >= 2 && isLenEquality(fn.Args[1])
		}
		return false
	}()

	if invalid {
		r.ReportUseFunction(pass, fn, "Len")
	}
}

func isLenCall(e ast.Expr) bool {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}

	fn, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return fn.Name == "len" && len(ce.Args) == 1
}

func isLenEquality(e ast.Expr) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return be.Op == token.EQL && xor(isLenCall(be.X), isLenCall(be.Y))
}
