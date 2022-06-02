package checker

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

func (ErrorIs) Priority() int {
	return 7
}

func (checker ErrorIs) Check(pass *analysis.Pass, call CallMeta) {
	if len(call.Args) < 2 {
		return
	}
	errArg := call.Args[1]

	switch call.Fn.Name {
	case "Error", "Errorf":
		if isError(pass, errArg) {
			proposed := "ErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", call.SelectorStr, proposed)
			r.Report(pass, checker.Name(), call, msg, newFixViaFnReplacement(call, proposed))
		}

	case "NoError", "NoErrorf":
		if isError(pass, errArg) {
			proposed := "NotErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", call.SelectorStr, proposed)
			r.Report(pass, checker.Name(), call, msg, newFixViaFnReplacement(call, proposed))
		}
	}
}
