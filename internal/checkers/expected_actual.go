package checkers

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"regexp"
	"strings"
)

var defaultExpectedVarPattern = regexp.MustCompile(strings.ReplaceAll(`(
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

func ExpectedActual(expPattern *regexp.Regexp) Checker {
	if expPattern == nil {
		expPattern = defaultExpectedVarPattern
	}

	return func(pass *analysis.Pass, fn FnMeta) {
		switch fn.Name {
		case "Equal", "Equalf", "NotEqual", "NotEqualf",
			"JSONEq", "JSONEqf", "YAMLEq", "YAMLEqf":
		default:
			return
		}

		if len(fn.Args) < 3 {
			return
		}

		if isWrongExpectedActualOrder(pass, expPattern, fn.Args[1], fn.Args[2]) {
			r.Report(pass, fn, "need to reverse actual and expected values")
		}
	}
}

func isWrongExpectedActualOrder(pass *analysis.Pass, expectedVarPattern *regexp.Regexp, _, second ast.Expr) bool {
	if ce, ok := second.(*ast.CallExpr); ok { // Protection from compile-known constants results of builtin functions, for example, len().
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
	case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"float32", "float64", "rune":
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
