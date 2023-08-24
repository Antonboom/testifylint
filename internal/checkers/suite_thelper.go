package checkers

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// SuiteTHelper wants t.Helper() call in suite helpers:
//
//	func (s *RoomSuite) assertRoomRound(roundID RoundID) {
//		s.T().Helper()
//		s.Equal(roundID, s.getRoom().CurrentRound.ID)
//	}
type SuiteTHelper struct{}

// NewSuiteTHelper constructs SuiteTHelper checker.
func NewSuiteTHelper() SuiteTHelper { return SuiteTHelper{} }
func (SuiteTHelper) Name() string   { return "suite-thelper" }

func (checker SuiteTHelper) Check(pass *analysis.Pass, inspector *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	inspector.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		fd := node.(*ast.FuncDecl)
		if !isTestifySuiteMethod(pass, fd) {
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

		if !containsSuiteCalls(pass, fd, rcvType) {
			return
		}

		firstStmt := getFirstStatement(fd)
		if firstStmt == nil {
			panic("containsSuiteCalls works incorrectly")
		}

		if isTHelperCall(pass, firstStmt, rcvType) {
			return
		}

		helper := fmt.Sprintf("%s.T().Helper()", rcvName)
		msg := fmt.Sprintf("suite helper method should start with " + helper)
		d := newDiagnostic(checker.Name(), fd, msg, &analysis.SuggestedFix{
			Message: "Insert " + helper,
			TextEdits: []analysis.TextEdit{
				{
					Pos:     firstStmt.Pos(),
					End:     firstStmt.Pos(),
					NewText: []byte(helper + "\n\n"),
				},
			},
		})
		diagnostics = append(diagnostics, *d)
	})
	return diagnostics
}

func isTestifySuiteMethod(pass *analysis.Pass, fDecl *ast.FuncDecl) bool {
	if fDecl.Recv == nil || len(fDecl.Recv.List) != 1 {
		return false
	}

	rcv := fDecl.Recv.List[0]
	return implementsTestifySuiteIface(pass, rcv.Type)
}

func isServiceSuiteMethod(name string) bool {
	// https://github.com/stretchr/testify/blob/master/suite/interfaces.go
	switch name {
	case "SetupSuite", "SetupTest", "TearDownSuite", "TearDownTest",
		"BeforeTest", "AfterTest", "HandleStats", "SetupSubTest", "TearDownSubTest":
		return true
	}
	return false
}

func containsSuiteCalls(pass *analysis.Pass, fn *ast.FuncDecl, rcvType types.Type) bool {
	if fn.Body == nil {
		return false
	}

	for _, s := range fn.Body.List {
		if isSuiteCall(pass, s, rcvType) {
			return true
		}
	}
	return false
}

func getFirstStatement(fn *ast.FuncDecl) ast.Stmt {
	if fn.Body == nil {
		return nil
	}

	if len(fn.Body.List) == 0 {
		return nil
	}
	return fn.Body.List[0]
}

func isTHelperCall(pass *analysis.Pass, stmt ast.Stmt, rcvType types.Type) bool {
	return isSuiteCall(pass, stmt, rcvType) &&
		(strings.HasSuffix(analysisutil.NodeString(pass.Fset, stmt), "T().Helper()"))
}

func isSuiteCall(pass *analysis.Pass, stmt ast.Stmt, rcvType types.Type) bool {
	expr, ok := stmt.(*ast.ExprStmt)
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
	return types.Identical(t, rcvType)
}

// unwrapSelector supports
//
//	s.True(b)
//	s.Assert().True(b)
//	s.Require().True(b)
//
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
