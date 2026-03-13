package checkers

import (
	"golang.org/x/tools/go/analysis"
)

// HTTPMethod detects situations like
//
//	assert.HTTPStatusCode(t, handler, "GET", "/index", nil, http.StatusOK)
//	assert.HTTPBodyContains(t, handler, "GET", "/index", nil, "counter")
//
// and requires
//
//	assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, http.StatusOK)
//	assert.HTTPBodyContains(t, handler, http.MethodGet, "/index", nil, "counter")
type HTTPMethod struct{}

// NewHTTPMethod constructs HTTPMethod checker.
func NewHTTPMethod() HTTPMethod { return HTTPMethod{} }
func (HTTPMethod) Name() string { return "http-method" }

func (checker HTTPMethod) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	switch call.Fn.NameFTrimmed {
	case "HTTPBody",
		"HTTPBodyContains",
		"HTTPBodyNotContains",
		"HTTPError",
		"HTTPRedirect",
		"HTTPSuccess",
		"HTTPStatusCode":
		if len(call.Args) < 2 {
			return nil
		}
		valueEdit, importEdit := newHTTPMethodTextEdit(pass, call.Args[1])
		if valueEdit == nil {
			return nil
		}
		textEdits := []analysis.TextEdit{*valueEdit}
		if importEdit != nil {
			textEdits = append(textEdits, *importEdit)
		}
		return newDiagnostic(checker.Name(), call, "use net/http constant instead of string literal",
			analysis.SuggestedFix{
				Message:   "Replace with net/http constant",
				TextEdits: textEdits,
			})
	}
	return nil
}
