package checkers

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

// CallMeta stores meta info about assertion function/method call, for example
//
//	assert.Equal(t, 42, result, "helpful comment")
type CallMeta struct {
	// Call stores the original AST call expression.
	Call *ast.CallExpr
	// Range contains start and end position of assertion call.
	analysis.Range
	// Obj contains Assertions or CollectT object, can be nil in case of package (not object) call.
	Obj types.Object
	// IsAssert true if this is "testify/assert" package (or object) call.
	IsAssert bool
	// Selector is the AST expression of "assert.Equal".
	Selector *ast.SelectorExpr
	// SelectorXStr is a string representation of Selector's left part – value before point, e.g. "assert".
	SelectorXStr string
	// Fn stores meta info about assertion function itself.
	Fn FnMeta
	// Args stores assertion call arguments but without `t *testing.T` argument.
	// E.g [42, result, "helpful comment"].
	Args []ast.Expr
	// ArgsRaw stores assertion call initial arguments.
	// E.g [t, 42, result, "helpful comment"].
	ArgsRaw []ast.Expr
}

func (c CallMeta) String() string {
	return c.SelectorXStr + "." + c.Fn.Name
}

// IsPkg true if this is package (not object) call.
func (c CallMeta) IsPkg() bool {
	return c.Obj == nil
}

// FnMeta stores meta info about assertion function itself, for example "Equal".
type FnMeta struct {
	// Range contains start and end position of function Name.
	analysis.Range
	// Name is a function name.
	Name string
	// NameFTrimmed is a function name without "f" suffix.
	NameFTrimmed string
	// IsFmt is true if function is formatted, e.g. "Equalf".
	IsFmt bool
	// Signature represents assertion signature.
	Signature *types.Signature
}

// NewCallMeta returns meta information about testify assertion call.
// Returns nil if ast.CallExpr is not testify call.
func NewCallMeta(pass *analysis.Pass, ce *ast.CallExpr) *CallMeta {
	funcObj, ok := typeutil.Callee(pass.TypesInfo, ce).(*types.Func)
	if !ok {
		return nil
	}

	sig := funcObj.Type().(*types.Signature) // NOTE(a.telyshev): Func's Type() is always a *Signature.

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || se.Sel == nil {
		return nil
	}
	fnName := se.Sel.Name

	var obj types.Object
	var pkg *types.Package

	if rcv := sig.Recv(); rcv != nil { //nolint:nestif // Types hell.
		if ptr, ok := rcv.Type().(*types.Pointer); ok {
			// Examples:
			// s.Assert         -> method of *suite.Suite        -> package suite ("vendor/github.com/stretchr/testify/suite")
			// s.Assert().Equal -> method of *assert.Assertions  -> package assert ("vendor/github.com/stretchr/testify/assert")
			// s.Equal          -> method of *assert.Assertions  -> package assert ("vendor/github.com/stretchr/testify/assert")
			// reqObj.Falsef    -> method of *require.Assertions -> package require ("vendor/github.com/stretchr/testify/require")
			// collect.Errorf   -> method of *assert.CollectT    -> package assert ("vendor/github.com/stretchr/testify/assert")
			obj = ptr.Elem().(*types.Named).Obj()
			pkg = obj.Pkg()
		}
	} else if id, isIdent := se.X.(*ast.Ident); isIdent {
		if selObj := pass.TypesInfo.ObjectOf(id); selObj != nil {
			if pkgName, ok := selObj.(*types.PkgName); ok {
				// Examples:
				// assert.False      -> assert  -> package assert ("vendor/github.com/stretchr/testify/assert")
				// require.NotEqualf -> require -> package require ("vendor/github.com/stretchr/testify/require")
				pkg = pkgName.Imported()
			}
		}
	}
	if pkg == nil {
		return nil
	}

	isAssert := analysisutil.IsPkg(pkg, testify.AssertPkgName, testify.AssertPkgPath)
	isRequire := analysisutil.IsPkg(pkg, testify.RequirePkgName, testify.RequirePkgPath)
	if !isAssert && !isRequire {
		return nil
	}

	return &CallMeta{
		Call:         ce,
		Range:        ce,
		Obj:          obj,
		IsAssert:     isAssert,
		Selector:     se,
		SelectorXStr: analysisutil.NodeString(pass.Fset, se.X),
		Fn: FnMeta{
			Range:        se.Sel,
			Name:         fnName,
			NameFTrimmed: strings.TrimSuffix(fnName, "f"),
			IsFmt:        strings.HasSuffix(fnName, "f"),
			Signature:    sig,
		},
		Args:    trimTArg(pass, ce.Args),
		ArgsRaw: ce.Args,
	}
}

func trimTArg(pass *analysis.Pass, args []ast.Expr) []ast.Expr {
	if len(args) == 0 {
		return args
	}

	if implementsTestingT(pass, args[0]) {
		return args[1:]
	}
	return args
}
