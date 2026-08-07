package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Regexp detects situations like
//
//	assert.Regexp(t, regexp.MustCompile(`\[.*\] DEBUG \(.*TestNew.*\): message`), out)
//	assert.NotRegexp(t, regexp.MustCompile(`\[.*\] TRACE message`), out)
//
// and requires
//
//	assert.Regexp(t, `\[.*\] DEBUG \(.*TestNew.*\): message`, out)
//	assert.NotRegexp(t, `\[.*\] TRACE message`, out)
//
// Also detects situations like
//
//	assert.Regexp(t, []byte(`\w+`), str)
//
// and requires
//
//	assert.Regexp(t, `\w+`, str) // or *regexp.Regexp
type Regexp struct{}

// NewRegexp constructs Regexp checker.
func NewRegexp() Regexp     { return Regexp{} }
func (Regexp) Name() string { return "regexp" }

func (checker Regexp) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.NameFTrimmed {
	default:
		return nil
	case "Regexp", "NotRegexp":
	}

	if len(call.Args) < 1 {
		return nil
	}

	arg := call.Args[0]

	ce, ok := arg.(*ast.CallExpr)
	if ok && len(ce.Args) == 1 && isRegexpMustCompileCall(pass, ce) {
		return newRemoveMustCompileDiagnostic(pass, checker.Name(), call, ce, ce.Args[0])
	}

	if !isStringOrRegexpType(pass, arg) {
		return newDiagnostic(checker.Name(), call,
			"use string or *regexp.Regexp as the first argument")
	}
	return nil
}

// isStringOrRegexpType returns true if the expression is of strict string type
// (including untyped strings and aliases to string), *regexp.Regexp (including aliases to
// *regexp.Regexp), or a type parameter accepted conservatively to avoid false positives for code
// like `func[T interface{ ~string }](rx T) { assert.Regexp(t, rx, str) }`.
//
// Note that *types.TypeParam and *types.Alias are different cases and are handled separately.
func isStringOrRegexpType(pass *analysis.Pass, e ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}

	unaliased := types.Unalias(t)

	// Conservatively accept type parameters to avoid false positives for code like
	// `func[T interface{ ~string }](rx T) { assert.Regexp(t, rx, str) }`.
	if _, ok := unaliased.(*types.TypeParam); ok {
		return true
	}

	// Accept only strict string and untyped string, plus aliases to those exact types.
	basic, ok := unaliased.(*types.Basic)
	if ok && (basic.Kind() == types.String || basic.Kind() == types.UntypedString) {
		return true
	}

	// Accept *regexp.Regexp and aliases to that exact pointer type.
	ptr, ok := unaliased.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "regexp" && obj.Name() == "Regexp"
}
