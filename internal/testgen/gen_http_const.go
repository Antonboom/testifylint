package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type HTTPConstTestsGenerator struct{}

func (HTTPConstTestsGenerator) Checker() checkers.Checker {
	return checkers.NewHTTPConst()
}

func (g HTTPConstTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "get", "/index", nil, 200`,
				ReportMsgf:    checker + ": use net/http constants instead of value",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			}, {
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "Get", "/index", nil, 200`,
				ReportMsgf:    checker + ": use net/http constants instead of value",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			}, {
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "GET", "/index", nil, 200`,
				ReportMsgf:    checker + ": use net/http constants instead of value",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			}, {
				Fn:            "HTTPBodyContains",
				Argsf:         `httpHelloName, "GET",  "/", url.Values{"name": []string{"World"}}, "Hello, World!"`,
				ReportMsgf:    checker + ": use net/http constants instead of value",
				ProposedArgsf: `httpHelloName, http.MethodGet,  "/", url.Values{"name": []string{"World"}}, "Hello, World!"`,
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "HTTPStatusCode", Argsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`},
			{Fn: "HTTPBodyContains", Argsf: `httpHelloName, http.MethodGet,  "/", url.Values{"name": []string{"World"}}, "Hello, World!"`},
		},
		IgnoredAssertions: []Assertion{
			// uncommon HTTP methods or HTTP status codes should be ignored
			{Fn: "HTTPStatusCode", Argsf: `httpOK, "FOO", "/index", nil, 999`},
			{Fn: "HTTPBodyContains", Argsf: `httpHelloName, "FOO",  "/", url.Values{"name": []string{"World"}}, "Hello, World!"`},
		},
	}
}

func (HTTPConstTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("HTTPConstTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(httpConstTestTmpl))
}

func (HTTPConstTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("HTTPConstTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(httpConstTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const httpConstTestTmpl = header + `
package {{ .CheckerName.AsPkgName }}
import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func httpOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func httpHelloName(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	_, _ = fmt.Fprintf(w, "Hello, %s!", name)
}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
	// Ignored.
	{
		{{- range $ai, $assrn := $.IgnoredAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}
`
