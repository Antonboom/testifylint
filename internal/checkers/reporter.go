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

func (r *reporter) Report(pass *analysis.Pass, meta FnMeta, msg string) {
	r.reportf(pass, meta.Pos, msg)
}

func (r *reporter) Reportf(pass *analysis.Pass, meta FnMeta, msg string, proposedFn string) {
	f := proposedFn
	if meta.IsFormatFn {
		f += "f"
	}
	r.reportf(pass, meta.Pos, msg, meta.Pkg, f)
}

func (r *reporter) ReportUseFunction(pass *analysis.Pass, meta FnMeta, proposedFn string) {
	r.Reportf(pass, meta, "use %s.%s", proposedFn)
}

func (r *reporter) reportf(p *analysis.Pass, pos token.Pos, format string, args ...any) {
	if _, ok := r.cache[pos]; ok {
		return
	}

	p.Reportf(pos, format, args...)
	r.cache[pos] = struct{}{}
}
