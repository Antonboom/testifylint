package checkers

import (
	"go/ast"
	"go/token"
	"go/types"

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
type Compares struct{}

// NewCompares constructs Compares checker.
func NewCompares() Compares   { return Compares{} }
func (Compares) Name() string { return "compares" }

func (checker Compares) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
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

	_, xp := isPointer(pass, be.X)
	_, yp := isPointer(pass, be.Y)
	if xp && yp {
		switch proposedFn {
		case "Equal":
			proposedFn = "Same"
		case "NotEqual":
			proposedFn = "NotSame"
		}
	}

	// When an untyped constant appears in a binary expression, the Go type
	// checker assigns it the contextual type (e.g., math.MaxInt16 gets type
	// int32 when compared with an int32 variable). However, when that same
	// constant is passed as interface{} to a testify function, it reverts to
	// its default type (int). If the default type differs from the other
	// operand's type, testify's comparison functions fail at runtime with
	// "Elements should be the same type".
	//
	// Fix: add an explicit cast for any untyped constant whose default type
	// differs from the contextual type, so both arguments have the same
	// runtime type when passed as interface{}.
	xObjType := untypedConstObjectType(pass, be.X)
	yObjType := untypedConstObjectType(pass, be.Y)
	xCtxType := pass.TypesInfo.TypeOf(be.X)
	yCtxType := pass.TypesInfo.TypeOf(be.Y)

	var aBytes, bBytes []byte
	switch {
	case xObjType != nil && !types.Identical(types.Default(xObjType), types.Default(xCtxType)):
		// X is an untyped constant; cast it to its contextual type.
		aBytes = formatAsCast(pass, be.X, types.Default(xCtxType))
		bBytes = analysisutil.NodeBytes(pass.Fset, be.Y)
	case yObjType != nil && !types.Identical(types.Default(yObjType), types.Default(yCtxType)):
		// Y is an untyped constant; cast it to its contextual type.
		aBytes = analysisutil.NodeBytes(pass.Fset, be.X)
		bBytes = formatAsCast(pass, be.Y, types.Default(yCtxType))
	default:
		aBytes = analysisutil.NodeBytes(pass.Fset, be.X)
		bBytes = analysisutil.NodeBytes(pass.Fset, be.Y)
	}

	return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
		analysis.TextEdit{
			Pos:     be.X.Pos(),
			End:     be.Y.End(),
			NewText: append(append(aBytes, []byte(", ")...), bBytes...),
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
