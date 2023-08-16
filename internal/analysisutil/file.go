package analysisutil

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strconv"
	"strings"

	"slices"
)

// IsTestFile returns true if file is test.
func IsTestFile(pass *analysis.Pass, file *ast.File) bool {
	fname := pass.Fset.Position(file.Pos()).Filename
	return strings.HasSuffix(fname, "_test.go")
}

// Imports tells if the file imports at least one of the pkgs.
func Imports(file *ast.File, pkgs ...string) bool {
	for _, i := range file.Imports {
		if i.Path == nil {
			continue
		}

		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if slices.Contains(pkgs, path) {
			return true
		}
	}
	return false
}
