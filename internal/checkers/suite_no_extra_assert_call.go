package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

const SuiteNoExtraAssertCallCheckerName = "suite-no-extra-assert-call"

// SuiteNoExtraAssertCall wants t.Helper() call in suite helpers:
//
//	func (s *RoomSuite) assertRoomRound(roundID RoundID) {
//		s.T().Helper()
//		s.Equal(roundID, s.getRoom().CurrentRound.ID)
//	}
type SuiteNoExtraAssertCall struct{}                    //
func NewSuiteNoExtraAssertCall() SuiteNoExtraAssertCall { return SuiteNoExtraAssertCall{} }
func (SuiteNoExtraAssertCall) Name() string             { return SuiteNoExtraAssertCallCheckerName }

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
