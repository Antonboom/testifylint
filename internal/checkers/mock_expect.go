package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// MockExpect detects situations like
//
//	m.On("CreateUser", mock.Anything, User{}).Return(nil)
//	m.On("CountUsers").Return(123)
//
// and requires
//
//	m.EXPECT().CreateUser(mock.Anything, User{}).Return(nil)
//	m.EXPECT().CountUsers().Return(123)
type MockExpect struct{}

// NewMockExpect constructs MockExpect checker.
func NewMockExpect() MockExpect { return MockExpect{} }
func (MockExpect) Name() string { return "mock-expect" }

func (checker MockExpect) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		if !isSuiteTestMethod(funcDecl.Name.Name) {
			// process only Test functions
			return
		}

		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok || selectorExpr.Sel.Name != "On" {
				return true
			}

			methodName := firstArg(callExpr)
			if methodName == "" {
				return false
			}

			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return false
			}

			pointer, ok := pass.TypesInfo.ObjectOf(ident).Type().(*types.Pointer)
			if !ok {
				return false
			}

			named, ok := pointer.Elem().(*types.Named)
			if !ok {
				return false
			}

			if !hasExpect(named, methodName) {
				return false
			}

			diagnostics = append(diagnostics, *newDiagnostic(
				checker.Name(), callExpr, "use "+ident.Name+".EXPECT."+methodName+"(...)",
				analysis.SuggestedFix{
					Message: "Replace mock.On with mock.EXPECT",
					TextEdits: []analysis.TextEdit{
						{
							Pos: callExpr.Pos(),
							End: callExpr.End(),
							NewText: []byte(fmt.Sprintf(
								"%s.EXPECT().%s(%s)", ident.Name, methodName, formatAsCallArgs(pass, callExpr.Args[1:]...),
							)),
						},
					},
				},
			))

			return true
		})
	})

	return diagnostics
}

func firstArg(expr *ast.CallExpr) string {
	arg1, ok := expr.Args[0].(*ast.BasicLit)
	if !ok || arg1.Kind != token.STRING || len(arg1.Value) < 3 {
		return ""
	}
	return arg1.Value[1 : len(arg1.Value)-1]
}

// hasExpect checks if instead of .On("MethodName", ...) there is callable .EXPECT().MethodName(...)
func hasExpect(named *types.Named, methodName string) bool {
	for i := range named.NumMethods() {
		if named.Method(i).Name() == "EXPECT" && expectHasMethod(named.Method(i), methodName) {
			return true
		}
	}
	return false
}

func expectHasMethod(method *types.Func, methodName string) bool {
	pointer, ok := method.Signature().Results().At(0).Type().(*types.Pointer)
	if !ok {
		return false
	}

	named, ok := pointer.Elem().(*types.Named)
	if !ok {
		return false
	}

	for i := range named.NumMethods() {
		if named.Method(i).Name() == methodName {
			return true
		}
	}

	return false
}
