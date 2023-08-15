package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// SuiteDontUsePkg checks situation like
//
//	func (s *MySuite) TestSomething() {
//		assert.Equal(s.T(), 42, value)
//	}
//
// and requires e.g.
//
//	func (s *MySuite) TestSomething() {
//		s.Equal(42, value)
//	}
type SuiteDontUsePkg struct{}

// NewSuiteDontUsePkg constructs SuiteDontUsePkg checker.
func NewSuiteDontUsePkg() SuiteDontUsePkg { return SuiteDontUsePkg{} }
func (SuiteDontUsePkg) Name() string      { return "suite-dont-use-pkg" }

func (checker SuiteDontUsePkg) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if !call.InsideSuiteMethod {
		return nil
	}
	if s := call.SelectorXStr; !(s == "assert" || s == "require") {
		return nil
	}

	args := call.ArgsRaw
	if len(args) < 2 {
		return nil
	}
	t := args[0]

	ce, ok := t.(*ast.CallExpr)
	if !ok {
		return nil
	}
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	if se.X == nil || !analysisutil.IsTestifySuiteObj(pass, se.X) {
		return nil
	}
	if se.Sel == nil || se.Sel.Name != "T" {
		return nil
	}
	rcv, ok := se.X.(*ast.Ident)
	if !ok {
		return nil
	}

	var newSelector string
	switch call.SelectorXStr {
	case "assert":
		newSelector = rcv.Name
	case "require":
		newSelector = rcv.Name + "." + "Require()"
	}

	msg := fmt.Sprintf("use %s.%s", newSelector, call.Fn.Name)
	return newDiagnostic(checker.Name(), call, msg, &analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace %s with %s", call.SelectorXStr, newSelector),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     call.Selector.Pos(),
				End:     call.Selector.End(),
				NewText: []byte(newSelector + "." + call.Fn.Name),
			},
			{
				Pos:     t.Pos(),
				End:     args[1].Pos(),
				NewText: []byte(""),
			},
		},
	})
}
