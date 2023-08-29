package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type BoolCompareTestsGenerator struct{}

func (BoolCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewBoolCompare()
}

func (g BoolCompareTestsGenerator) TemplateData() any {
	var (
		checker        = g.Checker().Name()
		reportUse      = checker + ": use %s.%s"
		reportSimplify = checker + ": need to simplify the assertion"
	)

	type test struct {
		Name              string
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}

	return struct {
		CheckerName CheckerName
		Tests       []test
	}{
		CheckerName: CheckerName(checker),
		Tests: []test{
			{
				Name: "assert.True cases",
				InvalidAssertions: []Assertion{
					{Fn: "Equal", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "Equal", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate == true", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "true == predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate == false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "false == predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate != true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "true != predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate != false", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "false != predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "!predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
				},
				ValidAssertions: []Assertion{
					{Fn: "True", Argsf: "predicate"},
				},
			},
			{
				Name: "assert.False cases",
				InvalidAssertions: []Assertion{
					{Fn: "Equal", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "Equal", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate == true", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "true == predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate == false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "false == predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate != true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "true != predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate != false", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "false != predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "!predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
				},
				ValidAssertions: []Assertion{
					{Fn: "False", Argsf: "predicate"},
				},
			},
		},
	}
}

func (BoolCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("BoolCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(boolCompareTestTmpl))
}

func (BoolCompareTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("BoolCompareTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(boolCompareTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const boolCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var predicate bool
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
