package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ComparesCasesGenerator struct{}

func (ComparesCasesGenerator) CheckerName() string {
	return checkers.NewCompares().Name()
}

func (ComparesCasesGenerator) Data() any {
	const (
		report = "compares: use %s.%s"
	)

	return struct {
		Pkgs, Objs     []string
		SuiteSelectors []string
		InvalidChecks  []Check
		ValidChecks    []Check
	}{
		Pkgs:           []string{"assert", "require"},
		Objs:           []string{"assObj", "reqObj"},
		SuiteSelectors: []string{"s", "s.Assert()", "assObj", "s.Require()", "reqObj"},
		InvalidChecks: []Check{
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
		},
		ValidChecks: []Check{
			{Fn: "Equal", Argsf: "a, b"},
			{Fn: "NotEqual", Argsf: "a, b"},
			{Fn: "Greater", Argsf: "a, b"},
			{Fn: "GreaterOrEqual", Argsf: "a, b"},
			{Fn: "Less", Argsf: "a, b"},
			{Fn: "LessOrEqual", Argsf: "a, b"},
		},
	}
}

func (ComparesCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("ComparesCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(comparesCasesTmplText))
}

func (ComparesCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("ComparesCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(comparesCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const comparesCasesTmplText = header + `

package {{ .CheckerName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCompares(t *testing.T) {
	{{- block "vars" . }}
	var a, b int
	{{- end }}

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg nil }}
			{{ end -}}
		}

		// Valid.

		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg nil }}
			{{ end -}}
		}
	})
	{{ end }}

	assObj, reqObj := assert.New(t), require.New(t)

	{{ range $pi, $obj := $.Objs }}
	t.Run("{{ $obj }}", func(t *testing.T) {
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj nil }}
			{{ end -}}
		}

		// Valid.

		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj nil }}
			{{ end -}}
		}
	})
	{{ end -}}
}

type ComparesSuite struct {
	suite.Suite
}

func TestComparesSuite(t *testing.T) {
	suite.Run(t, new(ComparesSuite))
}

func (s *ComparesSuite) TestAll() {
	{{- template "vars" .}}

	assObj, reqObj := s.Assert(), s.Require()

	{{- range $si, $sel := $.SuiteSelectors }}
	{
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel nil }}
			{{ end -}}
		}

		// Valid.

		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel nil }}
			{{ end -}}
		}
	}
	{{ end -}}
}
`
