package analysisutil

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	errIface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

	falseObj = types.Universe.Lookup("false")
	lenObj   = types.Universe.Lookup("len")
	trueObj  = types.Universe.Lookup("true")
)

func IsUntypedTrue(pass *analysis.Pass, e ast.Expr) bool {
	return isObj(pass, e, trueObj)
}

func IsUntypedFalse(pass *analysis.Pass, e ast.Expr) bool {
	return isObj(pass, e, falseObj)
}

func IsComparisonWithTrue(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}
	if be.Op != op {
		return nil, false
	}

	t1, t2 := IsUntypedTrue(pass, be.X), IsUntypedTrue(pass, be.Y)
	if xor(t1, t2) {
		if t1 {
			return be.Y, true
		}
		return be.X, true
	}
	return nil, false
}

func IsComparisonWithFalse(pass *analysis.Pass, e ast.Expr, op token.Token) (ast.Expr, bool) {
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}
	if be.Op != op {
		return nil, false
	}

	f1, f2 := IsUntypedFalse(pass, be.X), IsUntypedFalse(pass, be.Y)
	if xor(f1, f2) {
		if f1 {
			return be.Y, true
		}
		return be.X, true
	}
	return nil, false
}

func IsTestifySuiteObj(pass *analysis.Pass, rcv ast.Expr) bool {
	suiteIface := objectOf(pass, "github.com/stretchr/testify/suite", "TestingSuite")
	if suiteIface == nil {
		return false
	}

	return types.Implements(
		pass.TypesInfo.TypeOf(rcv),
		suiteIface.Type().Underlying().(*types.Interface),
	)
}

func IsTestifySuiteMethod(pass *analysis.Pass, fDecl *ast.FuncDecl) bool {
	if fDecl.Recv == nil || len(fDecl.Recv.List) != 1 {
		return false
	}

	rcv := fDecl.Recv.List[0]
	return IsTestifySuiteObj(pass, rcv.Type)
}

func IsTestFile(pass *analysis.Pass, file *ast.File) bool {
	fname := pass.Fset.Position(file.Pos()).Filename
	return strings.HasSuffix(fname, "_test.go")
}

func IsTestingTPtr(pass *analysis.Pass, arg ast.Expr) bool {
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

func IsNegation(e ast.Expr) (ast.Expr, bool) {
	ue, ok := e.(*ast.UnaryExpr)
	if !ok {
		return nil, false
	}
	return ue.X, ue.Op == token.NOT
}

func IsBuiltinLenCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	if isObj(pass, ce.Fun, lenObj) && len(ce.Args) == 1 {
		return ce.Args[0], true
	}
	return nil, false
}

func IsIntNumber(e ast.Expr, v int) bool {
	bl, ok := e.(*ast.BasicLit)
	return ok && bl.Kind == token.INT && bl.Value == fmt.Sprintf("%d", v)
}

func IsError(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}

	_, ok := t.Underlying().(*types.Interface)
	return ok && types.Implements(t, errIface)
}

// objectOf returns types.Object for the given package and name
// and nil if the object is not found.
func objectOf(pass *analysis.Pass, pkg, name string) types.Object {
	if pass.Pkg.Path() == pkg {
		return pass.Pkg.Scope().Lookup(name)
	}

	for _, i := range pass.Pkg.Imports() {
		if trimVendor(i.Path()) == pkg {
			return i.Scope().Lookup(name)
		}
	}
	return nil
}

func isObj(pass *analysis.Pass, e ast.Expr, expected types.Object) bool {
	if expected == nil {
		panic("expect obj must be defined")
	}

	id, ok := e.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.ObjectOf(id)
	return obj == expected
}

func xor(a, b bool) bool {
	return a != b
}
