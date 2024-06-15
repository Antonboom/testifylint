package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type EmptyTestsGenerator struct{}

func (EmptyTestsGenerator) Checker() checkers.Checker {
	return checkers.NewEmpty()
}

func (g EmptyTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	type test struct {
		Name              string
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}

	type lenTest struct {
		Vars  []string
		Assrn Assertion
	}

	return struct {
		CheckerName       CheckerName
		LenTest           lenTest
		Tests             []test
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		LenTest: lenTest{
			Vars:  []string{"arr", "arrPtr", "sl", "mp", "str", "ch"},
			Assrn: Assertion{Fn: "Equal", Argsf: "0, len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
		},
		Tests: []test{
			{
				Name: "assert.Empty cases",
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n == 0, n <= 0, n < 1
					// 0 == n, 0 >= n, 1 > n
					{Fn: "Len", Argsf: "elems, 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Equal", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Equal", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "EqualValues", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "EqualValues", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Exactly", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Exactly", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "LessOrEqual", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "GreaterOrEqual", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Less", Argsf: "len(elems), 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "1, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Zero", Argsf: "len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Empty", Argsf: "len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},

					// Bullshit, but supported by the checker:
					// n < 0, n <= 0
					// 0 > n, 0 >= n
					{Fn: "Less", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "LessOrEqual", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "GreaterOrEqual", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
				},
				ValidAssertions: []Assertion{
					{Fn: "Empty", Argsf: "elems"},
				},
			},
			{
				Name: "assert.NotEmpty cases",
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n != 0, n > 0
					// 0 != n, 0 < n
					{Fn: "NotEqual", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqual", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqualValues", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqualValues", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Less", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Positive", Argsf: "len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotZero", Argsf: "len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEmpty", Argsf: "len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
				},
				ValidAssertions: []Assertion{
					{Fn: "NotEmpty", Argsf: "elems"},
				},
			},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Len", Argsf: "elems, len(elems)"},
			{Fn: "Len", Argsf: "elems, 1"},

			{Fn: "Equal", Argsf: "len(elems), len(elems)"},
			{Fn: "Equal", Argsf: "len(elems), 1"},
			{Fn: "Equal", Argsf: "1, len(elems)"},

			{Fn: "NotEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "NotEqual", Argsf: "len(elems), 1"},
			{Fn: "NotEqual", Argsf: "1, len(elems)"},

			{Fn: "Greater", Argsf: "len(elems), len(elems)"},
			{Fn: "Greater", Argsf: "len(elems), 2"},
			{Fn: "Greater", Argsf: "2, len(elems)"},

			{Fn: "GreaterOrEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 0"},
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 2"},
			{Fn: "GreaterOrEqual", Argsf: "2, len(elems)"},

			{Fn: "Less", Argsf: "len(elems), len(elems)"},
			{Fn: "Less", Argsf: "len(elems), 2"},
			{Fn: "Less", Argsf: "2, len(elems)"},

			{Fn: "LessOrEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "LessOrEqual", Argsf: "0, len(elems)"},
			{Fn: "LessOrEqual", Argsf: "len(elems), 2"},
			{Fn: "LessOrEqual", Argsf: "2, len(elems)"},

			// The linter ignores n > 1 case, because it is not exactly equivalent of NotEmpty.
			{Fn: "Greater", Argsf: "len(elems), 1"},
			{Fn: "Less", Argsf: "1, len(elems)"},
			// The linter ignores n >= 1 case, because NotEmpty in such case may impair the readability of the test.
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 1"},
			{Fn: "LessOrEqual", Argsf: "1, len(elems)"},
		},
	}
}

func (EmptyTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("EmptyTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(emptyTestTmpl))
}

func (EmptyTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("EmptyTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(emptyTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const emptyTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var elems []string
	{{ range $ti, $test := $.Tests }}
		// {{ $test.Name }}.
		{
			// Invalid.
			{{- range $ai, $assrn := $test.InvalidAssertions }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}

			// Valid.
			{{- range $ai, $assrn := $test.ValidAssertions }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}
		}
	{{ end -}}
}

func {{ .CheckerName.AsTestName }}_LenVarIndependence(t *testing.T) {
	var (
		arr    [0]int
		arrPtr *[0]int
		sl     []int
		mp     map[int]int
		str    string
		ch     chan int
	)
	{{ range $vi, $var := $.LenTest.Vars }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $.LenTest.Assrn "assert" "t" (arr $var) }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	var elems []string

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
