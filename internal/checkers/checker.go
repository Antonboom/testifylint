package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// CallMeta stores meta info about assert function/method call, for example
//   assert.Equal(t, 42, result, "helpful comment")
type CallMeta struct {
	analysis.Range
	Selector    *ast.SelectorExpr
	SelectorStr string
	Fn          FnMeta
	Args        []ast.Expr // Without t argument.
}

type FnMeta struct {
	analysis.Range
	Name  string
	IsFmt bool
}

type Checker interface {
	Name() string
	Check(pass *analysis.Pass, call CallMeta)
}
