package checkers

import "go/ast"

type predicate func(expr ast.Expr) bool

func xor(a, b bool) bool {
	return a != b
}
