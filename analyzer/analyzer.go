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
	"github.com/Antonboom/testifylint/internal/testify"
)

const (
	name = "testifylint"
	doc  = "Checks usage of " + testify.ModulePath + "."
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

type testifyLint struct {
	regularCheckers  []checkers.RegularChecker
	advancedCheckers []checkers.AdvancedChecker
}

func (tl *testifyLint) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for _, f := range pass.Files {
		if !analysisutil.IsTestFile(pass.Fset, f) {
			continue
		}
		if !analysisutil.Imports(f, testify.AssertPkgPath, testify.RequirePkgPath, testify.SuitePkgPath) {
			continue
		}

		// Regular checkers.
		insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
			tl.regularCheck(pass, node.(*ast.CallExpr))
		})

		// Advanced checkers.
		for _, ch := range tl.advancedCheckers {
			for _, d := range ch.Check(pass, insp) {
				pass.Report(d)
			}
		}
	}
	return nil, nil
}

func (tl *testifyLint) regularCheck(pass *analysis.Pass, ce *ast.CallExpr) {
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if se.Sel == nil {
		return
	}
	fn := se.Sel.Name

	pkg := func() *types.Package {
		// Examples:
		// s.Assert         -> *suite.Suite        -> package suite ("vendor/github.com/stretchr/testify/suite")
		// s.Assert().Equal -> *assert.Assertions  -> package assert ("vendor/github.com/stretchr/testify/assert")
		// reqObj.Falsef    -> *require.Assertions -> package require ("vendor/github.com/stretchr/testify/require")
		if sel, ok := pass.TypesInfo.Selections[se]; ok {
			return sel.Obj().Pkg()
		}

		// Examples:
		// assert.False      -> assert  -> package assert ("vendor/github.com/stretchr/testify/assert")
		// require.NotEqualf -> require -> package require ("vendor/github.com/stretchr/testify/require")
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

	isAssert := analysisutil.IsPkg(pkg, testify.AssertPkgName, testify.AssertPkgPath)
	isRequire := analysisutil.IsPkg(pkg, testify.RequirePkgName, testify.RequirePkgPath)
	if !(isAssert || isRequire) {
		return
	}

	call := &checkers.CallMeta{
		Range:        ce,
		IsAssert:     isAssert,
		IsRequire:    isRequire,
		Selector:     se,
		SelectorXStr: analysisutil.NodeString(pass.Fset, se.X),
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

	if isTestingTPtr(pass, args[0]) {
		return args[1:]
	}
	return args
}

func isTestingTPtr(pass *analysis.Pass, arg ast.Expr) bool {
	ttObj := analysisutil.ObjectOf(pass.Pkg, "testing", "T")
	if ttObj == nil {
		return false
	}

	argType := pass.TypesInfo.TypeOf(arg)
	if argType == nil {
		return false
	}

	return types.Identical(argType, types.NewPointer(ttObj.Type()))
}
