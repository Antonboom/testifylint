package checkers

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isNamedType(t types.Type, pkgPath string, name string) bool {
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}

	obj := n.Obj()
	if obj == nil || obj.Pkg() == nil || obj.Pkg().Path() != pkgPath {
		return false
	}

	return obj.Name() == name
}

func lookupMethod(pass *analysis.Pass, typ types.Type, addressable bool, name string) *types.Func {
	if typ == nil {
		return nil
	}

	obj, _, _ := types.LookupFieldOrMethod(typ, addressable, pass.Pkg, name)
	fn, _ := obj.(*types.Func)
	return fn
}
