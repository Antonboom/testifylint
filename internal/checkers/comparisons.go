package checkers

type Comparisons struct{}

/*
import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

func Comparisons(pass *analysis.Pass, fn CallMeta) {
	if len(fn.Args) < 2 {
		return
	}

	be, ok := fn.Args[1].(*ast.BinaryExpr)
	if !ok {
		return
	}

	switch fn.Name {
	case "True", "Truef":
		if proposed, ok := tokenToProposedFnInsteadOfTrue[be.Op]; ok {
			r.ReportUseFunction(pass, fn, proposed)
		}

	case "False", "Falsef":
		if proposed, ok := tokenToProposedFnInsteadOfFalse[be.Op]; ok {
			r.ReportUseFunction(pass, fn, proposed)
		}
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
*/
