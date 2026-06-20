package checkers

import (
	"go/ast"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// Zero detects situations like
//
//	assert.Equal(t, time.Time{}, ts)
//	assert.EqualValues(t, time.Time{}, ts)
//	assert.Exactly(t, time.Time{}, ts)
//	assert.True(t, ts.IsZero())
//	assert.True(t, ts.Equal(time.Time{}))
//
//	assert.NotEqual(t, time.Time{}, ts)
//	assert.NotEqualValues(t, time.Time{}, ts)
//	assert.False(t, ts.IsZero())
//	assert.False(t, ts.Equal(time.Time{}))
//
// and requires
//
//	assert.Zero(t, ts)
//	assert.NotZero(t, ts)
//
// Additionally, variables starting with `zero` are considered to have zero values.
type Zero struct{}

// NewZero constructs Zero checker.
func NewZero() Zero       { return Zero{} }
func (Zero) Name() string { return "zero" }

func (checker Zero) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	newUseFnDiagnostic := func(proposed string, survivingArg ast.Expr, replaceStart, replaceEnd token.Pos) *analysis.Diagnostic {
		return newUseFunctionDiagnostic(checker.Name(), call, proposed, analysis.TextEdit{
			Pos:     replaceStart,
			End:     replaceEnd,
			NewText: analysisutil.NodeBytes(pass.Fset, survivingArg),
		})
	}

	switch call.Fn.NameFTrimmed {
	case "Equal", "EqualValues", "Exactly":
		if len(call.Args) < 2 {
			return nil
		}
		lhs, rhs := call.Args[0], call.Args[1]

		f1, f2 := isZeroTimeInstance(pass, lhs), isZeroTimeInstance(pass, rhs)
		if xor(f1, f2) {
			survivingArg, _ := anyVal([]bool{f1, f2}, rhs, lhs)
			return newUseFnDiagnostic("Zero", survivingArg, lhs.Pos(), rhs.End())
		}

	case "NotEqual", "NotEqualValues":
		if len(call.Args) < 2 {
			return nil
		}
		lhs, rhs := call.Args[0], call.Args[1]

		f1, f2 := isZeroTimeInstance(pass, lhs), isZeroTimeInstance(pass, rhs)
		if xor(f1, f2) {
			survivingArg, _ := anyVal([]bool{f1, f2}, rhs, lhs)
			return newUseFnDiagnostic("NotZero", survivingArg, lhs.Pos(), rhs.End())
		}

	case "True":
		if len(call.Args) < 1 {
			return nil
		}
		arg := call.Args[0]

		if t, ok := isTimeIsZeroCall(pass, arg); ok {
			return newUseFnDiagnostic("Zero", t, arg.Pos(), arg.End())
		}
		if t, ok := isTimeEqualZeroCall(pass, arg); ok {
			return newUseFnDiagnostic("Zero", t, arg.Pos(), arg.End())
		}

	case "False":
		if len(call.Args) < 1 {
			return nil
		}
		arg := call.Args[0]

		if t, ok := isTimeIsZeroCall(pass, arg); ok {
			return newUseFnDiagnostic("NotZero", t, arg.Pos(), arg.End())
		}
		if t, ok := isTimeEqualZeroCall(pass, arg); ok {
			return newUseFnDiagnostic("NotZero", t, arg.Pos(), arg.End())
		}
	}
	return nil
}

func isTimeIsZeroCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	if isTimeInstance(pass, se.X) && se.Sel.Name == "IsZero" && len(ce.Args) == 0 {
		return se.X, true
	}
	return nil, false
}

func isTimeEqualZeroCall(pass *analysis.Pass, e ast.Expr) (ast.Expr, bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	if isTimeInstance(pass, se.X) && se.Sel.Name == "Equal" &&
		len(ce.Args) == 1 && isZeroTimeInstance(pass, ce.Args[0]) {
		return se.X, true
	}
	return nil, false
}

var zeroVarPattern = regexp.MustCompile("^zero")

func isZeroTimeInstance(pass *analysis.Pass, e ast.Expr) bool {
	// `var zero time.Time` case.
	if isTimeInstance(pass, e) && isIdentNamedAfterPattern(zeroVarPattern, e) {
		return true
	}

	// `time.Time{}` case.
	l, ok := e.(*ast.CompositeLit)
	return ok && isTimeInstance(pass, l.Type) && len(l.Elts) == 0
}
