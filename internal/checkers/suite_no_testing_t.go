package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// SuiteNoTestingT requires no testing.T parameter in suite helpers:
//
//	func (s *RoomSuite) assertRoomRound(t *testing.T) {}
type SuiteNoTestingT struct{}

// NewSuiteNoTestingT constructs SuiteNoTestingT checker.
func NewSuiteNoTestingT() SuiteNoTestingT { return SuiteNoTestingT{} }
func (SuiteNoTestingT) Name() string      { return "suite-no-testing-t" }

func (checker SuiteNoTestingT) Check(pass *analysis.Pass, inspector *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
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
