package checkers

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// ErrorCompare detects situations like
//
//	assert.Contains(t, err.Error(), "user not found")
//	assert.Contains(t, err.Error(), errSentinel.Error())
//	assert.Equal(t, err.Error(), "user not found")
//	assert.Equal(t, "user not found", err.Error())
//	assert.Equal(t, err, errSentinel)
//	assert.NotEqual(t, err, errSentinel)
//
// and requires
//
//	assert.ErrorContains(t, err, "user not found")
//	assert.ErrorIs(t, err, errSentinel)
//	assert.EqualError(t, err, "user not found")
//	assert.EqualError(t, err, "user not found")
//	assert.ErrorIs(t, err, errSentinel)
//	assert.NotErrorIs(t, err, errSentinel)
type ErrorCompare struct{}

// NewErrorCompare constructs ErrorCompare checker.
func NewErrorCompare() ErrorCompare { return ErrorCompare{} }
func (ErrorCompare) Name() string   { return "error-compare" }

func (checker ErrorCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.NameFTrimmed {
	case "Contains":
		return checker.checkContains(pass, call)
	case "Equal":
		return checker.checkEqual(pass, call)
	case "NotEqual":
		return checker.checkNotEqual(pass, call)
	}
	return nil
}

func (checker ErrorCompare) checkContains(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}

	errReceiver, ok := isErrErrorCall(pass, call.Args[0])
	if !ok {
		return nil
	}

	errArg, ok := errArgFromErrErrorCallReceiver(pass, errReceiver)
	if !ok {
		return nil
	}

	// assert.Contains(t, err.Error(), sentinel.Error()) → assert.ErrorIs(t, err, sentinel)
	if sentinelReceiver, ok := isErrErrorCall(pass, call.Args[1]); ok {
		sentinelArg, ok := errArgFromErrErrorCallReceiver(pass, sentinelReceiver)
		if !ok {
			return nil
		}
		return newUseFunctionDiagnostic(checker.Name(), call, "ErrorIs",
			analysis.TextEdit{
				Pos:     call.Args[0].Pos(),
				End:     call.Args[0].End(),
				NewText: errArg,
			},
			analysis.TextEdit{
				Pos:     call.Args[1].Pos(),
				End:     call.Args[1].End(),
				NewText: sentinelArg,
			},
		)
	}

	if !isAssignableToString(pass, call.Args[1]) {
		return nil
	}

	// assert.Contains(t, err.Error(), "substr") → assert.ErrorContains(t, err, "substr")
	return newUseFunctionDiagnostic(checker.Name(), call, "ErrorContains",
		analysis.TextEdit{
			Pos:     call.Args[0].Pos(),
			End:     call.Args[0].End(),
			NewText: errArg,
		},
	)
}

func (checker ErrorCompare) checkEqual(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}

	a, b := call.Args[0], call.Args[1]

	// assert.Equal(t, err.Error(), "expected") → assert.EqualError(t, err, "expected")
	if errReceiver, ok := isErrErrorCall(pass, a); ok {
		if !isAssignableToString(pass, b) {
			return nil
		}

		errArg, ok := errArgFromErrErrorCallReceiver(pass, errReceiver)
		if !ok {
			return nil
		}

		return newUseFunctionDiagnostic(checker.Name(), call, "EqualError",
			analysis.TextEdit{
				Pos:     a.Pos(),
				End:     a.End(),
				NewText: errArg,
			},
		)
	}

	// assert.Equal(t, "expected", err.Error()) → assert.EqualError(t, err, "expected")
	if errReceiver, ok := isErrErrorCall(pass, b); ok {
		if !isAssignableToString(pass, a) {
			return nil
		}

		errArg, ok := errArgFromErrErrorCallReceiver(pass, errReceiver)
		if !ok {
			return nil
		}

		return newUseFunctionDiagnostic(checker.Name(), call, "EqualError",
			analysis.TextEdit{
				Pos:     a.Pos(),
				End:     b.End(),
				NewText: append(append(errArg, []byte(", ")...), analysisutil.NodeBytes(pass.Fset, a)...),
			},
		)
	}

	// assert.Equal(t, err, errSentinel) → assert.ErrorIs(t, err, errSentinel)
	// (only when both arguments implement the error interface and neither is nil)
	// No autofix: argument order cannot be determined automatically (ErrorIs is not symmetric).
	if !isNil(a) && !isNil(b) && implementsErrorIface(pass, a) && implementsErrorIface(pass, b) {
		msg := fmt.Sprintf("use %s.ErrorIs", call.SelectorXStr)
		return newDiagnostic(checker.Name(), call, msg)
	}

	return nil
}

func (checker ErrorCompare) checkNotEqual(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}

	a, b := call.Args[0], call.Args[1]

	// assert.NotEqual(t, err, errSentinel) → assert.NotErrorIs(t, err, errSentinel)
	// No autofix: argument order cannot be determined automatically (NotErrorIs is not symmetric).
	if !isNil(a) && !isNil(b) && implementsErrorIface(pass, a) && implementsErrorIface(pass, b) {
		msg := fmt.Sprintf("use %s.NotErrorIs", call.SelectorXStr)
		return newDiagnostic(checker.Name(), call, msg)
	}

	return nil
}
