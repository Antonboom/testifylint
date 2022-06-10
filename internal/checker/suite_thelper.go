package checker

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

type SuiteTHelper struct{}

func NewSuiteTHelper() SuiteTHelper {
	return SuiteTHelper{}
}

func (SuiteTHelper) Name() string       { return "suite-thelper" }
func (SuiteTHelper) Priority() int      { return 11 }
func (SuiteTHelper) DisabledByDefault() { return }

func (checker SuiteTHelper) Check(pass *analysis.Pass, inspector *inspector.Inspector) {
	inspector.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		fd := node.(*ast.FuncDecl)
		if !analysisutil.IsSuiteMethod(pass, fd) {
			return
		}

		if ident := fd.Name; ident == nil ||
			strings.HasPrefix(ident.Name, "Test") || isServiceSuiteMethod(ident.Name) {
			return
		}

		rcv := fd.Recv.List[0]
		if len(rcv.Names) != 1 || rcv.Names[0] == nil {
			return
		}
		rcvName := rcv.Names[0].Name

		rcvType := pass.TypesInfo.TypeOf(rcv.Type)
		if rcvType == nil {
			return
		}

		if !containsSuiteCalls(pass, fd, rcvName, rcvType) {
			return
		}

		if !firstStmtIsTHelperCall(pass, fd, rcvName, rcvType) {
			msg := fmt.Sprintf("suite helper method should start with %s.T().Helper()", rcvName)
			r.Report(pass, checker.Name(), fd, msg, nil)
		}
	})
}

func isServiceSuiteMethod(name string) bool {
	// github.com/stretchr/testify/suite/interfaces.go
	switch name {
	case "SetupSuite", "SetupTest", "TearDownSuite", "TearDownTest", "BeforeTest", "AfterTest", "HandleStats":
		return true
	}
	return false
}

func containsSuiteCalls(pass *analysis.Pass, fn *ast.FuncDecl, rcvName string, rcvType types.Type) bool {
	if fn.Body == nil {
		return false
	}

	for _, s := range fn.Body.List {
		if isSuiteCall(pass, rcvName, rcvType, s) {
			return true
		}
	}
	return false
}

func firstStmtIsTHelperCall(pass *analysis.Pass, fn *ast.FuncDecl, rcvName string, rcvType types.Type) bool {
	if fn.Body == nil {
		return false
	}

	if len(fn.Body.List) == 0 {
		return false
	}
	s := fn.Body.List[0]

	expr, ok := s.(*ast.ExprStmt)
	if !ok {
		return false
	}
	return isSuiteCall(pass, rcvName, rcvType, s) &&
		types.ExprString(expr.X) == fmt.Sprintf("%s.T().Helper()", rcvName)
}

func isSuiteCall(pass *analysis.Pass, rcvName string, rcvType types.Type, s ast.Stmt) bool {
	expr, ok := s.(*ast.ExprStmt)
	if !ok {
		return false
	}

	x := unwrapSelector(expr.X)
	if x == nil {
		return false
	}

	t := pass.TypesInfo.TypeOf(x)
	if t == nil {
		return false
	}
	return x.Name == rcvName && types.Identical(t, rcvType)
}

// unwrapSelector supports
//  s.True(b)
//  s.Assert().True(b)
//  s.Require().True(b)
// and returns "s" identifier.
func unwrapSelector(e ast.Expr) *ast.Ident {
	for {
		switch v := e.(type) {
		case *ast.CallExpr:
			e = v.Fun

		case *ast.SelectorExpr:
			e = v.X

		case *ast.Ident:
			return v

		default:
			// Protection against strange constructs that cause an infinite loop.
			return nil
		}
	}
}
