package checkers

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

type FnMeta struct {
	Pos        token.Pos
	Pkg        string
	Name       string
	IsFormatFn bool
	Args       []ast.Expr // First arg is always *testing.T
}

type Checker func(pass *analysis.Pass, fn FnMeta)
