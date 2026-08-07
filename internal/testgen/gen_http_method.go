package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type HTTPMethodTestsGenerator struct{}

func (HTTPMethodTestsGenerator) Checker() checkers.Checker {
	return checkers.NewHTTPMethod()
}

func (g HTTPMethodTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			// HTTPMethod: method only
			{
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "get", "/index", nil, http.StatusOK`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			},
			{
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "Get", "/index", nil, http.StatusOK`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			},
			{
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, "GET", "/index", nil, http.StatusOK`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			},
			// HTTPBodyContains: method only
			{
				Fn:            "HTTPBodyContains",
				Argsf:         `httpHelloName, "GET", "/", url.Values{"name": []string{"World"}}, "Hello, World!"`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpHelloName, http.MethodGet, "/", url.Values{"name": []string{"World"}}, "Hello, World!"`,
			},
			// HTTPBodyNotContains: method only
			{
				Fn:            "HTTPBodyNotContains",
				Argsf:         `httpHelloName, "POST", "/", nil, "Goodbye"`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpHelloName, http.MethodPost, "/", nil, "Goodbye"`,
			},
			// HTTPError: method only
			{
				Fn:            "HTTPError",
				Argsf:         `httpError, "DELETE", "/resource", nil`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpError, http.MethodDelete, "/resource", nil`,
			},
			// HTTPRedirect: method only
			{
				Fn:            "HTTPRedirect",
				Argsf:         `httpRedirect, "PUT", "/old", nil`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpRedirect, http.MethodPut, "/old", nil`,
			},
			// HTTPSuccess: method only
			{
				Fn:            "HTTPSuccess",
				Argsf:         `httpOK, "PATCH", "/update", nil`,
				ReportMsgf:    checker + ": use net/http constant instead of string literal",
				ProposedArgsf: `httpOK, http.MethodPatch, "/update", nil`,
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "HTTPStatusCode", Argsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`},
			{Fn: "HTTPBodyContains", Argsf: `httpHelloName, http.MethodGet, "/", url.Values{"name": []string{"World"}}, "Hello, World!"`},
			{Fn: "HTTPBodyNotContains", Argsf: `httpHelloName, http.MethodPost, "/", nil, "Goodbye"`},
			{Fn: "HTTPError", Argsf: `httpError, http.MethodDelete, "/resource", nil`},
			{Fn: "HTTPRedirect", Argsf: `httpRedirect, http.MethodPut, "/old", nil`},
			{Fn: "HTTPSuccess", Argsf: `httpOK, http.MethodPatch, "/update", nil`},
		},
		IgnoredAssertions: []Assertion{
			// Uncommon HTTP methods should be ignored.
			{Fn: "HTTPStatusCode", Argsf: `httpOK, "FOO", "/index", nil, http.StatusOK`},
			{Fn: "HTTPBodyContains", Argsf: `httpHelloName, "FOO", "/", url.Values{"name": []string{"World"}}, "Hello, World!"`},
		},
	}
}

func (HTTPMethodTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("HTTPMethodTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(httpMethodTestTmpl))
}

func (HTTPMethodTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("HTTPMethodTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(httpMethodTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const httpMethodTestTmpl = header + `

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

func httpError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error", http.StatusInternalServerError)
}

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/new", http.StatusMovedPermanently)
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
