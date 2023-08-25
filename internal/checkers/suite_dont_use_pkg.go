package checkers

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

// SuiteDontUsePkg detects situation like
//
//	func (s *MySuite) TestSomething() {
//		assert.Equal(s.T(), 42, value)
//	}
//
// and requires
//
//	func (s *MySuite) TestSomething() {
//		s.Equal(42, value)
//	}
type SuiteDontUsePkg struct{}

// NewSuiteDontUsePkg constructs SuiteDontUsePkg checker.
func NewSuiteDontUsePkg() SuiteDontUsePkg { return SuiteDontUsePkg{} }
func (SuiteDontUsePkg) Name() string      { return "suite-dont-use-pkg" }

func (checker SuiteDontUsePkg) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if s := call.SelectorXStr; !(s == testify.AssertPkgName || s == testify.RequirePkgName) {
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
	if se.X == nil || !implementsTestifySuiteIface(pass, se.X) {
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
	case testify.AssertPkgName:
		newSelector = rcv.Name
	case testify.RequirePkgName:
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

func implementsTestifySuiteIface(pass *analysis.Pass, rcv ast.Expr) bool {
	suiteIface := analysisutil.ObjectOf(pass.Pkg, testify.SuitePkgPath, "TestingSuite")
	if suiteIface == nil {
		return false
	}

	return types.Implements(
		pass.TypesInfo.TypeOf(rcv),
		suiteIface.Type().Underlying().(*types.Interface),
	)
}
