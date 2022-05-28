package analyzer

import (
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func isPkg(pkg *types.Package, name, path string) bool {
	return pkg.Name() == name && trimVendor(pkg.Path()) == path
}

func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	fname := pass.Fset.Position(file.Pos()).Filename
	return strings.HasSuffix(fname, "_test.go")
}

func objectOf(pass *analysis.Pass, pkg, name string) types.Object {
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

func imports(file *ast.File, pkg string) bool {
	for _, i := range file.Imports {
		if i.Path == nil {
			continue
		}

		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}

		if trimVendor(path) == pkg {
			return true
		}
	}
	return false
}

func trimVendor(path string) string {
	if strings.HasPrefix(path, "vendor/") {
		return path[len("vendor/"):]
	}
	return path
}
