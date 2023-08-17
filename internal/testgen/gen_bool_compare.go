package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type BoolCompareCasesGenerator struct{}

func (BoolCompareCasesGenerator) CheckerName() string {
	return checkers.NewBoolCompare().Name()
}

func (g BoolCompareCasesGenerator) Data() any {
	var (
		reportUse      = g.CheckerName() + ": use %s.%s"
		reportSimplify = g.CheckerName() + ": need to simplify the check"
	)

	type test struct {
		Name          string
		InvalidChecks []Check
		ValidChecks   []Check
	}

	return struct {
		CheckerName CheckerName
		Tests       []test
	}{
		CheckerName: CheckerName(g.CheckerName()),
		Tests: []test{
			{
				Name: "assert.True cases",
				InvalidChecks: []Check{
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
				ValidChecks: []Check{
					{Fn: "True", Argsf: "predicate"},
				},
			},
			{
				Name: "assert.False cases",
				InvalidChecks: []Check{
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
				ValidChecks: []Check{
					{Fn: "False", Argsf: "predicate"},
				},
			},
		},
	}
}

func (BoolCompareCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("BoolCompareCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(boolCompareCasesTmplText))
}

func (BoolCompareCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("BoolCompareCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(boolCompareCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const boolCompareCasesTmplText = header + `

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
			{{- range $ci, $check := $test.InvalidChecks }}
				{{ NewCheckerExpander.Expand $check "assert" "t" nil }}
			{{- end }}
	
			// Valid.
			{{- range $ci, $check := $test.ValidChecks }}
				{{ NewCheckerExpander.Expand $check "assert" "t" nil }}
			{{- end }}
		}
	{{ end -}}
}
`
