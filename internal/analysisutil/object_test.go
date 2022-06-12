package analysisutil_test

import (
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"

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

	pkg := types.NewPackage("mycoolapp/service", "servicetest")
	pkg.SetImports([]*types.Package{testingPkg, assertPkg, vendoredSuitePkg})
	timeoutIface := types.NewTypeName(
		token.NoPos,
		pkg,
		"Timeout",
		types.NewInterfaceType(nil, nil),
	)
	pkg.Scope().Insert(timeoutIface)

	pass := &analysis.Pass{Pkg: pkg}

	cases := []struct {
		pkg, name string
		expObj    types.Object
	}{
		{pkg: "mycoolapp/service", name: "Timeout", expObj: timeoutIface},
		{pkg: "testing", name: "T", expObj: testingT},
		{pkg: "github.com/stretchr/testify/assert", name: "Equal", expObj: assertEqual},
		{pkg: "github.com/stretchr/testify/suite", name: "TestingSuite", expObj: testingSuiteIface},

		// Negative.
		{pkg: "net/http", name: "Timeout", expObj: nil},
		{pkg: "testing", name: "TT", expObj: nil},
		{pkg: "github.com/stretchr/testify/assert", name: "NotEqual", expObj: nil},
		{pkg: "vendor/github.com/stretchr/testify/assert", name: "Equal", expObj: nil},
	}
	for _, tt := range cases {
		t.Run(tt.pkg+"."+tt.name, func(t *testing.T) {
			obj := analysisutil.ObjectOf(pass, tt.pkg, tt.name)
			if obj != tt.expObj {
				t.Fatalf("unexpected: %v", obj)
			}
		})
	}
}
