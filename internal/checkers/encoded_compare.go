package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// EncodedCompare detects situations like
//
//	assert.Equal(t, `{"foo": "bar"}`, body)
//	assert.EqualValues(t, `{"foo": "bar"}`, body)
//	assert.Exactly(t, `{"foo": "bar"}`, body)
//	assert.Equal(t, expectedJSON, resultJSON)
//	assert.Equal(t, expBodyConst, w.Body.String())
//	assert.Equal(t, fmt.Sprintf(`{"value":"%s"}`, hexString), result)
//	assert.Equal(t, "{}", json.RawMessage(resp))
//	assert.Equal(t, expJSON, strings.Trim(string(resultJSONBytes), "\n")) // + Replace, ReplaceAll, TrimSpace
//
//	assert.Equal(t, expectedYML, conf)
//
// and requires
//
//	assert.JSONEq(t, `{"foo": "bar"}`, body)
//	assert.YAMLEq(t, expectedYML, conf)
type EncodedCompare struct{}

// NewEncodedCompare constructs EncodedCompare checker.
func NewEncodedCompare() EncodedCompare { return EncodedCompare{} }
func (EncodedCompare) Name() string     { return "encoded-compare" }

func (checker EncodedCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.NameFTrimmed {
	case "Equal", "EqualValues", "Exactly":
	default:
		return nil
	}

	if len(call.Args) < 2 {
		return nil
	}
	lhs, rhs := call.Args[0], call.Args[1]

	a := checker.unwrap(pass, call.Args[0])
	b := checker.unwrap(pass, call.Args[1])

	var proposed string
	switch {
	case isJSONStyleExpr(pass, a), isJSONStyleExpr(pass, b):
		proposed = "JSONEq"
	case isYAMLStyleExpr(a), isYAMLStyleExpr(b):
		proposed = "YAMLEq"
	}

	if proposed != "" {
		return newUseFunctionDiagnostic(checker.Name(), call, proposed,
			newSuggestedFuncReplacement(call, proposed, analysis.TextEdit{
				Pos:     lhs.Pos(),
				End:     lhs.End(),
				NewText: formatWithStringCastForBytes(pass, a),
			}, analysis.TextEdit{
				Pos:     rhs.Pos(),
				End:     rhs.End(),
				NewText: formatWithStringCastForBytes(pass, b),
			}))
	}
	return nil
}

func (checker EncodedCompare) unwrap(pass *analysis.Pass, e ast.Expr) ast.Expr {
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return e
	}
	if len(ce.Args) == 0 {
		return e
	}

	if isIdentWithName("string", ce.Fun) ||
		isByteArray(ce.Fun) ||
		isStringsReplaceCall(pass, ce) ||
		isStringsReplaceAllCall(pass, ce) ||
		isStringsTrimCall(pass, ce) ||
		isStringsTrimSpaceCall(pass, ce) ||
		isJSONRawMessageCast(pass, ce) {
		return checker.unwrap(pass, ce.Args[0])
	}
	return e
}
