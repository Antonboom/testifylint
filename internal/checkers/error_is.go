package checkers

import (
	"fmt"
	"golang.org/x/tools/go/analysis"
)

const errorIsCheckerName = "errors-is"

func ErrorIs(pass *analysis.Pass, fn FnMeta) {
	if len(fn.Args) < 3 {
		return
	}

	switch fn.Name {
	case "Error", "Errorf":
		if isError(pass, fn.Args[2]) {
			pass.Report(analysis.Diagnostic{
				Pos:      fn.Pos.Pos(),
				Category: errorIsCheckerName,
				Message:  fmt.Sprintf("invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", fn.Pkg, fn.Pkg, "ErrorIs"),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Replace Error with ErrorIs",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     fn.Pos.Pos(),
							End:     fn.Pos.End(),
							NewText: []byte("ErrorIs"),
						},
					},
				}},
				Related: nil,
			})
			//r.Reportf(pass, fn, "invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", "ErrorIs")
		}

	case "NoError", "NoErrorf":
		if isError(pass, fn.Args[2]) {
			r.Reportf(pass, fn, "invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", "NotErrorIs")
		}
	}
}
