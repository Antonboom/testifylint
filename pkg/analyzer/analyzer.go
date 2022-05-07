package analyzer

import (
	"github.com/Antonboom/testifylint/internal/checkers"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

const (
	name = "testifylint"
	doc  = "Checks usage of github.com/stretchr/testify."
)

func New() *analysis.Analyzer {
	tl := &testifyLint{
		checkers: []checkers.Checker{ // Order is important!
			//checkers.BoolCompare,
			//checkers.FloatCompare,
			//checkers.Empty,
			//checkers.Len,
			//checkers.Comparisons,
			//checkers.Error,
			checkers.NewErrorIs(),
			//checkers.ExpectedActual(nil),
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

func (tl *testifyLint) run(pass *analysis.Pass) (any, error) {
	// TODO: inspector
	for _, f := range pass.Files {
		if !isTestFile(pass, f) {
			continue
		}

		if imports(f, "github.com/stretchr/testify/assert") ||
			imports(f, "github.com/stretchr/testify/require") {
			ast.Inspect(f, tl.newCheckersRunner(pass))
		}

		if imports(f, "github.com/stretchr/testify/suite") {
			ast.Inspect(f, tl.newSuiteSpecificCheckersRunner(pass))
		}
	}
	return nil, nil
}

func (tl *testifyLint) newCheckersRunner(pass *analysis.Pass) func(ast.Node) bool {
	return func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if se.Sel == nil {
			return true
		}
		fn := se.Sel.Name

		pkg := func() *types.Package {
			if sel, ok := pass.TypesInfo.Selections[se]; ok {
				return sel.Obj().Pkg()
			}

			if id, ok := se.X.(*ast.Ident); ok {
				if selObj := pass.TypesInfo.ObjectOf(id); selObj != nil {
					if pkg, ok := selObj.(*types.PkgName); ok {
						return pkg.Imported()
					}
				}
			}
			return nil
		}()
		if pkg == nil {
			return true
		}
		if !(isPkg(pkg, "assert", "github.com/stretchr/testify/assert") ||
			isPkg(pkg, "require", "github.com/stretchr/testify/require")) {
			return true
		}

		call := checkers.CallMeta{
			Range:       ce,
			Selector:    se,
			SelectorStr: types.ExprString(se.X),
			Fn: checkers.FnMeta{
				Range: se.Sel,
				Name:  fn,
				IsFmt: strings.HasSuffix(fn, "f"),
			},
			Args: trimTArg(pass, ce.Args),
		}
		for _, checker := range tl.checkers {
			checker.Check(pass, call)
		}
		return true
	}
}

func (tl *testifyLint) newSuiteSpecificCheckersRunner(pass *analysis.Pass) func(ast.Node) bool {
	return func(node ast.Node) bool {
		return false
	}
}

func trimTArg(pass *analysis.Pass, args []ast.Expr) []ast.Expr {
	if len(args) == 0 {
		return args
	}

	if isTestingTPtr(pass, args[0]) {
		return args[1:]
	}
	return args
}

func isTestingTPtr(pass *analysis.Pass, arg ast.Expr) bool {
	ttObj := objectOf(pass, "testing", "T")
	if ttObj == nil {
		return false
	}

	argType := pass.TypesInfo.TypeOf(arg)
	if argType == nil {
		return false
	}

	return types.Identical(argType, types.NewPointer(ttObj.Type()))
}

func implementsTestingSuite(pass *analysis.Pass) bool {
	//tSuiteObj := objectOf(pass, "github.com/stretchr/testify/suite", "TestingSuite")

	//tSuiteIface
	return false
}
