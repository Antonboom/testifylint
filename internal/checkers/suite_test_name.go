package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

// SuiteTestName detects situations where the suite runner function name does not match the suite name:
//
//	type BalanceSubscriptionSuite struct { suite.Suite }
//
//	❌ func TestBalanceSubs_Run(t *testing.T) { suite.Run(t, new(BalanceSubscriptionSuite)) }
//
//	✅ func TestBalanceSubscriptionSuite(t *testing.T) { suite.Run(t, new(BalanceSubscriptionSuite)) }
type SuiteTestName struct{}

// NewSuiteTestName constructs SuiteTestName checker.
func NewSuiteTestName() SuiteTestName { return SuiteTestName{} }
func (SuiteTestName) Name() string    { return "suite-test-name" }

func (checker SuiteTestName) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	suiteRunObj := analysisutil.ObjectOf(pass.Pkg, testify.SuitePkgPath, "Run")
	if suiteRunObj == nil {
		return nil
	}

	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		ce := node.(*ast.CallExpr) // e.g. suite.Run(t, new(SomeSuite))

		// Check that the call is suite.Run.
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if !analysisutil.IsObj(pass.TypesInfo, se.Sel, suiteRunObj) {
			return true
		}

		// suite.Run(t, suiteInstance) – at least 2 args required.
		if len(ce.Args) < 2 {
			return true
		}

		// Extract the suite type name from the second argument.
		// Supported forms: new(SuiteName) or &SuiteName{}.
		suiteName := extractSuiteTypeName(ce.Args[1])
		if suiteName == "" {
			return true
		}

		expectedFnName := "Test" + suiteName

		// Find the containing function declaration.
		for i := len(stack) - 2; i >= 0; i-- {
			fd, ok := stack[i].(*ast.FuncDecl)
			if !ok {
				continue
			}

			if fd.Name == nil {
				return false
			}

			// Only lint top-level test runner functions:
			// no receiver, name starts with "Test", and takes a *testing.T parameter.
			if fd.Recv != nil {
				return false
			}
			if !strings.HasPrefix(fd.Name.Name, "Test") {
				return false
			}
			if !hasStdTestingTParam(pass, fd.Type) {
				return false
			}

			if fd.Name.Name == expectedFnName {
				return false
			}

			msg := fmt.Sprintf("suite test function name %s does not match suite name (expected %s)",
				fd.Name.Name, expectedFnName)
			d := newDiagnostic(checker.Name(), fd.Name, msg, analysis.SuggestedFix{
				Message: "Rename to " + expectedFnName,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     fd.Name.Pos(),
						End:     fd.Name.End(),
						NewText: []byte(expectedFnName),
					},
				},
			})
			diagnostics = append(diagnostics, *d)
			return false
		}
		return true
	})
	return diagnostics
}

// extractSuiteTypeName returns the suite struct type name from an expression.
// It supports:
//   - new(SomeSuite)     → "SomeSuite"
//   - &SomeSuite{}       → "SomeSuite"
func extractSuiteTypeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.CallExpr:
		// Handles new(SomeSuite) form.
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "new" {
			if len(e.Args) == 1 {
				if id, ok := e.Args[0].(*ast.Ident); ok {
					return id.Name
				}
			}
		}
	case *ast.UnaryExpr:
		// &SomeSuite{...}
		if e.Op == token.AND {
			if cl, ok := e.X.(*ast.CompositeLit); ok {
				if id, ok := cl.Type.(*ast.Ident); ok {
					return id.Name
				}
			}
		}
	}
	return ""
}
