package analysisutil

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func IsSuiteObj(pass *analysis.Pass, rcv ast.Expr) bool {
	suiteIface := ObjectOf(pass, "github.com/stretchr/testify/suite", "TestingSuite")
	if suiteIface == nil {
		return false
	}

	return types.Implements(
		pass.TypesInfo.TypeOf(rcv),
		suiteIface.Type().Underlying().(*types.Interface),
	)
}

func IsSuiteMethod(pass *analysis.Pass, fDecl *ast.FuncDecl) bool {
	if fDecl.Recv == nil || len(fDecl.Recv.List) != 1 {
		return false
	}

	rcv := fDecl.Recv.List[0]
	return IsSuiteObj(pass, rcv.Type)
}

func IsTestingTPtr(pass *analysis.Pass, arg ast.Expr) bool {
	ttObj := ObjectOf(pass, "testing", "T")
	if ttObj == nil {
		return false
	}

	argType := pass.TypesInfo.TypeOf(arg)
	if argType == nil {
		return false
	}

	return types.Identical(argType, types.NewPointer(ttObj.Type()))
}
