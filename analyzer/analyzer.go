package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/config"
	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/checkers"
)

const (
	name = "testifylint"
	doc  = "Checks usage of " + testifyPath + "."
)

// New accepts config.Config and returns testifylint analyzer.
func New(cfg config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		Run: func(pass *analysis.Pass) (any, error) {
			if err := config.Validate(cfg); err != nil {
				return nil, fmt.Errorf("invalid config: %v", err)
			}

			regularCheckers, advancedCheckers, err := newCheckers(cfg)
			if err != nil {
				return nil, fmt.Errorf("build checkers: %v", err)
			}

			tl := &testifyLint{
				regularCheckers:  regularCheckers,
				advancedCheckers: advancedCheckers,
			}
			return tl.run(pass)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

const (
	testifyPath = "github.com/stretchr/testify"

	assertPkgName  = "assert"
	requirePkgName = "require"
	suitePkgName   = "suite"

	testifyAssertPath  = testifyPath + "/" + assertPkgName
	testifyRequirePath = testifyPath + "/" + requirePkgName
	testifySuitePath   = testifyPath + "/" + suitePkgName
)

type testifyLint struct {
	regularCheckers  []checkers.RegularChecker
	advancedCheckers []checkers.AdvancedChecker
}

func (tl *testifyLint) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for _, f := range pass.Files {
		if !analysisutil.IsTestFile(pass, f) {
			continue
		}

		if analysisutil.Imports(f, testifyAssertPath, testifyRequirePath, testifySuitePath) {
			insp.Nodes([]ast.Node{
				(*ast.CallExpr)(nil),
				(*ast.FuncDecl)(nil),
			}, tl.newRegularCheckersRunner(pass))

			for _, ch := range tl.advancedCheckers {
				for _, d := range ch.Check(pass, insp) {
					pass.Report(d)
				}
			}
		}
	}
	return nil, nil
}

func (tl *testifyLint) newRegularCheckersRunner(pass *analysis.Pass) func(ast.Node, bool) bool {
	var insideSuiteMethod bool

	return func(node ast.Node, push bool) (proceed bool) {
		switch v := node.(type) {
		case *ast.FuncDecl:
			if analysisutil.IsTestifySuiteMethod(pass, v) {
				if push {
					insideSuiteMethod = true
				} else {
					insideSuiteMethod = false
				}
			}

		case *ast.CallExpr:
			// NOTE(a.telyshev): Process call expressions once.
			if push {
				tl.checkCall(v, pass, insideSuiteMethod)
			}
		}
		return true
	}
}

func (tl *testifyLint) checkCall(ce *ast.CallExpr, pass *analysis.Pass, insideSuiteMethod bool) {
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if se.Sel == nil {
		return
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
		return
	}

	isAssert := analysisutil.IsPkg(pkg, assertPkgName, testifyAssertPath)
	isRequire := analysisutil.IsPkg(pkg, requirePkgName, testifyRequirePath)
	if !(isAssert || isRequire) {
		return
	}

	call := &checkers.CallMeta{
		Range:             ce,
		IsAssert:          isAssert,
		IsRequire:         isRequire,
		InsideSuiteMethod: insideSuiteMethod,
		Selector:          se,
		SelectorXStr:      types.ExprString(se.X),
		Fn: checkers.FnMeta{
			Range: se.Sel,
			Name:  fn,
			IsFmt: strings.HasSuffix(fn, "f"),
		},
		Args:    trimTArg(pass, ce.Args),
		ArgsRaw: ce.Args,
	}
	for _, ch := range tl.regularCheckers {
		if d := ch.Check(pass, call); d != nil {
			pass.Report(*d)
			// NOTE(a.telyshev): I'm not interested in multiple diagnostics per assertion.
			// This simplifies the code and also makes the linter more efficient.
			return
		}
	}
}

func trimTArg(pass *analysis.Pass, args []ast.Expr) []ast.Expr {
	if len(args) == 0 {
		return args
	}

	if analysisutil.IsTestingTPtr(pass, args[0]) {
		return args[1:]
	}
	return args
}
