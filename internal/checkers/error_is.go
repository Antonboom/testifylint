package checkers

import (
	"fmt"
	"golang.org/x/tools/go/analysis"
)

type ErrorIs struct{}

func NewErrorIs() ErrorIs {
	return ErrorIs{}
}

func (ErrorIs) Name() string {
	return "error-is"
}

func (checker ErrorIs) Check(pass *analysis.Pass, call CallMeta) {
	if len(call.Args) < 2 {
		return
	}

	switch call.Fn.Name {
	case "Error", "Errorf":
		if isError(pass, call.Args[1]) {
			pass.Report(analysis.Diagnostic{
				Pos:      call.Pos(),
				End:      call.End(),
				Category: checker.Name(),
				Message: fmt.Sprintf(
					"%[1]s: invalid usage of %[2]s.Error, use %[2]s.%[3]s instead",
					checker.Name(),
					call.SelectorStr,
					"ErrorIs",
				),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Replace Error with ErrorIs",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     call.Fn.Pos(),
							End:     call.Fn.End(),
							NewText: []byte("ErrorIs"),
						},
					},
				}},
				Related: nil,
			})
			//r.Reportf(pass, fn, "invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", "ErrorIs")
		}

	case "NoError", "NoErrorf":
		if isError(pass, call.Args[1]) {
			pass.Report(analysis.Diagnostic{
				Pos:      call.Pos(),
				End:      call.End(),
				Category: checker.Name(),
				Message: fmt.Sprintf(
					"%[1]s: invalid usage of %[2]s.NoError, use %[2]s.%[3]s instead",
					checker.Name(),
					call.SelectorStr,
					"NotErrorIs",
				),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Replace NoError with NotErrorIs",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     call.Fn.Pos(),
							End:     call.Fn.End(),
							NewText: []byte("NotErrorIs"),
						},
					},
				}},
				Related: nil,
			})

			//r.Reportf(pass, call, "invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", "NotErrorIs")
		}
	}
}
