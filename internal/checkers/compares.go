package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Compares detects situations like
//
//	assert.True(t, a == b)
//	assert.True(t, a != b)
//	assert.True(t, a > b)
//	assert.True(t, a >= b)
//	assert.True(t, a < b)
//	assert.True(t, a <= b)
//	assert.False(t, a == b)
//	...
//
// and requires
//
//	assert.Equal(t, a, b)
//	assert.NotEqual(t, a, b)
//	assert.Greater(t, a, b)
//	assert.GreaterOrEqual(t, a, b)
//	assert.Less(t, a, b)
//	assert.LessOrEqual(t, a, b)
//
// If `a` and `b` are pointers then `assert.Same`/`NotSame` is required instead,
// due to the inappropriate recursive nature of `assert.Equal` (based on `reflect.DeepEqual`).
type Compares struct{}

// NewCompares constructs Compares checker.
func NewCompares() Compares   { return Compares{} }
func (Compares) Name() string { return "compares" }

func (checker Compares) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}

	be, ok := call.Args[0].(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	_, ptrComparison := pass.TypesInfo.TypeOf(be.X).(*types.Pointer)

	var tokenToProposedFn map[token.Token]string

	switch call.Fn.NameFTrimmed {
	case "True":
		if ptrComparison {
			tokenToProposedFn = tokenToProposedFnInsteadOfTrueForPtr
		} else {
			tokenToProposedFn = tokenToProposedFnInsteadOfTrue
		}

	case "False":
		if ptrComparison {
			tokenToProposedFn = tokenToProposedFnInsteadOfFalseForPtr
		} else {
			tokenToProposedFn = tokenToProposedFnInsteadOfFalse
		}
	default:
		return nil
	}

	if proposedFn, ok := tokenToProposedFn[be.Op]; ok {
		a, b := be.X, be.Y
		return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
			newSuggestedFuncReplacement(call, proposedFn, analysis.TextEdit{
				Pos:     be.X.Pos(),
				End:     be.Y.End(),
				NewText: formatAsCallArgs(pass, a, b),
			}),
		)
	}
	return nil
}

var tokenToProposedFnInsteadOfTrue = map[token.Token]string{
	token.EQL: "Equal",
	token.NEQ: "NotEqual",
	token.GTR: "Greater",
	token.GEQ: "GreaterOrEqual",
	token.LSS: "Less",
	token.LEQ: "LessOrEqual",
}

var tokenToProposedFnInsteadOfFalse = map[token.Token]string{
	token.EQL: "NotEqual",
	token.NEQ: "Equal",
	token.GTR: "LessOrEqual",
	token.GEQ: "Less",
	token.LSS: "GreaterOrEqual",
	token.LEQ: "Greater",
}

var tokenToProposedFnInsteadOfTrueForPtr = map[token.Token]string{
	token.EQL: "Same",
	token.NEQ: "NotSame",
}

var tokenToProposedFnInsteadOfFalseForPtr = map[token.Token]string{
	token.EQL: "NotSame",
	token.NEQ: "Same",
}
