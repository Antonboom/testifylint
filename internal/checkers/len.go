package checkers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Len detects situations like
//
//	assert.Equal(t, 3, len(arr))
//	assert.True(t, len(arr) == 5)
//
// and requires
//
//	assert.Len(t, arr, 3)
//	assert.Len(t, arr, 5)
type Len struct{}

// NewLen constructs Len checker.
func NewLen() Len        { return Len{} }
func (Len) Name() string { return "len" }

func (checker Len) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.Name {
	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, targetVal, ok := xorLenCall(pass, a, b); ok {
			return newUseFunctionDiagnostic(checker.Name(), call, "Len",
				newSuggestedFuncReplacement(call, "Len", analysis.TextEdit{
					Pos:     a.Pos(),
					End:     b.End(),
					NewText: asCallArguments(pass, lenArg, targetVal),
				}),
			)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		if lenArg, targetVal, ok := isLenEquality(pass, expr); ok {
			return newUseFunctionDiagnostic(checker.Name(), call, "Len",
				newSuggestedFuncReplacement(call, "Len", analysis.TextEdit{
					Pos:     expr.Pos(),
					End:     expr.End(),
					NewText: asCallArguments(pass, lenArg, targetVal),
				}),
			)
		}
	}
	return nil
}

func xorLenCall(pass *analysis.Pass, a, b ast.Expr) (lenArg ast.Expr, targetVal ast.Expr, ok bool) {
	arg1, ok1 := isBuiltinLenCall(pass, a)
	arg2, ok2 := isBuiltinLenCall(pass, b)

	if xor(ok1, ok2) {
		if ok1 {
			return arg1, b, true
		}
		if ok2 {
			return arg2, a, true
		}
	}
	return nil, nil, false
}

func isLenEquality(pass *analysis.Pass, e ast.Expr) (ast.Expr, ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, nil, false
	}

	if be.Op != token.EQL {
		return nil, nil, false
	}
	return xorLenCall(pass, be.X, be.Y)
}
