package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// ErrorNil detects situations like
//
//	assert.Nil(t, err)
//	assert.NotNil(t, err)
//
// and requires
//
//	assert.NoError(t, err)
//	assert.ErrorNil(t, err)
type ErrorNil struct{}

// NewErrorNil constructs ErrorNil checker.
func NewErrorNil() ErrorNil   { return ErrorNil{} }
func (ErrorNil) Name() string { return "error-nil" }

func (checker ErrorNil) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}
	errArg := call.Args[0]

	switch call.Fn.Name {
	case "NotNil", "NotNilf":
		if isError(pass, errArg) {
			return newUseFunctionDiagnostic(checker.Name(), call, "Error",
				newSuggestedFuncReplacement(call, "Error"))
		}

	case "Nil", "Nilf":
		if isError(pass, errArg) {
			return newUseFunctionDiagnostic(checker.Name(), call, "NoError",
				newSuggestedFuncReplacement(call, "NoError"))
		}
	}
	return nil
}

var errIface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isError(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	_, ok := t.Underlying().(*types.Interface)
	return ok && types.Implements(t, errIface)
}
