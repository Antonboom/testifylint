package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Contains detects situations like
//
//	assert.True(t, strings.Contains(a, "abc123"))
//	assert.False(t, !strings.Contains(a, "abc123"))
//	assert.False(t, strings.Contains(a, "abc123"))
//	assert.True(t, !strings.Contains(a, "abc123"))
//	assert.Contains(t, arr, 1, 2)
//	assert.NotContains(t, arr, 1, 2)
//
// and requires
//
//	assert.Contains(t, a, "abc123")
//	assert.NotContains(t, a, "abc123")
//	assert.Subset(t, arr, 1, 2)
//	assert.NotSubset(t, arr, 1, 2)
type Contains struct{}

// NewContains constructs Contains checker.
func NewContains() Contains   { return Contains{} }
func (Contains) Name() string { return "contains" }

func (checker Contains) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if d := checker.checkStringContains(pass, call); d != nil {
		return d
	}

	return checker.checkSubset(pass, call)
}

func (checker Contains) checkStringContains(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}

	expr := call.Args[0]
	unpacked, isNeg := isNegation(expr)
	if isNeg {
		expr = unpacked
	}

	ce, ok := expr.(*ast.CallExpr)
	if !ok || len(ce.Args) != 2 {
		return nil
	}

	if !isStringsContainsCall(pass, ce) {
		return nil
	}

	var proposed string
	switch call.Fn.NameFTrimmed {
	default:
		return nil

	case "True":
		proposed = "Contains"
		if isNeg {
			proposed = "NotContains"
		}

	case "False":
		proposed = "NotContains"
		if isNeg {
			proposed = "Contains"
		}
	}

	return newUseFunctionDiagnostic(checker.Name(), call, proposed,
		analysis.TextEdit{
			Pos:     call.Args[0].Pos(),
			End:     call.Args[0].End(),
			NewText: formatAsCallArgs(pass, ce.Args[0], ce.Args[1]),
		})
}

func (checker Contains) checkSubset(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if call.Call.Ellipsis.IsValid() { // msgAndArgs... case.
		return nil
	}
	if len(call.Args) < 3 {
		return nil
	}

	third := call.Args[2]
	if hasStringType(pass, third) {
		// Possible false positives because of format string.
		return nil
	}

	var proposed string
	switch call.Fn.NameFTrimmed {
	default:
		return nil

	case "Contains":
		proposed = "Subset"

	case "NotContains":
		proposed = "NotSubset"
	}
	if call.Fn.IsFmt {
		proposed += "f"
	}

	msg := fmt.Sprintf("invalid usage of %s, use %s.%s for multi elements assertion",
		call.String(), call.SelectorXStr, proposed)

	return newDiagnostic(checker.Name(), call, msg)
}
