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
		CheckerName CheckerName
		LenTest     lenTest
		Tests       []test
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
					{Fn: "Len", Argsf: "elems, 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},

					{Fn: "Equal", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Equal", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},

					{Fn: "Less", Argsf: "len(elems), 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "1, len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},

					{Fn: "True", Argsf: "len(elems) == 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "0 == len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "len(elems) < 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "1 > len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},

					{Fn: "False", Argsf: "len(elems) != 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "False", Argsf: "0 != len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "False", Argsf: "len(elems) >= 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "False", Argsf: "1 <= len(elems)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "elems"},
				},
				ValidAssertions: []Assertion{
					{Fn: "Empty", Argsf: "elems"},
				},
			},
			{
				Name: "assert.NotEmpty cases",
				InvalidAssertions: []Assertion{
					{Fn: "NotEqual", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqual", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},

					{Fn: "Greater", Argsf: "len(elems), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Less", Argsf: "0, len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},

					{Fn: "True", Argsf: "len(elems) != 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "0 != len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "len(elems) > 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "True", Argsf: "0 < len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},

					{Fn: "False", Argsf: "len(elems) == 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "False", Argsf: "0 == len(elems)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
				},
				ValidAssertions: []Assertion{
					{Fn: "NotEmpty", Argsf: "elems"},
				},
			},
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
`
