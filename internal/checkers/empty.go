package checkers

/*
type Empty struct{}

func (Empty) Name() string {
	return "empty"
}

func (checker Empty) Check(pass *analysis.Pass, call CallMeta) {
	if invalid := checkEmpty(call); invalid {
		r.ReportUseFunction(pass, checker.Name(), call,"Empty", nil)
	}

	if invalid := checkNotEmpty(call); invalid {
		r.ReportUseFunction(pass, checker.Name(), call,"NotEmpty", nil)
	}
}

func checkEmpty(call CallMeta) bool {
	switch call.Fn.Name {
	case "Len", "Lenf":
		return len(call.Args) >= 3 && isZero(call.Args[2])

	case "Equal", "Equalf":
		return len(fn.Args) >= 3 &&
			(isLenCall(fn.Args[1]) && isZero(fn.Args[2]) || isZero(fn.Args[1]) && isLenCall(fn.Args[2]))

	case "True", "Truef":
		return len(fn.Args) >= 2 && isZeroLenCheck(fn.Args[1])
	}
	return false
}

func checkNotEmpty(fn CallMeta) bool {
	switch fn.Name {
	case "NotEqual", "NotEqualf":
		return len(fn.Args) >= 3 &&
			(isLenCall(fn.Args[1]) && isZero(fn.Args[2]) || isZero(fn.Args[1]) && isLenCall(fn.Args[2]))

	case "Greater", "Greaterf":
		return len(fn.Args) >= 3 && (isLenCall(fn.Args[1]) && isZero(fn.Args[2]))

	case "GreaterOrEqual", "GreaterOrEqualf":
		return len(fn.Args) >= 3 && (isLenCall(fn.Args[1]) && isOne(fn.Args[2]))

	case "True", "Truef":
		return len(fn.Args) >= 2 &&
			(isBinaryExpr(fn.Args[1], isLenCall, token.NEQ, isZero) ||
				isBinaryExpr(fn.Args[1], isZero, token.NEQ, isLenCall) ||
				isBinaryExpr(fn.Args[1], isLenCall, token.GTR, isZero) ||
				isBinaryExpr(fn.Args[1], isZero, token.LSS, isLenCall) ||
				isBinaryExpr(fn.Args[1], isLenCall, token.GEQ, isOne) ||
				isBinaryExpr(fn.Args[1], isOne, token.LEQ, isLenCall))
	}
	return false
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

func isZeroLenCheck(e ast.Expr) bool {
	return isBinaryExpr(e, isLenCall, token.EQL, isZero) ||
		isBinaryExpr(e, isZero, token.EQL, isLenCall)
}

func isBinaryExpr(e ast.Expr, lhs predicate, op token.Token, rhs predicate) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	return (be.Op == op) && lhs(be.X) && rhs(be.Y)
}
*/
