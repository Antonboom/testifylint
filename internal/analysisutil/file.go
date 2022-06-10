package analysisutil

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func IsTestFile(pass *analysis.Pass, file *ast.File) bool {
	fname := pass.Fset.Position(file.Pos()).Filename
	return strings.HasSuffix(fname, "_test.go")
}
