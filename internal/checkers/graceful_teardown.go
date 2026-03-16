package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/testify"
)

const gracefulTeardownReport = "do not use require in cleanup code"

// GracefulTeardown detects usage of require package's functions in t.Cleanup functions and suite teardown methods.
//
//	func (s *ServiceIntegrationSuite) TearDownTest() {
//		if p := s.verdictsProducer; p != nil {
//			s.Require().NoError(p.Close()) // ← not ok, use s.Assert() instead
//		}
//		if c := s.dlqVerdictsConsumer; c != nil {
//			s.NoError(c.Close()) // ← ok
//		}
//	}
type GracefulTeardown struct{}

// NewGracefulTeardown constructs GracefulTeardown checker.
func NewGracefulTeardown() GracefulTeardown { return GracefulTeardown{} }
func (GracefulTeardown) Name() string       { return "graceful-teardown" }

func (checker GracefulTeardown) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}
		if len(stack) < 3 {
			return true
		}

		fID := findSurroundingFunc(pass, stack)
		if fID == nil || !fID.meta.isTestCleanup {
			return true
		}

		call := NewCallMeta(pass, node.(*ast.CallExpr))
		if call == nil {
			return true
		}

		if !call.IsAssert {
			var fixes []analysis.SuggestedFix
			if fix := checker.buildFix(pass, call); fix != nil {
				fixes = append(fixes, *fix)
			}
			d := newDiagnostic(checker.Name(), call, gracefulTeardownReport, fixes...)
			diagnostics = append(diagnostics, *d)
			return false
		}
		return true
	})
	return diagnostics
}

func (GracefulTeardown) buildFix(pass *analysis.Pass, call *CallMeta) *analysis.SuggestedFix {
	if call.IsPkg() {
		// require.NoError(t, err) → assert.NoError(t, err)
		assertLocalName, importEdit, ok := addImportFix(pass.Files, call.Pos(), testify.AssertPkgPath)
		if !ok {
			return nil
		}

		textEdits := []analysis.TextEdit{
			{
				Pos:     call.Selector.X.Pos(),
				End:     call.Selector.X.End(),
				NewText: []byte(assertLocalName),
			},
		}
		if importEdit != nil {
			textEdits = append(textEdits, *importEdit)
		}
		return &analysis.SuggestedFix{
			Message:   fmt.Sprintf("Replace `%s` with `%s`", call.SelectorXStr, assertLocalName),
			TextEdits: textEdits,
		}
	}

	// s.Require().NoError(nil) → s.Assert().NoError(nil)
	requireCallExpr, ok := call.Selector.X.(*ast.CallExpr)
	if !ok {
		return nil
	}
	requireSel, ok := requireCallExpr.Fun.(*ast.SelectorExpr)
	if !ok || requireSel.Sel == nil || requireSel.Sel.Name != "Require" {
		return nil
	}
	return &analysis.SuggestedFix{
		Message: "Replace `Require` with `Assert`",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     requireSel.Sel.Pos(),
				End:     requireSel.Sel.End(),
				NewText: []byte("Assert"),
			},
		},
	}
}
