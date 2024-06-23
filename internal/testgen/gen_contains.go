package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ContainsTestsGenerator struct{}

func (ContainsTestsGenerator) Checker() checkers.Checker {
	return checkers.NewContains()
}

func (g ContainsTestsGenerator) TemplateData() any {
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
				Fn:            "True",
				Argsf:         `strings.Contains(a, "abc123")`,
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "Contains",
				ProposedArgsf: `a, "abc123"`,
			}, {
				Fn:            "False",
				Argsf:         `strings.Contains(a, "456")`,
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "NotContains",
				ProposedArgsf: `a, "456"`,
			}, {
				Fn:            "True",
				Argsf:         `strings.Contains(string(b), "abc123")`,
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "Contains",
				ProposedArgsf: `string(b), "abc123"`,
			}, {
				Fn:            "False",
				Argsf:         `strings.Contains(string(b), "456")`,
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "NotContains",
				ProposedArgsf: `string(b), "456"`,
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Contains", Argsf: `a, "abc123"`},
			{Fn: "NotContains", Argsf: `a, "456"`},
			{Fn: "Contains", Argsf: `string(b), "abc123"`},
			{Fn: "NotContains", Argsf: `string(b), "456"`},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Contains", Argsf: `errSentinel.Error(), "user"`},      // Requested by https://github.com/Antonboom/testifylint/issues/47
			{Fn: "Equal", Argsf: `strings.Contains(a, "abc123"), true`}, // Handled by bool-compare
			{Fn: "False", Argsf: `!strings.Contains(a, "abc123")`},      // Handled by bool-compare
			{Fn: "True", Argsf: `!strings.Contains(a, "456")`},          // Handled by bool-compare
		},
	}
}

func (ContainsTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ContainsTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(constainsTestTmpl))
}

func (ContainsTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ContainsTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(constainsTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const constainsTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		a           = "abc123"
		b           = []byte(a)
		errSentinel = errors.New("user not found")
	)

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
