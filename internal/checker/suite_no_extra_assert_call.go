package checker

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

type SuiteNoExtraAssertCall struct{}

func NewSuiteNoExtraAssertCall() SuiteNoExtraAssertCall {
	return SuiteNoExtraAssertCall{}
}

func (SuiteNoExtraAssertCall) Name() string            { return "suite-no-extra-assert-call" }
func (SuiteNoExtraAssertCall) Priority() int           { return 9 }
func (SuiteNoExtraAssertCall) DisabledByDefault() bool { return true }

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
	if se.X == nil || !isSuiteObj(pass, se.X) {
		return
	}
	if se.Sel != nil && se.Sel.Name != "Assert" {
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

func isSuiteObj(pass *analysis.Pass, rcv ast.Expr) bool {
	suiteIface := analysisutil.ObjectOf(pass, "github.com/stretchr/testify/suite", "TestingSuite")
	if suiteIface == nil {
		return false
	}

	return types.Implements(
		pass.TypesInfo.TypeOf(rcv),
		suiteIface.Type().Underlying().(*types.Interface),
	)
}
