package main

import (
	"fmt"
	"go/token"
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ComparesTestsGenerator struct{}

func (ComparesTestsGenerator) Checker() checkers.Checker {
	return checkers.NewCompares()
}

func (g ComparesTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	boolOps := []token.Token{token.LAND, token.LOR}
	ignored := make([]Assertion, 0, len(boolOps)*2)
	for _, tok := range boolOps {
		ignored = append(ignored,
			Assertion{Fn: "True", Argsf: fmt.Sprintf("c %s d", tok)},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("d %s c", tok)},
		)
	}

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "True", Argsf: "a == b", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a != b", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a > b", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a >= b", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a < b", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a <= b", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "a, b"},

			{Fn: "False", Argsf: "a == b", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a != b", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a > b", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a >= b", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a < b", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a <= b", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "a, b"},

			{Fn: "True", Argsf: "ptrA == ptrB", ReportMsgf: report, ProposedFn: "Same", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "True", Argsf: "ptrA != ptrB", ReportMsgf: report, ProposedFn: "NotSame", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "False", Argsf: "ptrA == ptrB", ReportMsgf: report, ProposedFn: "NotSame", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "False", Argsf: "ptrA != ptrB", ReportMsgf: report, ProposedFn: "Same", ProposedArgsf: "ptrA, ptrB"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "a, b"},
			{Fn: "NotEqual", Argsf: "a, b"},
			{Fn: "Greater", Argsf: "a, b"},
			{Fn: "GreaterOrEqual", Argsf: "a, b"},
			{Fn: "Less", Argsf: "a, b"},
			{Fn: "LessOrEqual", Argsf: "a, b"},

			{Fn: "Same", Argsf: "ptrA, ptrB"},
			{Fn: "NotSame", Argsf: "ptrA, ptrB"},
		},
		IgnoredAssertions: ignored,
	}
}

func (ComparesTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ComparesTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(comparesTestTmpl))
}

func (ComparesTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ComparesTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(comparesTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const comparesTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var a, b int
	var c, d bool
	var ptrA, ptrB *int

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
