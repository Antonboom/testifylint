package checkers

import (
	"go/ast"
	"go/types"

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
