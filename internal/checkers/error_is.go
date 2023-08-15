package checkers

import (
	"fmt"
	util "github.com/Antonboom/testifylint/internal/analysisutil"

	"golang.org/x/tools/go/analysis"
)

// ErrorIs checks situation like
//
//	assert.Equal(t, len(arr), 3)
//
// and requires e.g.
//
//	assert.Len(t, arr, 3)
type ErrorIs struct{}

// NewErrorIs constructs ErrorIs checker.
func NewErrorIs() ErrorIs    { return ErrorIs{} }
func (ErrorIs) Name() string { return "error-is" }

func (checker ErrorIs) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}
	errArg := call.Args[1]

	switch call.Fn.Name {
	case "Error", "Errorf":
		if util.IsError(pass, errArg) {
			proposed := "ErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", call.SelectorXStr, proposed)
			return newDiagnostic(checker.Name(), call, msg, newSuggestedFuncReplacement(call, proposed))
		}

	case "NoError", "NoErrorf":
		if util.IsError(pass, errArg) {
			proposed := "NotErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", call.SelectorXStr, proposed)
			return newDiagnostic(checker.Name(), call, msg, newSuggestedFuncReplacement(call, proposed))
		}
	}
	return nil
}
