package checkers

import (
	"golang.org/x/tools/go/analysis"
)

// HTTPStatusCode detects situations like
//
//	assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, 200)
//
// and requires
//
//	assert.HTTPStatusCode(t, handler, http.MethodGet, "/index", nil, http.StatusOK)
type HTTPStatusCode struct{}

// NewHTTPStatusCode constructs HTTPStatusCode checker.
func NewHTTPStatusCode() HTTPStatusCode { return HTTPStatusCode{} }
func (HTTPStatusCode) Name() string     { return "http-status-code" }

func (checker HTTPStatusCode) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if call.Fn.NameFTrimmed != "HTTPStatusCode" {
		return nil
	}
	if len(call.Args) < 5 {
		return nil
	}
	valueEdit, importEdit := newHTTPStatusCodeTextEdit(pass, call.Args[4])
	if valueEdit == nil {
		return nil
	}
	textEdits := []analysis.TextEdit{*valueEdit}
	if importEdit != nil {
		textEdits = append(textEdits, *importEdit)
	}
	return newDiagnostic(checker.Name(), call, "use net/http constant instead of integer literal",
		analysis.SuggestedFix{
			Message:   "Replace with net/http constant",
			TextEdits: textEdits,
		})
}
