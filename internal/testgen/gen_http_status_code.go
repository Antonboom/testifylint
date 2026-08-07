package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type HTTPStatusCodeTestsGenerator struct{}

func (HTTPStatusCodeTestsGenerator) Checker() checkers.Checker {
	return checkers.NewHTTPStatusCode()
}

func (g HTTPStatusCodeTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			// HTTPStatusCode: status code only
			{
				Fn:            "HTTPStatusCode",
				Argsf:         `httpOK, http.MethodGet, "/index", nil, 200`,
				ReportMsgf:    checker + ": use net/http constant instead of integer literal",
				ProposedArgsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`,
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "HTTPStatusCode", Argsf: `httpOK, http.MethodGet, "/index", nil, http.StatusOK`},
		},
		IgnoredAssertions: []Assertion{
			// Uncommon HTTP status codes should be ignored.
			{Fn: "HTTPStatusCode", Argsf: `httpOK, http.MethodGet, "/index", nil, 999`},
		},
	}
}

func (HTTPStatusCodeTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("HTTPStatusCodeTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(httpStatusCodeTestTmpl))
}

func (HTTPStatusCodeTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("HTTPStatusCodeTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(httpStatusCodeTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const httpStatusCodeTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func httpOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
