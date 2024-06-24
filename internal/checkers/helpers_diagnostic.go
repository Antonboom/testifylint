package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func newUseFunctionDiagnostic(
	checker string,
	call *CallMeta,
	proposedFn string,
	fix *analysis.SuggestedFix,
) *analysis.Diagnostic {
	f := proposedFn
	if call.Fn.IsFmt {
		f += "f"
	}
	msg := fmt.Sprintf("use %s.%s", call.SelectorXStr, f)

	return newDiagnostic(checker, call, msg, fix)
}

func newRemoveSprintfDiagnostic(
	pass *analysis.Pass,
	checker string,
	call *CallMeta,
	fnPos analysis.Range,
	fnArgs []ast.Expr,
) *analysis.Diagnostic {
	return newRemoveFnDiagnostic(pass, checker, call, "fmt.Sprintf", fnPos, fnArgs...)
}

func newRemoveMustCompileDiagnostic(
	pass *analysis.Pass,
	checker string,
	call *CallMeta,
	fnPos analysis.Range,
	fnArg ast.Expr,
) *analysis.Diagnostic {
	return newRemoveFnDiagnostic(pass, checker, call, "regexp.MustCompile", fnPos, fnArg)
}

func newRemoveFnDiagnostic(
	pass *analysis.Pass,
	checker string,
	call *CallMeta,
	fnName string,
	fnPos analysis.Range,
	fnArgs ...ast.Expr,
) *analysis.Diagnostic {
	return newDiagnostic(checker, call, "remove unnecessary "+fnName, &analysis.SuggestedFix{
		Message: fmt.Sprintf("Remove `%s`", fnName),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     fnPos.Pos(),
				End:     fnPos.End(),
				NewText: formatAsCallArgs(pass, fnArgs...),
			},
		},
	})
}

func newDiagnostic(
	checker string,
	rng analysis.Range,
	msg string,
	fix *analysis.SuggestedFix,
) *analysis.Diagnostic {
	var suggestedFixes []analysis.SuggestedFix
	if fix != nil {
		suggestedFixes = append(suggestedFixes, *fix)
	}
	return newAnalysisDiagnostic(checker, rng, msg, suggestedFixes)
}

func newAnalysisDiagnostic(
	checker string,
	rng analysis.Range,
	msg string,
	suggestedFixes []analysis.SuggestedFix,
) *analysis.Diagnostic {
	return &analysis.Diagnostic{
		Pos:            rng.Pos(),
		End:            rng.End(),
		Category:       checker,
		Message:        checker + ": " + msg,
		SuggestedFixes: suggestedFixes,
	}
}

func newSuggestedFuncReplacement(
	call *CallMeta,
	proposedFn string,
	additionalEdits ...analysis.TextEdit,
) *analysis.SuggestedFix {
	if call.Fn.IsFmt {
		proposedFn += "f"
	}
	return &analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace `%s` with `%s`", call.Fn.Name, proposedFn),
		TextEdits: append([]analysis.TextEdit{
			{
				Pos:     call.Fn.Pos(),
				End:     call.Fn.End(),
				NewText: []byte(proposedFn),
			},
		}, additionalEdits...),
	}
}
