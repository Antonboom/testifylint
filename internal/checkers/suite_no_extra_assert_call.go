package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// SuiteNoExtraAssertCall wants t.Helper() call in suite helpers:
//
//	func (s *RoomSuite) assertRoomRound(roundID RoundID) {
//		s.T().Helper()
//		s.Equal(roundID, s.getRoom().CurrentRound.ID)
//	}
type SuiteNoExtraAssertCall struct{}

// NewSuiteNoExtraAssertCall constructs SuiteNoExtraAssertCall checker.
func NewSuiteNoExtraAssertCall() SuiteNoExtraAssertCall { return SuiteNoExtraAssertCall{} }
func (SuiteNoExtraAssertCall) Name() string             { return "suite-no-extra-assert-call" }

func (checker SuiteNoExtraAssertCall) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	ce, ok := call.Selector.X.(*ast.CallExpr)
	if !ok {
		return nil
	}
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	if se.X == nil || !implementsTestifySuiteIface(pass, se.X) {
		return nil
	}
	if se.Sel == nil || se.Sel.Name != "Assert" {
		return nil
	}

	return newDiagnostic(checker.Name(), call, "need to simplify the check", &analysis.SuggestedFix{
		Message: "Remove Assert() call",
		TextEdits: []analysis.TextEdit{{
			Pos:     se.Sel.Pos(),
			End:     ce.End() + 1, // +1 for dot.
			NewText: []byte(""),
		}},
	})
}
