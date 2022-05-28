package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Compares struct{}

func NewCompares() Compares {
	return Compares{}
}

func (Compares) Name() string {
	return "compares"
}

func (checker Compares) Check(pass *analysis.Pass, call CallMeta) {
	if len(call.Args) < 1 {
		return
	}

	be, ok := call.Args[0].(*ast.BinaryExpr)
	if !ok {
		return
	}

	var tokenToProposedFn map[token.Token]string

	switch call.Fn.Name {
	case "True", "Truef":
		tokenToProposedFn = tokenToProposedFnInsteadOfTrue
	case "False", "Falsef":
		tokenToProposedFn = tokenToProposedFnInsteadOfFalse
	default:
		return
	}

	if proposedFn, ok := tokenToProposedFn[be.Op]; ok {
		a, b := be.X, be.Y
		r.ReportUseFunction(pass, checker.Name(), call, proposedFn,
			newFixViaFnReplacement(call, proposedFn, analysis.TextEdit{
				Pos:     be.X.Pos(),
				End:     be.Y.End(),
				NewText: []byte(types.ExprString(a) + ", " + types.ExprString(b)),
			}),
		)
	}
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
