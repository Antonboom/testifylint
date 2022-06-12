package analysisutil

import (
	"go/ast"
	"go/types"
	"strconv"
	"strings"
)

// IsPkg checks that pkg has provided name & path.
// Supports vendored packages.
func IsPkg(pkg *types.Package, name, path string) bool {
	return pkg.Name() == name && trimVendor(pkg.Path()) == path
}

// Imports tells if the file imports the pkg.
func Imports(file *ast.File, pkg string) bool {
	for _, i := range file.Imports {
		if i.Path == nil {
			continue
		}

		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if path == pkg {
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
