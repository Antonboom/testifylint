package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Contains detects situations like
//
//	assert.True(t, strings.Contains(a, "abc123"))
//	assert.False(t, strings.Contains(a, "456"))
//	assert.True(t, strings.Contains(string(b), "abc123"))
//	assert.False(t, strings.Contains(string(b), "456"))
//
// and requires
//
//	assert.Contains(t, a, "abc123")
//	assert.NotContains(t, a, "456")
//	assert.Contains(t, string(b), "abc123")
//	assert.NotContains(t, string(b), "456")
type Contains struct{}

// NewContains constructs Contains checker.
func NewContains() Contains   { return Contains{} }
func (Contains) Name() string { return "contains" }

func (checker Contains) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.NameFTrimmed {
	case "True":
		if len(call.Args) < 1 {
			return nil
		}

		ce, ok := call.Args[0].(*ast.CallExpr)
		if !ok {
			return nil
		}
		if len(ce.Args) != 2 {
			return nil
		}

		if isStringsContainsCall(pass, ce) {
			const proposed = "Contains"
			return newUseFunctionDiagnostic(checker.Name(), call, proposed,
				newSuggestedFuncReplacement(call, proposed, analysis.TextEdit{
					Pos:     ce.Pos(),
					End:     ce.End(),
					NewText: formatAsCallArgs(pass, ce.Args[0], ce.Args[1]),
				}),
			)
		}

	case "False":
		if len(call.Args) < 1 {
			return nil
		}

		ce, ok := call.Args[0].(*ast.CallExpr)
		if !ok {
			return nil
		}
		if len(ce.Args) != 2 {
			return nil
		}

		if isStringsContainsCall(pass, ce) {
			const proposed = "NotContains"
			return newUseFunctionDiagnostic(checker.Name(), call, proposed,
				newSuggestedFuncReplacement(call, proposed, analysis.TextEdit{
					Pos:     ce.Pos(),
					End:     ce.End(),
					NewText: formatAsCallArgs(pass, ce.Args[0], ce.Args[1]),
				}),
			)
		}
	}
	return nil
}
