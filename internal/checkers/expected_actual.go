package checkers

import (
	"go/ast"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// TODO: упростить регулярку + тесты.
var DefaultExpectedVarPattern = regexp.MustCompile(strings.ReplaceAll(`(
	^exp$|
	^expected$|
	^exp[A-Z0-9].*|
	^expected[A-Z0-9].*|
	.*[a-z0-9]Exp$|
	.*[a-z0-9]Expected$|
	^want$|
	^wanted$|
	^want[A-Z0-9].*|
	^wanted[A-Z0-9].*|
	.*[a-z0-9]Want$|
	.*[a-z0-9]Wanted$)`, "\n\t", ""))

// ExpectedActual checks situation like
//
//	assert.NotEqual(t, result, "expected")
//
// and requires e.g.
//
//	assert.NotEqual(t, "expected", result)
type ExpectedActual struct {
	expPattern *regexp.Regexp
}

// NewExpectedActual constructs ExpectedActual checker using DefaultExpectedVarPattern.
func NewExpectedActual() *ExpectedActual {
	return &ExpectedActual{expPattern: DefaultExpectedVarPattern}
}

func (ExpectedActual) Name() string { return "expected-actual" }

func (checker *ExpectedActual) SetExpPattern(p *regexp.Regexp) *ExpectedActual {
	if p != nil {
		checker.expPattern = p
	}
	return checker
}

func (checker ExpectedActual) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.Name {
	case "Equal", "Equalf", "NotEqual", "NotEqualf",
		"JSONEq", "JSONEqf", "YAMLEq", "YAMLEqf":
	default:
		return nil
	}

	if len(call.Args) < 2 {
		return nil
	}
	first, second := call.Args[0], call.Args[1]

	if isWrongExpectedActualOrder(pass, checker.expPattern, first, second) {
		return newDiagnostic(checker.Name(), call, "need to reverse actual and expected values", &analysis.SuggestedFix{
			Message: "Reverse actual and expected values",
			TextEdits: []analysis.TextEdit{
				{
					Pos:     first.Pos(),
					End:     second.End(),
					NewText: []byte(types.ExprString(second) + ", " + types.ExprString(first)),
				},
			},
		})
	}
	return nil
}

func isWrongExpectedActualOrder(pass *analysis.Pass, expectedVarPattern *regexp.Regexp, _, second ast.Expr) bool {
	// Protection from compile-known constants –
	// results of builtin functions, for example, len().
	if ce, ok := second.(*ast.CallExpr); ok {
		return isCastedBasicLit(ce)
	}

	return isBasicLit(second) ||
		isUntypedConst(pass, second) ||
		isTypedConst(pass, second) ||
		isIdentNamedAsExpected(expectedVarPattern, second) ||
		isStructFieldNamedAsExpected(expectedVarPattern, second)
}

func isCastedBasicLit(ce *ast.CallExpr) bool {
	if len(ce.Args) != 1 {
		return false
	}

	fn, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	switch fn.Name {
	case "uint", "uint8", "uint16", "uint32", "uint64",
		"int", "int8", "int16", "int32", "int64",
		"float32", "float64",
		"rune":
		return isBasicLit(ce.Args[0])

	case "string":
		return isBasicLit(ce.Args[0]) || isIdent(ce.Args[0])

	case "complex64", "complex128":
		return isBasicLit(ce.Args[0]) || isAnyBinaryExpr(ce.Args[0])
	}
	return false
}

func isBasicLit(e ast.Expr) bool {
	_, ok := e.(*ast.BasicLit)
	return ok
}

func isAnyBinaryExpr(e ast.Expr) bool {
	_, ok := e.(*ast.BinaryExpr)
	return ok
}

func isUntypedConst(p *analysis.Pass, e ast.Expr) bool {
	t := p.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}

	b, ok := t.(*types.Basic)
	return ok && b.Info()&types.IsUntyped > 0
}

func isTypedConst(p *analysis.Pass, e ast.Expr) bool {
	tt, ok := p.TypesInfo.Types[e]
	return ok && tt.IsValue() && tt.Value != nil
}

func isIdentNamedAsExpected(pattern *regexp.Regexp, e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && pattern.MatchString(id.Name)
}

func isIdent(e ast.Expr) bool {
	_, ok := e.(*ast.Ident)
	return ok
}

func isStructFieldNamedAsExpected(pattern *regexp.Regexp, e ast.Expr) bool {
	s, ok := e.(*ast.SelectorExpr)
	return ok && isIdentNamedAsExpected(pattern, s.Sel)
}
