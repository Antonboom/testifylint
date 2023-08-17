package analysisutil

import (
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

// IsTestFile returns true if the file from the file set has the `_test.go` suffix.
// If the file does not belong to the set, then the function will return false.
func IsTestFile(fset *token.FileSet, file *ast.File) bool {
	fname := fset.Position(file.Pos()).Filename
	return strings.HasSuffix(fname, "_test.go")
}

// Imports tells if the file imports at least one of the packages.
// If no packages provided then function returns false.
func Imports(file *ast.File, pkgs ...string) bool {
	for _, i := range file.Imports {
		if i.Path == nil {
			continue
		}

		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if slices.Contains(pkgs, path) { // Small O(n).
			return true
		}
	}
	return false
}
