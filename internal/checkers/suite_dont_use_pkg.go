package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// SuiteDontUsePkg detects situations like
//
//	func (s *MySuite) TestSomething() {
//		assert.Equal(s.T(), 42, value)
//	}
//
// or
//
//	func (s *MySuite) checkSomething(t *testing.T) {
//		assert.Equal(t, 42, value)
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

func (checker SuiteDontUsePkg) CheckUsingPkg(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if !call.IsPkg {
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
	if se.X == nil || !implementsTestifySuite(pass, se.X) {
		return nil
	}
	if se.Sel == nil || se.Sel.Name != "T" {
		return nil
	}
	rcv, ok := se.X.(*ast.Ident) // At this point we ensure that `s.T()` is used as the first argument of assertion.
	if !ok {
		return nil
	}

	newSelector := rcv.Name
	if !call.IsAssert {
		newSelector += "." + "Require()"
	}

	msg := fmt.Sprintf("use %s.%s", newSelector, call.Fn.Name)
	return newDiagnostic(checker.Name(), call, msg, analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace `%s` with `%s`", call.SelectorXStr, newSelector),
		TextEdits: []analysis.TextEdit{
			// Replace package function with suite method.
			{
				Pos:     call.Selector.X.Pos(),
				End:     call.Selector.X.End(),
				NewText: []byte(newSelector),
			},
			// Remove `s.T()`.
			{
				Pos:     t.Pos(),
				End:     args[1].Pos(),
				NewText: []byte(""),
			},
		},
	})
}

func (checker SuiteDontUsePkg) CheckBareTestingT(pass *analysis.Pass, inspector *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	inspector.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		fd := node.(*ast.FuncDecl)
		if !isSuiteMethod(pass, fd) {
			return
		}

		if ident := fd.Name; ident == nil || isSuiteTestMethod(ident.Name) || isSuiteServiceMethod(ident.Name) {
			return
		}

		for i, param := range fd.Type.Params.List {
			var sel *ast.SelectorExpr
			switch t := param.Type.(type) {
			case *ast.SelectorExpr:
				sel = t
			case *ast.StarExpr:
				if s, ok := t.X.(*ast.SelectorExpr); ok {
					sel = s
				}
			}

			if sel == nil {
				continue
			}
			if pkgIdent, ok := sel.X.(*ast.Ident); !ok || pkgIdent.Name != "testing" {
				continue
			}
			if sel.Sel.Name != "T" {
				continue
			}

			textEditEnd := param.End()
			if i < len(fd.Type.Params.List)-1 {
				textEditEnd = fd.Type.Params.List[i+1].Pos()
			}

			msg := "suite method must not include a testing.T parameter"
			d := newDiagnostic(checker.Name(), fd, msg, analysis.SuggestedFix{
				Message: "Remove testing.T parameter",
				TextEdits: []analysis.TextEdit{
					{
						Pos: param.Pos(),
						End: textEditEnd,
					},
				},
			})
			diagnostics = append(diagnostics, *d)
		}
	})
	return diagnostics
}

func (checker SuiteDontUsePkg) Check(pass *analysis.Pass, inspector *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	inspector.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		call := NewCallMeta(pass, node.(*ast.CallExpr))
		if call == nil {
			return
		}
		pkgCall := checker.CheckUsingPkg(pass, call)
		if pkgCall != nil {
			diagnostics = append(diagnostics, *pkgCall)
		}
	})

	bareTestingTCall := checker.CheckBareTestingT(pass, inspector)
	diagnostics = append(diagnostics, bareTestingTCall...)

	return diagnostics
}
