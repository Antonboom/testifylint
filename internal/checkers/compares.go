package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// Compares detects situations like
//
//	assert.True(t, a == b)
//	assert.True(t, a != b)
//	assert.True(t, a > b)
//	assert.True(t, a >= b)
//	assert.True(t, a < b)
//	assert.True(t, a <= b)
//	assert.False(t, a == b)
//	...
//
// and requires
//
//	assert.Equal(t, a, b)
//	assert.NotEqual(t, a, b)
//	assert.Greater(t, a, b)
//	assert.GreaterOrEqual(t, a, b)
//	assert.Less(t, a, b)
//	assert.LessOrEqual(t, a, b)
//
// If `a` and `b` are pointers then `assert.Same`/`NotSame` is required instead,
// due to the inappropriate recursive nature of `assert.Equal` (based on `reflect.DeepEqual`).
//
// Compares also detects and simplifies equivalent `time.Time` comparisons, like
//
//	assert.True(t, t1.After(t2))
//	assert.Greater(t, t1.Compare(t2), 0)
//	assert.Greater(t, t1, t2)
//
// For `assert.Equal(t, 0, t1.Compare(t2))` case `compares` suggests `WithinDuration` as the assertion with the most
// readable message.
type Compares struct{}

// NewCompares constructs Compares checker.
func NewCompares() Compares   { return Compares{} }
func (Compares) Name() string { return "compares" }

func (checker Compares) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if d := checker.checkBinaryExpr(pass, call); d != nil {
		return d
	}
	return checker.checkTimeCompares(pass, call)
}

func (checker Compares) checkBinaryExpr(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 1 {
		return nil
	}

	be, ok := call.Args[0].(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	var tokenToProposedFn map[token.Token]string

	switch call.Fn.NameFTrimmed {
	case "True":
		tokenToProposedFn = tokenToProposedFnInsteadOfTrue
	case "False":
		tokenToProposedFn = tokenToProposedFnInsteadOfFalse
	default:
		return nil
	}

	proposedFn, ok := tokenToProposedFn[be.Op]
	if !ok {
		return nil
	}

	a, b := be.X, be.Y

	_, xp := isPointer(pass, a)
	_, yp := isPointer(pass, b)
	if xp && yp {
		switch proposedFn {
		case "Equal":
			proposedFn = "Same"
		case "NotEqual":
			proposedFn = "NotSame"
		}
	}

	return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
		analysis.TextEdit{
			Pos:     a.Pos(),
			End:     b.End(),
			NewText: formatAsCallArgs(pass, a, b),
		})
}

var tokenToProposedFnInsteadOfTrue = map[token.Token]string{
	token.EQL: "Equal",
	token.NEQ: "NotEqual",
	token.GTR: "Greater",
	token.GEQ: "GreaterOrEqual",
	token.LSS: "Less",
	token.LEQ: "LessOrEqual",
}

var tokenToProposedFnInsteadOfFalse = map[token.Token]string{
	token.EQL: "NotEqual",
	token.NEQ: "Equal",
	token.GTR: "LessOrEqual",
	token.GEQ: "Less",
	token.LSS: "GreaterOrEqual",
	token.LEQ: "Greater",
}

