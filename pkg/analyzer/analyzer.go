package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/checker"
	"github.com/Antonboom/testifylint/pkg/config"
)

const (
	name = "testifylint"
	doc  = "Checks usage of github.com/stretchr/testify."
)

// New accepts validated config.Config and returns testifylint analyzer.
func New(cfg config.Config) *analysis.Analyzer {
	tl := &testifyLint{
		checkers: newCheckers(cfg),
	}
	return &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		Run:      tl.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type testifyLint struct {
	checkers []checker.Checker
}

func (tl *testifyLint) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for _, f := range pass.Files {
		if !isTestFile(pass, f) {
			continue
		}

		if imports(f, "github.com/stretchr/testify/assert") ||
			imports(f, "github.com/stretchr/testify/require") ||
			imports(f, "github.com/stretchr/testify/suite") {
			insp.Nodes([]ast.Node{
				(*ast.CallExpr)(nil),
				(*ast.FuncDecl)(nil),
			}, tl.newCheckersRunner(pass))
		}
	}
	return nil, nil
}

func (tl *testifyLint) newCheckersRunner(pass *analysis.Pass) func(ast.Node, bool) bool {
	var insideSuiteMethod bool

	return func(node ast.Node, push bool) (proceed bool) {
		switch v := node.(type) {
		case *ast.FuncDecl:
			if isSuiteMethod(v, pass) {
				if push {
					insideSuiteMethod = true
				} else {
					insideSuiteMethod = false
				}
			}

		case *ast.CallExpr:
			tl.checkCall(v, pass, insideSuiteMethod)
		}
		return true
	}
}

func isSuiteMethod(fDecl *ast.FuncDecl, pass *analysis.Pass) bool {
	if fDecl.Recv == nil || len(fDecl.Recv.List) == 0 {
		return false
	}

	suiteIface := analysisutil.ObjectOf(pass, "github.com/stretchr/testify/suite", "TestingSuite")
	if suiteIface == nil {
		return false
	}

	rcv := fDecl.Recv.List[0]
	return types.Implements(
		pass.TypesInfo.TypeOf(rcv.Type),
		suiteIface.Type().Underlying().(*types.Interface),
	)
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

	isAssert := isPkg(pkg, "assert", "github.com/stretchr/testify/assert")
	isRequire := isPkg(pkg, "require", "github.com/stretchr/testify/require")
	if !(isAssert || isRequire) {
		return
	}

	call := checker.CallMeta{
		Range:             ce,
		Selector:          se,
		IsAssert:          isAssert,
		IsRequire:         isRequire,
		InsideSuiteMethod: insideSuiteMethod,
		SelectorStr:       types.ExprString(se.X),
		Fn: checker.FnMeta{
			Range: se.Sel,
			Name:  fn,
			IsFmt: strings.HasSuffix(fn, "f"),
		},
		Args: trimTArg(pass, ce.Args),
	}
	for _, ch := range tl.checkers {
		ch.Check(pass, call)
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
	ttObj := analysisutil.ObjectOf(pass, "testing", "T")
	if ttObj == nil {
		return false
	}

	argType := pass.TypesInfo.TypeOf(arg)
	if argType == nil {
		return false
	}

	return types.Identical(argType, types.NewPointer(ttObj.Type()))
}
