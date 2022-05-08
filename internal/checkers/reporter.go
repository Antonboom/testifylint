package checkers

import (
	"go/token"
	"golang.org/x/tools/go/analysis"
)

var r = newReporter()

type reporter struct {
	cache map[token.Pos]struct{} // TODO: можно добавить приоритеты чекерам
}

func newReporter() *reporter {
	return &reporter{cache: map[token.Pos]struct{}{}}
}

func (r *reporter) Report(pass *analysis.Pass, checker string, meta CallMeta, msg string) {
	r.reportf(pass, checker, meta.Range.Pos(), msg)
}

func (r *reporter) ReportUseFunction(pass *analysis.Pass, checker string, meta CallMeta, proposedFn string) {
	r.Reportf(pass, checker, meta, "use %s.%s", proposedFn)
}

func (r *reporter) Reportf(pass *analysis.Pass, checker string, meta CallMeta, msg string, proposedFn string) {
	f := proposedFn
	if meta.Fn.IsFmt {
		f += "f"
	}
	r.reportf(pass, checker, meta.Range.Pos(), msg, meta.SelectorStr, f)
}

func (r *reporter) reportf(p *analysis.Pass, checker string, pos token.Pos, format string, args ...any) {
	if _, ok := r.cache[pos]; ok {
		return
	}

	p.Reportf(pos, checker+": "+format, args...)
	r.cache[pos] = struct{}{}
}
