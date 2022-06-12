//go:build debug

package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func skipFile(pass *analysis.Pass, node ast.Node) bool {
	return !strings.HasSuffix(pass.Fset.Position(node.Pos()).Filename, "float_compare_test.go")
}
