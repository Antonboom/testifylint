package checkers

import (
	"bytes"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// formatAsCallArgs joins a, b and c and returns bytes like `a, b, c`.
func formatAsCallArgs(pass *analysis.Pass, args ...ast.Expr) []byte {
	if len(args) == 0 {
		return []byte("")
	}

	var buf bytes.Buffer
	for i, arg := range args {
		buf.Write(analysisutil.NodeBytes(pass.Fset, arg))
		if i != len(args)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.Bytes()
}

// formatWithStringCastForBytes formats the expression as a string argument for JSONEq/YAMLEq.
// It wraps []byte expressions with string() and replaces buffer.Bytes() with buffer.String().
// Returns nil if the expression cannot be formatted as a string.
func formatWithStringCastForBytes(pass *analysis.Pass, e ast.Expr) []byte {
	if !hasBytesType(pass, e) {
		t := pass.TypesInfo.TypeOf(e)
		if t == nil {
			// We cannot format this expression as a string.
			return nil
		}
		bt, ok := t.(*types.Basic)
		if !ok || bt.Info()&types.IsString == 0 {
			// We cannot format this expression as a string.
			return nil
		}
		return analysisutil.NodeBytes(pass.Fset, e)
	}

	if se, ok := isBufferBytesCall(pass, e); ok {
		return []byte(analysisutil.NodeString(pass.Fset, se) + ".String()")
	}
	return []byte("string(" + analysisutil.NodeString(pass.Fset, e) + ")")
}

func isBufferBytesCall(pass *analysis.Pass, e ast.Expr) (ast.Node, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	if !isIdentWithName("Bytes", se.Sel) {
		return nil, false
	}
	if t := pass.TypesInfo.TypeOf(se.X); t != nil {
		// NOTE(a.telyshev): This is hack, because `bytes` package can be not imported,
		// and we cannot do "true" comparison with `Buffer` object.
		return se.X, strings.TrimPrefix(t.String(), "*") == "bytes.Buffer"
	}

	return nil, false
}
