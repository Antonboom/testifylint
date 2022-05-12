package checkers

import "go/ast"

type predicate func(expr ast.Expr) bool

func xor(a, b bool) bool {
	return a != b
}

func anyVal[T any](bools []bool, vals ...T) (T, bool) {
	if len(bools) != len(vals) {
		panic("inconsistent usage of valOr")
	}

	for i, b := range bools {
		if b {
			return vals[i], true
		}
	}

	var _default T
	return _default, false
}
