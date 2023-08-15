package checkers

import (
	util "github.com/Antonboom/testifylint/internal/analysisutil"
	"golang.org/x/tools/go/analysis"
)

// Error checks situation like
//
//	assert.Nil(t, err)
//
// and requires e.g.
//
//	assert.NoError(t, err)
type Error struct{}

// NewError constructs Error checker.
func NewError() Error      { return Error{} }
func (Error) Name() string { return "error" }

func (checker Error) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}
	errArg := call.Args[0]

	switch call.Fn.Name {
	case "NotNil", "NotNilf":
		if util.IsError(pass, errArg) {
			return newUseFunctionDiagnostic(checker.Name(), call, "Error",
				newSuggestedFuncReplacement(call, "Error"))
		}

	case "Nil", "Nilf":
		if util.IsError(pass, errArg) {
			return newUseFunctionDiagnostic(checker.Name(), call, "NoError",
				newSuggestedFuncReplacement(call, "NoError"))
		}
	}
	return nil
}
