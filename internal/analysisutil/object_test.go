package analysisutil_test

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestObjectOf(t *testing.T) {
	testingPkg := types.NewPackage("testing", "testing")
	testingT := types.NewVar(token.NoPos, testingPkg, "T", types.NewPointer(types.NewStruct(nil, nil)))
	testingPkg.Scope().Insert(testingT)

	assertPkg := types.NewPackage("github.com/stretchr/testify/assert", "assert")
	assertEqual := types.NewFunc(token.NoPos, assertPkg, "Equal", nil)
	assertPkg.Scope().Insert(assertEqual)

	vendoredSuitePkg := types.NewPackage("vendor/github.com/stretchr/testify/suite", "suite")
	testingSuiteIface := types.NewTypeName(
		token.NoPos,
		vendoredSuitePkg,
		"TestingSuite",
		types.NewInterfaceType(nil, nil),
	)
	vendoredSuitePkg.Scope().Insert(testingSuiteIface)

	pkg := types.NewPackage("mycoolapp/service", "service_test")
	pkg.SetImports([]*types.Package{testingPkg, assertPkg, vendoredSuitePkg})
	timeoutIface := types.NewTypeName(
		token.NoPos,
		pkg,
		"Timeout",
		types.NewInterfaceType(nil, nil),
	)
	pkg.Scope().Insert(timeoutIface)

	cases := []struct {
		objPkg, objName string
		expObj          types.Object
	}{
		// Positive.
		{objPkg: "mycoolapp/service", objName: "Timeout", expObj: timeoutIface},
		{objPkg: "testing", objName: "T", expObj: testingT},
		{objPkg: "github.com/stretchr/testify/assert", objName: "Equal", expObj: assertEqual},
		{objPkg: "github.com/stretchr/testify/suite", objName: "TestingSuite", expObj: testingSuiteIface},

		// Negative.
		{objPkg: "net/http", objName: "Timeout", expObj: nil},
		{objPkg: "testing", objName: "TT", expObj: nil},
		{objPkg: "github.com/stretchr/testify/assert", objName: "NotEqual", expObj: nil},
		{objPkg: "vendor/github.com/stretchr/testify/assert", objName: "Equal", expObj: nil},
	}
	for _, tt := range cases {
		t.Run(tt.objPkg+"."+tt.objName, func(t *testing.T) {
			obj := analysisutil.ObjectOf(pkg, tt.objPkg, tt.objName)
			if obj != tt.expObj {
				t.Fatalf("unexpected: %v", obj)
			}
		})
	}
}

func TestIsObj(t *testing.T) {
	lenIdent, lenObj := ast.NewIdent("len"), types.Universe.Lookup("len")
	falseIdent, falseObj := ast.NewIdent("false"), types.Universe.Lookup("false")

	typesInfo := &types.Info{
		Defs: map[*ast.Ident]types.Object{
			lenIdent:   lenObj,
			falseIdent: falseObj,
		},
	}

	cases := []struct {
		expr        ast.Expr
		expectedObj types.Object
		isObj       bool
	}{
		{
			expr:        lenIdent,
			expectedObj: lenObj,
			isObj:       true,
		},
		{
			expr:        falseIdent,
			expectedObj: falseObj,
			isObj:       true,
		},
		{
			expr:        lenIdent,
			expectedObj: falseObj,
			isObj:       false,
		},
		{
			expr:        falseIdent,
			expectedObj: lenObj,
			isObj:       false,
		},
		{
			expr:        &ast.BasicLit{Value: "42"},
			expectedObj: lenObj,
			isObj:       false,
		},
	}
	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			v := analysisutil.IsObj(typesInfo, tt.expr, tt.expectedObj)
			if v != tt.isObj {
				t.FailNow()
			}
		})
	}
}

func TestIsObj_NamesakesFromDifferentPackages(t *testing.T) {
	lhs := types.NewFunc(token.NoPos, types.NewPackage("errors", "errors"), "Is", nil)
	rhs := types.NewFunc(token.NoPos, types.NewPackage("pkg/errors", "errors"), "Is", nil)

	ident := new(ast.Ident)
	typesInfo := &types.Info{
		Defs: map[*ast.Ident]types.Object{
			ident: lhs,
		},
	}
	if analysisutil.IsObj(typesInfo, ident, rhs) {
		t.Fatalf("objects should not be equal")
	}
}
