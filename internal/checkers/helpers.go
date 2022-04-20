package checkers

import "golang.org/x/tools/go/analysis"

func xor(a, b bool) bool {
	return a != b
}

func reportUseFunction(pass *analysis.Pass, meta FnMeta, proposedFn string) {
	if meta.IsFormatFn {
		pass.Reportf(meta.Pos, "use %s.%sf", meta.Pkg, proposedFn)
	} else {
		pass.Reportf(meta.Pos, "use %s.%s", meta.Pkg, proposedFn)
	}
}

func reportNeedSimplifyCheck(pass *analysis.Pass, meta FnMeta) {
	pass.Reportf(meta.Pos, "need to simplify the check")
}
