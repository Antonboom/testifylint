package checker

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var r = newReporter()

type positionID struct {
	pkg string
	pos token.Pos
}

type reporter struct {
	cache map[positionID]struct{}
}

func newReporter() *reporter {
	return &reporter{cache: map[positionID]struct{}{}}
}

func (r *reporter) ReportUseFunction(
	pass *analysis.Pass,
	checker string,
	call CallMeta,
	proposedFn string,
	fix *analysis.SuggestedFix,
) {
	f := proposedFn
	if call.Fn.IsFmt {
		f += "f"
	}
	msg := fmt.Sprintf("use %s.%s", call.SelectorStr, f)

	r.Report(pass, checker, call.Range, msg, fix)
}

func newFixViaFnReplacement(call CallMeta, proposedFn string, additionalEdits ...analysis.TextEdit) *analysis.SuggestedFix {
	if call.Fn.IsFmt {
		proposedFn += "f"
	}
	return &analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace %s with %s", call.Fn.Name, proposedFn),
		TextEdits: append([]analysis.TextEdit{
			{
				Pos:     call.Fn.Pos(),
				End:     call.Fn.End(),
				NewText: []byte(proposedFn),
			},
		}, additionalEdits...),
	}
}

func (r *reporter) Report(pass *analysis.Pass, checker string, rng analysis.Range, msg string, fix *analysis.SuggestedFix) {
	if !rng.Pos().IsValid() || !rng.End().IsValid() {
		panic("invalid report position")
	}

	posID := positionID{pkg: pass.Pkg.String(), pos: rng.Pos()}
	if _, ok := r.cache[posID]; ok {
		return
	}

	d := analysis.Diagnostic{
		Pos:      rng.Pos(),
		End:      rng.End(),
		Category: checker,
		Message:  checker + ": " + msg,
	}
	if fix != nil {
		d.SuggestedFixes = []analysis.SuggestedFix{*fix}
	}

	pass.Report(d)
	r.cache[posID] = struct{}{}
}
