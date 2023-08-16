package analysisutil

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// ObjectOf returns types.Object for the given package and name
// and nil if the object is not found.
func ObjectOf(pass *analysis.Pass, pkg, name string) types.Object {
	if pass.Pkg.Path() == pkg {
		return pass.Pkg.Scope().Lookup(name)
	}

	for _, i := range pass.Pkg.Imports() {
		if trimVendor(i.Path()) == pkg {
			return i.Scope().Lookup(name)
		}
	}
	return nil
}

func IsObj(pass *analysis.Pass, e ast.Expr, expected types.Object) bool {
	if expected == nil {
		panic("expected obj must be defined")
	}

	id, ok := e.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.ObjectOf(id)
	return obj == expected
}
