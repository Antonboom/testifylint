package checkers

/*
import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

func BoolCompare(pass *analysis.Pass, fn CallMeta) {
	switch fn.Name {
	case "Equal", "Equalf":
		if len(fn.Args) < 3 {
			return
		}

		switch lhs, rhs := fn.Args[1], fn.Args[2]; {
		case xor(isUntypedTrue(lhs), isUntypedTrue(rhs)):
			r.ReportUseFunction(pass, fn, "True")

		case xor(isUntypedFalse(lhs), isUntypedFalse(rhs)):
			r.ReportUseFunction(pass, fn, "False")
		}

	case "NotEqual", "NotEqualf":
		if len(fn.Args) < 3 {
			return
		}

		switch lhs, rhs := fn.Args[1], fn.Args[2]; {
		case xor(isUntypedTrue(lhs), isUntypedTrue(rhs)):
			r.ReportUseFunction(pass, fn, "False")

		case xor(isUntypedFalse(lhs), isUntypedFalse(rhs)):
			r.ReportUseFunction(pass, fn, "True")
		}

	case "True", "Truef":
		if len(fn.Args) < 2 {
			return
		}

		switch arg := fn.Args[1]; {
		case isComparisonWithTrue(arg, token.EQL), isComparisonWithFalse(arg, token.NEQ):
			r.Report(pass, fn, "need to simplify the check")

		case isComparisonWithTrue(arg, token.NEQ), isComparisonWithFalse(arg, token.EQL), isNegation(arg):
			r.ReportUseFunction(pass, fn, "False")
		}

	case "False", "Falsef":
		if len(fn.Args) < 2 {
			return
		}

		switch arg := fn.Args[1]; {
		case isComparisonWithTrue(arg, token.EQL), isComparisonWithFalse(arg, token.NEQ):
			r.Report(pass, fn, "need to simplify the check")

		case isComparisonWithTrue(arg, token.NEQ), isComparisonWithFalse(arg, token.EQL), isNegation(arg):
			r.ReportUseFunction(pass, fn, "True")
		}
	}
}

func isUntypedTrue(e ast.Expr) bool {
	val, ok := e.(*ast.Ident)
	return ok && val.Name == "true"
}

func isUntypedFalse(e ast.Expr) bool { // TODO: use types.Universe
	val, ok := e.(*ast.Ident)
	return ok && val.Name == "false"
}

func isComparisonWithTrue(e ast.Expr, op token.Token) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return be.Op == op && xor(isUntypedTrue(be.X), isUntypedTrue(be.Y))
}

func isComparisonWithFalse(e ast.Expr, op token.Token) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return be.Op == op && xor(isUntypedFalse(be.X), isUntypedFalse(be.Y))
}

func isNegation(e ast.Expr) bool {
	ue, ok := e.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	return ue.Op == token.NOT
}
*/
