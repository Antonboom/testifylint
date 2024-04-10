package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type NegativePostiveTestsGenerator struct{}

func (NegativePostiveTestsGenerator) Checker() checkers.Checker {
	return checkers.NewNegativePositive()
}

func (g NegativePostiveTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "Greater", Argsf: "0, a", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
			{Fn: "Greater", Argsf: "a, 0", ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},
			{Fn: "Less", Argsf: "0, a", ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},
			{Fn: "Less", Argsf: "a, 0", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Greater", Argsf: "1, a"},
			{Fn: "Less", Argsf: "1, a"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Greater", Argsf: "1, a"},
			{Fn: "Greater", Argsf: "-1, a"}, // avoid regression
			{Fn: "Greater", Argsf: "0, 0"},  // this one will be reported by useless-assert
			{Fn: "GreaterOrEqual", Argsf: "0, a"},
			{Fn: "GreaterOrEqual", Argsf: "a, 0"},
			{Fn: "Less", Argsf: "1, a"},
			{Fn: "Less", Argsf: "-1, a"}, // avoid regression
			{Fn: "Less", Argsf: "0, 0"},  // this one will be reported by useless-assert
			{Fn: "LessOrEqual", Argsf: "0, a"},
			{Fn: "LessOrEqual", Argsf: "a, 0"},
		},
	}
}

func (NegativePostiveTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("NegativePostiveTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(negativePositiveTestTmpl))
}

func (NegativePostiveTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("NegativePostiveTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(negativePositiveTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const negativePositiveTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var a int

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
