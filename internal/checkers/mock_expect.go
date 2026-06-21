package checkers

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
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
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		// -> u.On("CountUsers").Return(123)
		callExpr := node.(*ast.CallExpr)
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || !isTestifyMockMethod(pass, callExpr, "On", "Mock") {
			return true
		}

		// -> "CountUsers"
		methodName, ok := getMockMethodName(pass, callExpr)
		if !ok {
			return false
		}

		// -> func (_m *MockUserIFace) EXPECT() *MockUserIFace_Expecter
		receiverInfo := pass.TypesInfo.Types[selectorExpr.X]
		expectMethod := lookupMethod(pass, receiverInfo.Type, receiverInfo.Addressable(), "EXPECT")
		if expectMethod == nil {
			return false
		}

		expectSignature := expectMethod.Signature()
		if expectSignature.Params().Len() != 0 || expectSignature.Results().Len() != 1 {
			return false
		}

		// -> func (_e *MockUserIFace_Expecter) CountUsers() *MockUserIFace_CountUsers_Call
		method := lookupMethod(pass, expectSignature.Results().At(0).Type(), false, methodName)
		if method == nil || !callArgsMatch(pass, method.Signature(), callExpr.Args[1:], callExpr.Ellipsis.IsValid()) {
			return false
		}

		receiver := analysisutil.NodeString(pass.Fset, selectorExpr.X)
		report := "use " + receiver + ".EXPECT()." + methodName + "(...)"

		if inMockRunCall(pass, stack) {
			diagnostics = append(diagnostics, *newDiagnostic(checker.Name(), callExpr, report))
			return false
		}

		// -> u.EXPECT().CountUsers().Return(123)
		args := string(formatAsCallArgs(pass, callExpr.Args[1:]...))
		if callExpr.Ellipsis.IsValid() {
			args += "..."
		}
		newCall := fmt.Sprintf("%s.EXPECT().%s(%s)", receiver, methodName, args)

		diagnostics = append(diagnostics, *newDiagnostic(
			checker.Name(), callExpr, report,
			analysis.SuggestedFix{
				Message: "Replace mock.On with mock.EXPECT",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     callExpr.Pos(),
						End:     callExpr.End(),
						NewText: []byte(newCall),
					},
				},
			},
		))

		return false
	})

	return diagnostics
}

func isTestifyMockMethod(pass *analysis.Pass, callExpr *ast.CallExpr, methodName, receiverName string) bool {
	fn, ok := typeutil.Callee(pass.TypesInfo, callExpr).(*types.Func)
	if !ok || fn.Name() != methodName || fn.Pkg() == nil {
		return false
	}
	if !analysisutil.IsPkg(fn.Pkg(), testify.MockPkgName, testify.MockPkgPath) {
		return false
	}

	rcv := fn.Signature().Recv()
	if rcv == nil {
		return false
	}

	rcvType := rcv.Type()
	if ptr, isPtr := rcvType.(*types.Pointer); isPtr {
		rcvType = ptr.Elem()
	}

	named, ok := rcvType.(*types.Named)
	return ok && named.Obj().Name() == receiverName
}

func getMockMethodName(pass *analysis.Pass, callExpr *ast.CallExpr) (string, bool) {
	if len(callExpr.Args) == 0 {
		return "", false
	}

	value := pass.TypesInfo.Types[callExpr.Args[0]].Value
	if value == nil || value.Kind() != constant.String {
		return "", false
	}

	name := constant.StringVal(value)
	return name, token.IsIdentifier(name)
}

func lookupMethod(pass *analysis.Pass, typ types.Type, addressable bool, name string) *types.Func {
	if typ == nil {
		return nil
	}

	obj, _, _ := types.LookupFieldOrMethod(typ, addressable, pass.Pkg, name)
	fn, _ := obj.(*types.Func)
	return fn
}

func inMockRunCall(pass *analysis.Pass, stack []ast.Node) bool {
	for _, node := range stack {
		callExpr, ok := node.(*ast.CallExpr)
		if ok && isTestifyMockMethod(pass, callExpr, "Run", "Call") {
			return true
		}
	}
	return false
}
