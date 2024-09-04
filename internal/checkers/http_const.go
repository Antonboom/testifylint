package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// HTTPConst detects situations like
//
//	assert.HTTPStatusCode(t, handler, "GET", "/index", nil, 200)
//	assert.HTTPBodyContains(t, handler, "GET", "/index", nil, "counter")
//
// and requires
//
//	assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, http.StatusOK)
//	assert.HTTPBodyContains(t, handler, http.MethodGet, "/index", nil, "counter")
type HTTPConst struct{}

// NewHTTPConst constructs Regexp checker.
func NewHTTPConst() HTTPConst  { return HTTPConst{} }
func (HTTPConst) Name() string { return "http-const" }

func (checker HTTPConst) Check(_ *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	var suggestedFixes []analysis.SuggestedFix
	switch call.Fn.NameFTrimmed {
	case "HTTPBody",
		"HTTPBodyContains",
		"HTTPBodyNotContains",
		"HTTPError",
		"HTTPRedirect",
		"HTTPSuccess":
		if len(call.Args) < 2 {
			return nil
		}
		suggestedFix := newHTTPMethodReplacement(call.Args[1])
		if suggestedFix == nil {
			return nil
		}
		suggestedFixes = append(suggestedFixes, *suggestedFix)
	case "HTTPStatusCode":
		if len(call.Args) < 5 {
			return nil
		}
		if suggestedFix := newHTTPMethodReplacement(call.Args[1]); suggestedFix != nil {
			suggestedFixes = append(suggestedFixes, *suggestedFix)
		}
		if suggestedFix := newHTTPStatusCodeReplacement(call.Args[4]); suggestedFix != nil {
			suggestedFixes = append(suggestedFixes, *suggestedFix)
		}
	}
	if len(suggestedFixes) == 0 {
		return nil
	}
	return newHTTPConstDiagnostic(checker.Name(), call, suggestedFixes)
}

func newHTTPMethodReplacement(e ast.Expr) *analysis.SuggestedFix {
	bt, ok := typeSafeBasicLit(e, token.STRING)
	if !ok {
		return nil
	}
	currentVal, ok := unquoteBasicLitValue(bt)
	if !ok {
		return nil
	}
	key := strings.ToUpper(currentVal)
	newVal, ok := httpMethod[key]
	if !ok {
		return nil
	}
	return newConstReplacement(bt, currentVal, newVal)
}

func newHTTPStatusCodeReplacement(e ast.Expr) *analysis.SuggestedFix {
	bt, ok := typeSafeBasicLit(e, token.INT)
	if !ok {
		return nil
	}
	currentVal := bt.Value
	key := strings.ToUpper(currentVal)
	newVal, ok := httpStatusCode[key]
	if !ok {
		return nil
	}
	return newConstReplacement(bt, currentVal, newVal)
}

func newConstReplacement(bt *ast.BasicLit, currentVal, newVal string) *analysis.SuggestedFix {
	return &analysis.SuggestedFix{
		Message: fmt.Sprintf("Replace %q with %s", currentVal, newVal),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     bt.Pos(),
				End:     bt.End(),
				NewText: []byte(newVal),
			},
		},
	}
}

func newHTTPConstDiagnostic(
	checker string,
	call *CallMeta,
	suggestedFixes []analysis.SuggestedFix,
) *analysis.Diagnostic {
	return newAnalysisDiagnostic(checker, call, "use net/http constants instead of value", suggestedFixes)
}
