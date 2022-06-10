package checker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

type SuiteNoExtraAssertCall struct{}

func NewSuiteNoExtraAssertCall() SuiteNoExtraAssertCall {
	return SuiteNoExtraAssertCall{}
}

func (SuiteNoExtraAssertCall) Name() string       { return "suite-no-extra-assert-call" }
func (SuiteNoExtraAssertCall) Priority() int      { return 9 }
func (SuiteNoExtraAssertCall) DisabledByDefault() { return }

func (checker SuiteNoExtraAssertCall) Check(pass *analysis.Pass, call CallMeta) {
	if !call.InsideSuiteMethod {
		return
	}

	ce, ok := call.Selector.X.(*ast.CallExpr)
	if !ok {
		return
	}
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if se.X == nil || !analysisutil.IsSuiteObj(pass, se.X) {
		return
	}
	if se.Sel == nil || se.Sel.Name != "Assert" {
		return
	}

	r.Report(pass, checker.Name(), call, "need to simplify the check", &analysis.SuggestedFix{
		Message: "Remove Assert() call",
		TextEdits: []analysis.TextEdit{{
			Pos:     se.Sel.Pos(),
			End:     ce.End() + 1, // +1 for dot.
			NewText: []byte(""),
		}},
	})
}
