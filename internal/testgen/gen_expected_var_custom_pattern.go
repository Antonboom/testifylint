package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ExpectedVarCustomPatternTestsGenerator struct{}

func (g ExpectedVarCustomPatternTestsGenerator) TemplateData() any {
	var (
		checker = checkers.NewExpectedActual().Name()
		report  = checker + ": need to reverse actual and expected values"
	)

	return struct {
		CheckerName           CheckerName
		InvalidAssertions     []Assertion
		ValidAssertions       []Assertion
		NotDetectedAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "result, goldenValue", ReportMsgf: report, ProposedArgsf: "goldenValue, result"},
			{Fn: "NotEqual", Argsf: "result, goldenValue", ReportMsgf: report, ProposedArgsf: "goldenValue, result"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "goldenValue, result"},
			{Fn: "NotEqual", Argsf: "goldenValue, result"},
		},
		NotDetectedAssertions: []Assertion{
			{Fn: "Equal", Argsf: "result, expected"},
			{Fn: "Equal", Argsf: "expected, result"},
			{Fn: "NotEqual", Argsf: "result, wanted"},
			{Fn: "NotEqual", Argsf: "wanted, result"},
		},
	}
}

func (ExpectedVarCustomPatternTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ExpectedVarCustomPatternTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(expectedVarCustomPatternTestTmpl))
}

func (ExpectedVarCustomPatternTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ExpectedVarCustomPatternTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(expectedVarCustomPatternTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const expectedVarCustomPatternTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var result any
	var goldenValue, expected, wanted string

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

	// Not detected.
	{
		{{- range $ai, $assrn := $.NotDetectedAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}
`
