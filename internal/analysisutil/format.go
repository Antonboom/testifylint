package analysisutil

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

// NodeString is a more powerful analogue of types.ExprString.
// Return empty string if node AST is invalid.
func NodeString(fset *token.FileSet, node ast.Node) string {
	if v := formatNode(fset, node); v != nil {
		return v.String()
	}
	return ""
}

func formatNode(fset *token.FileSet, node ast.Node) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if err := format.Node(buf, fset, node); err != nil {
		return nil
	}
	return buf
}
