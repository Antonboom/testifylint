package checkers

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

/*
	type TestingSuite interface {
		T() *testing.T
		SetT(*testing.T)
	}
*/

//testingPkg := types.NewPackage("testing", "testing")
//
//testingSuiteIface := types.NewInterfaceType([]*types.Func{
//	types.NewFunc(
//		token.NoPos,
//		testingPkg,
//		"T",
//		types.NewSignatureType(nil, nil, nil, nil,
//			types.NewTuple(types.NewVar(token.NoPos, testingPkg, "T", nil)), false),
//	),
//	types.NewFunc(
//		token.NoPos,
//		testingPkg,
//		"SetT",
//		types.NewSignatureType(nil, nil, nil, nil,
//			types.NewTuple(types.NewVar(token.NoPos, testingPkg, "T", nil)), false),
//	),
//}, nil)

//suiteTHelperInspect := func(node ast.Node) bool {
//	fd, ok := node.(*ast.FuncDecl)
//	if !ok {
//		return true
//	}
//	if fd.Recv == nil {
//		return true
//	}
//	if len(fd.Recv.List) != 1 {
//		return true
//	}
//	rcv := fd.Recv.List[0]
//	rcvType := pass.TypesInfo.TypeOf(rcv.Type)
//	if rcvType == nil {
//		return true
//	}
//
//	// Через Implements
//	// Через LookupFieldOrMethod и поиск методов
//	// Через LookupFieldOrMethod и поиск suite.Suite
//
//	s, _, _ := types.LookupFieldOrMethod(
//		rcvType,
//		false,
//		types.NewPackage("github.com/stretchr/testify/suite", "suite"),
//		"Suite",
//	)
//	if s == nil {
//		return true
//	}
//
//	if strings.HasPrefix(fd.Name.Name, "Test") {
//		return true
//	}
//	if unicode.IsUpper([]rune(fd.Name.Name)[0]) {
//		return true
//	}
//
//	if containsSuiteCalls(pass, rcv.Names[0].Name, rcvType, fd) {
//		if !firstStmtIsTHelperCall(pass, rcv.Names[0].Name, rcvType, fd) {
//			pass.Reportf(fd.Range(), "suite helper function should start from %s.T().Helper()", rcv.Names[0].Name)
//		}
//	}
//
//	return true
//}

// testingT := scope.Lookup("TestingT").Type().Underlying().(*types.Interface)

// https://github.com/stretchr/testify/blob/77977386932ab1866a4b9556a7af1ae347531786/_codegen/main.go#L158
// https://github.com/gostaticanalysis/analysisutil/blob/ccfdecf515f47e636ba164ce0e5f26810eaf8747/pkg.go#L21

func firstStmtIsTHelperCall(pass *analysis.Pass, rcvName string, rcvType types.Type, fn *ast.FuncDecl) bool {
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

	return types.ExprString(expr.X) == fmt.Sprintf("%s.T().Helper()", rcvName)
}

func containsSuiteCalls(pass *analysis.Pass, rcvName string, rcvType types.Type, fn *ast.FuncDecl) bool {
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

func isSuiteCall(pass *analysis.Pass, rcvName string, rcvType types.Type, s ast.Stmt) bool {
	expr, ok := s.(*ast.ExprStmt)
	if !ok {
		return false
	}

	ce, ok := expr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	x := unwrapSelector(ce.Fun)

	t := pass.TypesInfo.TypeOf(x)
	if t == nil {
		return false
	}

	return x.Name == rcvName && types.Identical(t, rcvType)
}

func unwrapSelector(e ast.Expr) *ast.Ident {
	for {
		switch v := e.(type) {
		case *ast.SelectorExpr:
			e = v.X

		case *ast.CallExpr:
			e = v.Fun

		case *ast.Ident:
			return v
		}
	}
}
