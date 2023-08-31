package analysisutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestNodeFormatting(t *testing.T) {
	const src = `
package p
var _ = make(chan int, 1)
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var chBytes []byte
	var chStr string

	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.GenDecl:
			chBytes = analysisutil.NodeBytes(fset, v)

		case *ast.CallExpr:
			chStr = analysisutil.NodeString(fset, v)
		}
		return true
	})

	t.Run("NodeBytes", func(t *testing.T) {
		if string(chBytes) != "var _ = make(chan int, 1)" {
			t.Fatalf("%s", chBytes)
		}
	})

	t.Run("NodeString", func(t *testing.T) {
		if chStr != "make(chan int, 1)" {
			t.Fatal(chStr)
		}
	})
}

func TestNodeFormatting_Invalid(t *testing.T) {
	t.Run("NodeBytes", func(t *testing.T) {
		b := analysisutil.NodeBytes(token.NewFileSet(), nil)
		if string(b) != "" {
			t.Fatalf("%s", b)
		}
	})

	t.Run("NodeString", func(t *testing.T) {
		str := analysisutil.NodeString(token.NewFileSet(), nil)
		if str != "" {
			t.Fatal(str)
		}
	})
}
