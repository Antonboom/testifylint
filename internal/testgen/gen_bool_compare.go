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
		CheckerName       CheckerName
		Tests             []test
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		Tests: []test{
			{
				Name: "assert.True cases",
				InvalidAssertions: []Assertion{
					{Fn: "Equal", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "Equal", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "EqualValues", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "EqualValues", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "Exactly", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "Exactly", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqualValues", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "NotEqualValues", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate == true", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "true == predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate == false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "false == predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate != true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "true != predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate != false", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "false != predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "!predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: `!result["flag"].(bool)`, ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: `result["flag"].(bool)`}, //nolint:lll
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
					{Fn: "EqualValues", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "EqualValues", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "Exactly", Argsf: "predicate, false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "Exactly", Argsf: "false, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqual", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqualValues", Argsf: "predicate, true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "NotEqualValues", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate == true", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "true == predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate == false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "false == predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "predicate != true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "true != predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "predicate != false", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "False", Argsf: "false != predicate", ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: "!predicate", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate"},
					{Fn: "True", Argsf: `!result["flag"].(bool)`, ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: `result["flag"].(bool)`}, //nolint:lll
				},
				ValidAssertions: []Assertion{
					{Fn: "False", Argsf: "predicate"},
				},
			},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "true, true"},
			{Fn: "Equal", Argsf: "false, false"},
			{Fn: "NotEqual", Argsf: "true, true"},
			{Fn: "NotEqual", Argsf: "false, false"},
			{Fn: "True", Argsf: "true == true"},
			{Fn: "True", Argsf: "false == false"},
			{Fn: "False", Argsf: "true == true"},
			{Fn: "False", Argsf: "false == false"},
			{Fn: "True", Argsf: "true != true"},
			{Fn: "True", Argsf: "false != false"},
			{Fn: "False", Argsf: "true != true"},
			{Fn: "False", Argsf: "false != false"},

			{Fn: "Equal", Argsf: "predicate, predicate"},
			{Fn: "NotEqual", Argsf: "predicate, predicate"},
			{Fn: "True", Argsf: "predicate == predicate"},
			{Fn: "False", Argsf: "predicate == predicate"},
			{Fn: "True", Argsf: "predicate != predicate"},
			{Fn: "False", Argsf: "predicate != predicate"},

			// `any` cases.

			{Fn: "Equal", Argsf: `true, result["flag"]`},
			{Fn: "Equal", Argsf: `result["flag"], true`},
			{Fn: "Equal", Argsf: `false, result["flag"]`},
			{Fn: "Equal", Argsf: `result["flag"], false`},
			{Fn: "NotEqual", Argsf: `true, result["flag"]`},
			{Fn: "NotEqual", Argsf: `result["flag"], true`},
			{Fn: "NotEqual", Argsf: `false, result["flag"]`},
			{Fn: "NotEqual", Argsf: `result["flag"], false`},
			// https://go.dev/ref/spec#Comparison_operators
			// A value x of non-interface type X and a value t of interface type T can be compared
			// if type X is comparable and X implements T.
			{Fn: "True", Argsf: `true == result["flag"]`},
			{Fn: "True", Argsf: `result["flag"] == true`},
			{Fn: "True", Argsf: `false == result["flag"]`},
			{Fn: "True", Argsf: `result["flag"] == false`},
			{Fn: "False", Argsf: `true == result["flag"]`},
			{Fn: "False", Argsf: `result["flag"] == true`},
			{Fn: "False", Argsf: `false == result["flag"]`},
			{Fn: "False", Argsf: `result["flag"] == false`},
			{Fn: "True", Argsf: `true != result["flag"]`},
			{Fn: "True", Argsf: `result["flag"] != true`},
			{Fn: "True", Argsf: `false != result["flag"]`},
			{Fn: "True", Argsf: `result["flag"] != false`},
			{Fn: "False", Argsf: `true != result["flag"]`},
			{Fn: "False", Argsf: `result["flag"] != true`},
			{Fn: "False", Argsf: `false != result["flag"]`},
			{Fn: "False", Argsf: `result["flag"] != false`},

			{Fn: "Equal", Argsf: "foo, foo"},
			{Fn: "NotEqual", Argsf: "foo, foo"},
			{Fn: "True", Argsf: "foo == foo"},
			{Fn: "False", Argsf: "foo == foo"},
			{Fn: "True", Argsf: "foo != foo"},
			{Fn: "False", Argsf: "foo != foo"},
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
	result := map[string]any{}

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

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	var predicate bool
	var foo any
	result := map[string]any{}

	foo = true
	assert.Equal(t, true, foo)

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
