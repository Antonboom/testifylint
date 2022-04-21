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

func (r *reporter) ReportUseFunction(pass *analysis.Pass, meta FnMeta, proposedFn string) {
	if meta.IsFormatFn {
		r.reportf(pass, meta.Pos, "use %s.%sf", meta.Pkg, proposedFn)
	} else {
		r.reportf(pass, meta.Pos, "use %s.%s", meta.Pkg, proposedFn)
	}
}

func (r *reporter) ReportNeedSimplifyCheck(pass *analysis.Pass, meta FnMeta) {
	r.reportf(pass, meta.Pos, "need to simplify the check")
}

func (r *reporter) reportf(p *analysis.Pass, pos token.Pos, format string, args ...interface{}) {
	if _, ok := r.cache[pos]; ok {
		return
	}

	p.Reportf(pos, format, args...)
	r.cache[pos] = struct{}{}
}