func (checker Compares) checkTimeCompares(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch fn := call.Fn.NameFTrimmed; fn {
	default:
		return nil

	case "True", "False":
		if len(call.Args) < 1 {
			return nil
		}
		expr := call.Args[0]

		var (
			proposed string
			a, b     ast.Expr
			ok       bool
		)

		if a, b, ok = isTimeMethodCall(pass, expr, "After"); ok { //nolint:nestif // Cannot be simplified.
			if fn == "True" {
				proposed = "Greater"
			} else {
				proposed = "LessOrEqual"
			}
		} else if a, b, ok = isTimeMethodCall(pass, expr, "Before"); ok {
			if fn == "True" {
				proposed = "Less"
			} else {
				proposed = "GreaterOrEqual"
			}
		}

		if proposed != "" {
			return newUseFunctionDiagnostic(checker.Name(), call, proposed, analysis.TextEdit{
				Pos:     expr.Pos(),
				End:     expr.End(),
				NewText: formatAsCallArgs(pass, a, b),
			})
		}

	case "Equal", "EqualValues", "Exactly", "NotEqual", "NotEqualValues",
		"Greater", "GreaterOrEqual", "Less", "LessOrEqual":
	}

	if len(call.Args) < 2 {
		return nil
	}

	var lhs, rhs int
	var cmpArg1, cmpArg2 ast.Expr

	// `t1.Compare(t2), 0`
	if a, b, ok := isTimeMethodCall(pass, call.Args[0], "Compare"); ok {
		n, ok := isIntBasicLit(call.Args[1])
		if !ok {
			return nil
		}
		lhs = timeCompareFn
		rhs = n
		cmpArg1, cmpArg2 = a, b
	}

	// `0, t1.Compare(t2)`
	if a, b, ok := isTimeMethodCall(pass, call.Args[1], "Compare"); ok {
		n, ok := isIntBasicLit(call.Args[0])
		if !ok {
			return nil
		}
		rhs = timeCompareFn
		lhs = n
		cmpArg1, cmpArg2 = a, b
	}

	proposed, ok := timeCompareTransformations[timeAssert{call.Fn.NameFTrimmed, lhs, rhs}]
	if !ok {
		return nil
	}

	var msg string
	if strings.Contains(proposed.argsFmt, ".Equal(") {
		msg = fmt.Sprintf("use %s.Equal", analysisutil.NodeString(pass.Fset, cmpArg1))
	} else {
		f := proposed.fn
		if call.Fn.IsFmt {
			f += "f"
		}
		msg = fmt.Sprintf("use %s.%s", call.SelectorXStr, f)
	}

	argsReplacement := analysis.TextEdit{
		Pos: call.Args[0].Pos(),
		End: call.Args[1].End(),
		NewText: fmt.Appendf(nil, proposed.argsFmt,
			analysisutil.NodeString(pass.Fset, cmpArg1),
			analysisutil.NodeString(pass.Fset, cmpArg2)),
	}
	return newDiagnostic(checker.Name(), call, msg,
		newSuggestedFuncReplacement(call, proposed.fn, argsReplacement))
}

type timeAssert struct {
	fn       string
	lhs, rhs int
}

type timeAssertProposed struct {
	fn      string
	argsFmt string
}

const (
	timeCompareFn = 42 // Alias for (time.Time).Compare.
)

var timeCompareTransformations = map[timeAssert]timeAssertProposed{
	{"Equal", 0, timeCompareFn}:          {"WithinDuration", "%[1]s, %[0]s, 0)"},
	{"EqualValues", 0, timeCompareFn}:    {"WithinDuration", "%[1]s, %[0]s, 0)"},
	{"Exactly", 0, timeCompareFn}:        {"WithinDuration", "%[1]s, %[0]s, 0)"},
	{"NotEqual", 0, timeCompareFn}:       {"False", "%s.Equal(%s)"},
	{"NotEqualValues", 0, timeCompareFn}: {"False", "%s.Equal(%s)"},

	{"Greater", timeCompareFn, 0}:        {"Greater", "%s, %s"},
	{"Less", 0, timeCompareFn}:           {"Greater", "%s, %s"},
	{"GreaterOrEqual", timeCompareFn, 0}: {"GreaterOrEqual", "%s, %s"},
	{"LessOrEqual", 0, timeCompareFn}:    {"GreaterOrEqual", "%s, %s"},
	{"Less", timeCompareFn, 0}:           {"Less", "%s, %s"},
	{"Greater", 0, timeCompareFn}:        {"Less", "%s, %s"},
	{"LessOrEqual", timeCompareFn, 0}:    {"LessOrEqual", "%s, %s"},
	{"GreaterOrEqual", 0, timeCompareFn}: {"LessOrEqual", "%s, %s"},

	{"Equal", timeCompareFn, 1}:     {"Greater", "%s, %s"},
	{"Equal", 1, timeCompareFn}:     {"Greater", "%s, %s"},
	{"NotEqual", timeCompareFn, -1}: {"GreaterOrEqual", "%s, %s"},
	{"NotEqual", -1, timeCompareFn}: {"GreaterOrEqual", "%s, %s"},
	{"Equal", timeCompareFn, -1}:    {"Less", "%s, %s"},
	{"Equal", -1, timeCompareFn}:    {"Less", "%s, %s"},
	{"NotEqual", timeCompareFn, 1}:  {"LessOrEqual", "%s, %s"},
	{"NotEqual", 1, timeCompareFn}:  {"LessOrEqual", "%s, %s"},
}

func isTimeMethodCall(pass *analysis.Pass, e ast.Expr, method string) (a ast.Expr, b ast.Expr, ok bool) {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, nil, false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, nil, false
	}

	if isTimeInstance(pass, se.X) && se.Sel.Name == method && len(ce.Args) == 1 {
		return se.X, ce.Args[0], true
	}
	return nil, nil, false
}
