package checkers

import "golang.org/x/tools/go/analysis"

func ErrorIs(pass *analysis.Pass, fn FnMeta) {
	if len(fn.Args) < 3 {
		return
	}

	switch fn.Name {
	case "Error", "Errorf":
		if isError(pass, fn.Args[2]) {
			r.Reportf(pass, fn, "invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", "ErrorIs")
		}

	case "NoError", "NoErrorf":
		if isError(pass, fn.Args[2]) {
			r.Reportf(pass, fn, "invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", "NotErrorIs")
		}
	}
}
