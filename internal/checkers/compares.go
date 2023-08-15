package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Compares checks situation like
//
//	assert.True(t, a >= b)
//
// and requires e.g.
//
//	assert.GreaterOrEqual(t, a, b)
type Compares struct{}

// NewCompares constructs Compares checker.
func NewCompares() Compares   { return Compares{} }
func (Compares) Name() string { return "compares" }

func (checker Compares) Check(_ *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}

	be, ok := call.Args[0].(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	var tokenToProposedFn map[token.Token]string

	switch call.Fn.Name {
	case "True", "Truef":
		tokenToProposedFn = tokenToProposedFnInsteadOfTrue
	case "False", "Falsef":
		tokenToProposedFn = tokenToProposedFnInsteadOfFalse
	default:
		return nil
	}

	if proposedFn, ok := tokenToProposedFn[be.Op]; ok {
		a, b := be.X, be.Y
		return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
			newSuggestedFuncReplacement(call, proposedFn, analysis.TextEdit{
				Pos:     be.X.Pos(),
				End:     be.Y.End(),
				NewText: []byte(types.ExprString(a) + ", " + types.ExprString(b)),
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
