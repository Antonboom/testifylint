package analyzer

import (
	"github.com/Antonboom/testifylint/internal/checkers"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

const (
	name = "testifylint"
	doc  = "Checks usage of github.com/stretchr/testify."
)

func New() *analysis.Analyzer {
	tl := &testifyLint{
		checkers: []checkers.Checker{
			checkers.Len,
		},
	}

	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		Run:  tl.run,
	}
}

type testifyLint struct {
	checkers []checkers.Checker
}

func (tl *testifyLint) run(pass *analysis.Pass) (interface{}, error) {
	// TODO: inspector

	inspect := func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg, pOk := se.X.(*ast.Ident)
		fn := se.Sel

		if pOk && isAssertOrRequire(pkg.Name) {
			fn := checkers.FnMeta{
				Pos:        ce.Lparen,
				Pkg:        pkg.Name,
				Name:       fn.Name,
				IsFormatFn: strings.HasSuffix(fn.Name, "f"),
				Args:       ce.Args,
			}
			for _, check := range tl.checkers {
				check(pass, fn)
			}
		}
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

func isAssertOrRequire(p string) bool {
	return p == "assert" || p == "require"
}
