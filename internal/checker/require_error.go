package checker

import "golang.org/x/tools/go/analysis"

type RequireError struct{}

func NewRequireError() RequireError {
	return RequireError{}
}

func (RequireError) Name() string  { return "require-error" }
func (RequireError) Priority() int { return 8 }

func (checker RequireError) Check(pass *analysis.Pass, call CallMeta) {
	switch call.Fn.Name {
	case "Error", "ErrorIs", "ErrorAs", "EqualError", "ErrorContains", "NoError", "NotErrorIs",
		"Errorf", "ErrorIsf", "ErrorAsf", "EqualErrorf", "ErrorContainsf", "NoErrorf", "NotErrorIsf":
	default:
		return
	}

	if call.IsAssert {
		r.Report(pass, checker.Name(), call, "for error assertions use the `require` package", nil)
	}
}
