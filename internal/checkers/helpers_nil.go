package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func xorNil(pass *analysis.Pass, first, second ast.Expr) (ast.Expr, bool) {
	a, b := isNilExpr(pass, first), isNilExpr(pass, second)
	if xor(a, b) {
		if a {
			return second, true
		}
		return first, true
	}
	return nil, false
}

// isNilExpr returns true if expr is the untyped nil identifier or a typed nil
// conversion such as (*int)(nil), []int(nil), (any)(nil), etc.
func isNilExpr(pass *analysis.Pass, expr ast.Expr) bool {
	return isNil(expr) || isTypedNil(pass, expr)
}

func isNil(expr ast.Expr) bool {
	return isIdentWithName("nil", expr)
}

// isTypedNil returns true if expr is a type conversion of nil to a specific type,
// e.g. (*int)(nil), []int(nil), (any)(nil), (chan struct{})(nil), etc.
func isTypedNil(pass *analysis.Pass, expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok || len(ce.Args) != 1 {
		return false
	}
	if !isNil(ce.Args[0]) {
		return false
	}
	// Check that the function part is a type expression (not a regular function call).
	tv, ok := pass.TypesInfo.Types[ce.Fun]
	return ok && tv.IsType()
}
