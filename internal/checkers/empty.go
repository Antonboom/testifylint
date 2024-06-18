package checkers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// Empty detects situations like
//
//	assert.Len(t, a, 0)
//	assert.Equal(t, 0, len(a))
//	assert.EqualValues(t, 0, len(a))
//	assert.Exactly(t, 0, len(a))
//	assert.LessOrEqual(t, len(a), 0)
//	assert.GreaterOrEqual(t, 0, len(a))
//	assert.Less(t, len(a), 0)
//	assert.Greater(t, 0, len(a))
//	assert.Less(t, len(a), 1)
//	assert.Greater(t, 1, len(a))
//	assert.Zero(t, len(a))
//	assert.Empty(t, len(a))
//	assert.Empty(t, string(a))
//
//	assert.NotEqual(t, 0, len(a))
//	assert.NotEqualValues(t, 0, len(a))
//	assert.Less(t, 0, len(a))
//	assert.Greater(t, len(a), 0)
//	assert.Positive(t, len(a))
//	assert.NotZero(t, len(a))
//	assert.NotEmpty(t, len(a))
//	assert.NotEmpty(t, string(a))
//
// and requires
//
//	assert.Empty(t, a)
//	assert.NotEmpty(t, a)
//
// String conversion (like `assert.Len(t, string(b), 0)`) are also supported.
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
		const proposed = "Empty"
		return newUseFunctionDiagnostic(checker.Name(), call, proposed,
			newSuggestedFuncReplacement(call, proposed, analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: analysisutil.NodeBytes(pass.Fset, replaceWith),
			}),
		)
	}

	if len(call.Args) == 0 {
		return nil
	}

	a := call.Args[0]
	switch call.Fn.NameFTrimmed {
	case "Empty":
		if simplified, ok := isStringConversion(a); ok {
			return newUseEmptyDiagnostic(a.Pos(), a.End(), simplified)
		}
		fallthrough
	case "Zero":
		lenArg, ok := isBuiltinLenCall(pass, a)
		if ok {
			return newUseEmptyDiagnostic(a.Pos(), a.End(), lenArg)
		}
	}

	if len(call.Args) < 2 {
		return nil
	}
	b := call.Args[1]

	switch call.Fn.NameFTrimmed {
	case "Len":
		if isZero(b) {
			if simplified, ok := isStringConversion(a); ok {
				return newUseEmptyDiagnostic(a.Pos(), b.End(), simplified)
			}
			return newUseEmptyDiagnostic(a.Pos(), b.End(), a)
		}

	case "Equal", "EqualValues", "Exactly":
		if isEmptyString(a) {
			if simplified, ok := isStringConversion(b); ok {
				return newUseEmptyDiagnostic(a.Pos(), b.End(), simplified)
			}
			return newUseEmptyDiagnostic(a.Pos(), b.End(), b)
		}

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "LessOrEqual":
		if lenArg, ok := isBuiltinLenCall(pass, a); ok && isZero(b) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "GreaterOrEqual":
		if lenArg, ok := isBuiltinLenCall(pass, b); ok && isZero(a) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Less":
		if lenArg, ok := isBuiltinLenCall(pass, a); ok && (isOne(b) || isZero(b)) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Greater":
		if lenArg, ok := isBuiltinLenCall(pass, b); ok && (isOne(a) || isZero(a)) {
			return newUseEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}
	}
	return nil
}

func (checker Empty) checkNotEmpty(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic { //nolint:gocognit
	newUseNotEmptyDiagnostic := func(replaceStart, replaceEnd token.Pos, replaceWith ast.Expr) *analysis.Diagnostic {
		const proposed = "NotEmpty"
		return newUseFunctionDiagnostic(checker.Name(), call, proposed,
			newSuggestedFuncReplacement(call, proposed, analysis.TextEdit{
				Pos:     replaceStart,
				End:     replaceEnd,
				NewText: analysisutil.NodeBytes(pass.Fset, replaceWith),
			}),
		)
	}

	if len(call.Args) == 0 {
		return nil
	}

	a := call.Args[0]
	switch call.Fn.NameFTrimmed {
	case "NotEmpty":
		if simplified, ok := isStringConversion(a); ok {
			return newUseNotEmptyDiagnostic(a.Pos(), a.End(), simplified)
		}
		fallthrough
	case "NotZero", "Positive":
		lenArg, ok := isBuiltinLenCall(pass, a)
		if ok {
			return newUseNotEmptyDiagnostic(a.Pos(), a.End(), lenArg)
		}
	}

	if len(call.Args) < 2 {
		return nil
	}
	b := call.Args[1]

	switch call.Fn.NameFTrimmed {
	case "NotEqual", "NotEqualValues":
		if isEmptyString(a) {
			if simplified, ok := isStringConversion(b); ok {
				return newUseNotEmptyDiagnostic(a.Pos(), b.End(), simplified)
			}

			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), b)
		}

		arg1, ok1 := isLenCallAndZero(pass, a, b)
		arg2, ok2 := isLenCallAndZero(pass, b, a)

		if lenArg, ok := anyVal([]bool{ok1, ok2}, arg1, arg2); ok {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Less":
		if lenArg, ok := isBuiltinLenCall(pass, b); ok && isZero(a) {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}

	case "Greater":
		if lenArg, ok := isBuiltinLenCall(pass, a); ok && isZero(b) {
			return newUseNotEmptyDiagnostic(a.Pos(), b.End(), lenArg)
		}
	}
	return nil
}
