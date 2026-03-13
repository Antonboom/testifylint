package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

var (
	errorObj   = types.Universe.Lookup("error")
	errorType  = errorObj.Type()
	errorIface = errorType.Underlying().(*types.Interface)
)

func isError(pass *analysis.Pass, expr ast.Expr) bool {
	return pass.TypesInfo.TypeOf(expr) == errorType
}

func isErrorsIsCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "errors", "Is")
}

func isErrorsAsCall(pass *analysis.Pass, ce *ast.CallExpr) bool {
	return isPkgFnCall(pass, ce, "errors", "As")
}

// isErrErrorCall returns the receiver expression if e is a method call of the form
// `receiver.Error()` where receiver implements the error interface.
func isErrErrorCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok || len(ce.Args) != 0 {
		return nil, false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || !isIdentWithName("Error", se.Sel) {
		return nil, false
	}

	if !implementsErrorIface(pass, se.X) {
		return nil, false
	}
	return se.X, true
}

// errArgFromErrErrorCallReceiver returns source bytes for an expression that can
// be passed as an `error` argument for receiver.Error() call replacement.
// It returns receiver itself when it implements error, or `&receiver` when only
// pointer type implements error and receiver is addressable.
func errArgFromErrErrorCallReceiver(pass *analysis.Pass, receiver ast.Expr) ([]byte, bool) {
	t := pass.TypesInfo.TypeOf(receiver)
	if t == nil {
		return nil, false
	}

	if types.Implements(t, errorIface) {
		return analysisutil.NodeBytes(pass.Fset, receiver), true
	}

	if !types.Implements(types.NewPointer(t), errorIface) {
		return nil, false
	}

	tv, ok := pass.TypesInfo.Types[receiver]
	if !ok || !tv.Addressable() {
		return nil, false
	}

	return append([]byte("&"), analysisutil.NodeBytes(pass.Fset, receiver)...), true
}

func isAssignableToString(pass *analysis.Pass, e ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(e)
	return t != nil && types.AssignableTo(t, types.Typ[types.String])
}

// implementsErrorIface returns true if the expression's type implements the error interface.
func implementsErrorIface(pass *analysis.Pass, e ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}
	return types.Implements(t, errorIface) || types.Implements(types.NewPointer(t), errorIface)
}
