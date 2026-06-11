package checkers

import "go/types"

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
