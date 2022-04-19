package checkers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

func Len(pass *analysis.Pass, fn FnMeta) {
	invalid := func() bool {
		switch fn.Name {
		case "Equal", "Equalf":
			return len(fn.Args) >= 3 && xor(isLenCall(fn.Args[1]), isLenCall(fn.Args[2]))

		case "True", "Truef":
			return len(fn.Args) >= 2 && isLenComparison(fn.Args[1])
		}
		return false
	}()

	if invalid {
		if fn.IsFormatFn {
			pass.Reportf(fn.Pos, "use %s.Lenf", fn.Pkg)
		} else {
			pass.Reportf(fn.Pos, "use %s.Len", fn.Pkg)
		}
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

func isLenComparison(e ast.Expr) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return xor(isLenCall(be.X), isLenCall(be.Y))
}
