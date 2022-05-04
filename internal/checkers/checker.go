package checkers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

type FnMeta struct {
	Pos        analysis.Range
	Pkg        string
	Name       string
	IsFormatFn bool
	Args       []ast.Expr // First arg is always *testing.T
}

//
//type Checker interface {
//	Name() string
//	Check(pass *analysis.Pass, fn FnMeta) error
//}

type Checker func(pass *analysis.Pass, fn FnMeta)
