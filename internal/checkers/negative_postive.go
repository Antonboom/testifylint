package checkers

import (
	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// NegativePositive detects situations like
//
//	assert.Greater(t, 0, value)
//	assert.Less(t, 0, value)
//
// and requires
//
//	assert.Positive(t, value)
//	assert.Negative(t, value)
type NegativePositive struct{}

// NewNegativePositive constructs NegativePositive checker.
func NewNegativePositive() NegativePositive { return NegativePositive{} }
func (NegativePositive) Name() string       { return "negative-positive" }

func (checker NegativePositive) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}

	secondArgZero := isZero(call.Args[1])
	if !xor(isZero(call.Args[0]), secondArgZero) {
		return nil
	}

	var proposedFn string
	switch call.Fn.NameFTrimmed {
	case "Greater":
		proposedFn = "Positive"
		if secondArgZero {
			proposedFn = "Negative"
		}
	case "Less":
		proposedFn = "Negative"
		if secondArgZero {
			proposedFn = "Positive"
		}
	default:
		return nil
	}

	survivingArg := call.Args[1]
	if secondArgZero {
		survivingArg = call.Args[0]
	}

	return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
		newSuggestedFuncReplacement(call, proposedFn, analysis.TextEdit{
			Pos:     call.Args[0].Pos(),
			End:     call.Args[1].End(),
			NewText: analysisutil.NodeBytes(pass.Fset, survivingArg),
		}),
	)
}
