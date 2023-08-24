package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// Empty checks situation like
//
//	assert.Equal(t, len(arr), 0)
//
// and requires e.g.
//
//	assert.Empty(t, arr)
type Empty struct{}

// NewEmpty constructs Empty checker.
func NewEmpty() Empty      { return Empty{} }
func (Empty) Name() string { return "empty" }

func (checker Empty) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if d := checker.checkEmpty(pass, call); d != nil {
		return d
	}
	return checker.checkNotEmpty(pass, call)
}

func (checker Empty) checkEmpty(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic { //nolint:gocognit
	newUseEmptyDiagnostic := func(replaceStart, replaceEnd token.Pos, replaceWith ast.Expr) *analysis.Diagnostic {
		return newUseFunctionDiagnostic(checker.Name(), call, "Empty",
			newSuggestedFuncReplacement(call, "Empty", analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: analysisutil.NodeBytes(pass.Fset, replaceWith),
			}),
		)
	}

	switch call.Fn.Name {
	case "Len", "Lenf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if isZero(b) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), a)
		}

	case "Equal", "Equalf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Less", "Lessf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isBuiltinLenCall(pass, a); ok && isOne(b) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Greater", "Greaterf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isBuiltinLenCall(pass, b); ok && isOne(a) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return nil
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) == 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.EQL

		// 0 == len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.EQL

		// len(%s) < 1
		arg3, ok3 := isBuiltinLenCall(pass, a)
		ok3 = ok3 && isOne(b) && op == token.LSS

		// 1 > len(%s)
		arg4, ok4 := isBuiltinLenCall(pass, b)
		ok4 = ok4 && isOne(a) && op == token.GTR

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return nil
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) != 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.NEQ

		// 0 != len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.NEQ

		// len(%s) >= 1
		arg3, ok3 := isBuiltinLenCall(pass, a)
		ok3 = ok3 && isOne(b) && op == token.GEQ

		// 1 <= len(%s)
		arg4, ok4 := isBuiltinLenCall(pass, b)
		ok4 = ok4 && isOne(a) && op == token.LEQ

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}
	}
	return nil
}

func (checker Empty) checkNotEmpty(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic { //nolint:gocognit
	newUseNotEmptyDiagnostic := func(replaceStart, replaceEnd token.Pos, replaceWith ast.Expr) *analysis.Diagnostic {
		return newUseFunctionDiagnostic(checker.Name(), call, "NotEmpty",
			newSuggestedFuncReplacement(call, "NotEmpty", analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: analysisutil.NodeBytes(pass.Fset, replaceWith),
			}),
		)
	}

	switch call.Fn.Name {
	case "NotEqual", "NotEqualf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Greater", "Greaterf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isBuiltinLenCall(pass, a); ok && isZero(b) {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Less", "Lessf":
		if len(call.Args) < 2 {
			return nil
		}
		a, b := call.Args[0], call.Args[1]

		if lenArg, ok := isBuiltinLenCall(pass, b); ok && isZero(a) {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "True", "Truef":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return nil
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) != 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.NEQ

		// 0 != len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.NEQ

		// len(%s) > 0
		arg3, ok3 := isBuiltinLenCall(pass, a)
		ok3 = ok3 && isZero(b) && op == token.GTR

		// 0 < len(%s)
		arg4, ok4 := isBuiltinLenCall(pass, b)
		ok4 = ok4 && isZero(a) && op == token.LSS

		if lenArg, ok := anyVal([]bool{ok1, ok2, ok3, ok4}, arg1, arg2, arg3, arg4); ok {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "False", "Falsef":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		be, ok := expr.(*ast.BinaryExpr)
		if !ok {
			return nil
		}
		a, b, op := be.X, be.Y, be.Op

		// len(%s) == 0
		arg1, ok1 := isLenCallAndZero(pass, a, b)
		ok1 = ok1 && op == token.EQL

		// 0 == len(%s)
		arg2, ok2 := isLenCallAndZero(pass, b, a)
		ok2 = ok2 && op == token.EQL

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}
	}
	return nil
}

var lenObj = types.Universe.Lookup("len")

func isLenCallAndZero(pass *analysis.Pass, a, b ast.Expr) (ast.Expr, bool) {
	lenArg, ok := isBuiltinLenCall(pass, a)
	return lenArg, ok && isZero(b)
}

func isBuiltinLenCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	if analysisutil.IsObj(pass.TypesInfo, ce.Fun, lenObj) && len(ce.Args) == 1 {
		return ce.Args[0], true
	}
	return nil, false
}

func isZero(e ast.Expr) bool {
	return isIntNumber(e, 0)
}

func isOne(e ast.Expr) bool {
	return isIntNumber(e, 1)
}

func isIntNumber(e ast.Expr, v int) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == fmt.Sprintf("%d", v)
}
