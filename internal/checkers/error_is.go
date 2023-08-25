package checkers

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
)

// ErrorIs detects situations like
//
//	assert.Error(t, err, errSentinel)
//	assert.NoError(t, err, errSentinel)
//
// and requires
//
//	assert.ErrorIs(t, err, errSentinel)
//	assert.NotErrorIs(t, err, errSentinel)
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
		if isError(pass, errArg) {
			const proposed = "ErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.Error, use %[1]s.%[2]s instead", call.SelectorXStr, proposed)
			return newDiagnostic(checker.Name(), call, msg, newSuggestedFuncReplacement(call, proposed))
		}

	case "NoError", "NoErrorf":
		if isError(pass, errArg) {
			const proposed = "NotErrorIs"
			msg := fmt.Sprintf("invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead", call.SelectorXStr, proposed)
			return newDiagnostic(checker.Name(), call, msg, newSuggestedFuncReplacement(call, proposed))
		}
	}
	return nil
}
