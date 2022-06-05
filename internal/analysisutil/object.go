package analysisutil

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

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

func trimVendor(path string) string {
	if strings.HasPrefix(path, "vendor/") {
		return path[len("vendor/"):]
	}
	return path
}
