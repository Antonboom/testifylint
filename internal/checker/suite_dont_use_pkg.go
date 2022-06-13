package checker

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

type SuiteDontUsePkg struct{}

func NewSuiteDontUsePkg() SuiteDontUsePkg {
	return SuiteDontUsePkg{}
}

func (SuiteDontUsePkg) Name() string  { return "suite-dont-use-pkg" }
func (SuiteDontUsePkg) Priority() int { return 12 }

func (checker SuiteDontUsePkg) Check(pass *analysis.Pass, call CallMeta) {
	if !call.InsideSuiteMethod {
		return
	}
	if s := call.SelectorStr; !(s == "assert" || s == "require") {
		return
	}

	args := call.ArgsRaw
	if len(args) < 2 {
		return
	}
	t := args[0]

	ce, ok := t.(*ast.CallExpr)
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
	if se.Sel == nil || se.Sel.Name != "T" {
		return
	}
	rcv, ok := se.X.(*ast.Ident)
	if !ok {
		return
	}

	var newSelector string
	switch call.SelectorStr {
	case "assert":
		newSelector = rcv.Name
	case "require":
		newSelector = rcv.Name + "." + "Require()"
	}

	msg := fmt.Sprintf("use %s.%s", newSelector, call.Fn.Name)
	r.Report(pass, checker.Name(), call, msg, &analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace %s with %s", call.SelectorStr, newSelector),
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
