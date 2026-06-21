package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// callArgsMatch returns true if the given arguments (usually from analysed call) match the given signature.
func callArgsMatch(pass *analysis.Pass, signature *types.Signature, args []ast.Expr, ellipsisCall bool) bool {
	// Log(prefix string, values ...int) vs Log("result", vals...).
	if ellipsisCall {
		return signature.Variadic() && nonVariadicArgsMatch(pass, signature, args)
	}

	// Log(prefix string, values ...int) vs Log("result", 1, 2, 3).
	if signature.Variadic() {
		return variadicArgsMatch(pass, signature, args)
	}

	// CreateUser(ctx context.Context, user User) vs CreateUser(ctx, usr)
	return nonVariadicArgsMatch(pass, signature, args)
}

func nonVariadicArgsMatch(pass *analysis.Pass, signature *types.Signature, args []ast.Expr) bool {
	params := signature.Params()
	if len(args) != params.Len() {
		return false
	}

	for i, arg := range args {
		if !isAssignable(pass, arg, params.At(i).Type()) {
			return false
		}
	}
	return true
}

func variadicArgsMatch(pass *analysis.Pass, signature *types.Signature, args []ast.Expr) bool {
	params := signature.Params()
	fixed := params.Len() - 1
	if len(args) < fixed {
		return false
	}

	// Args, before variadic arg.
	for i, arg := range args[:fixed] {
		if !isAssignable(pass, arg, params.At(i).Type()) {
			return false
		}
	}

	// Variadic arg.
	elem := params.At(fixed).Type().(*types.Slice).Elem()
	for _, arg := range args[fixed:] {
		if !isAssignable(pass, arg, elem) {
			return false
		}
	}
	return true
}

func isAssignable(pass *analysis.Pass, expr ast.Expr, target types.Type) bool {
	typ := pass.TypesInfo.TypeOf(expr)
	return typ != nil && types.AssignableTo(typ, target)
}
