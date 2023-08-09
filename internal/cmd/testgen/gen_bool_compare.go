package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type BoolCompareCasesGenerator struct{}

func (BoolCompareCasesGenerator) CheckerName() string {
	return checkers.BoolCompareCheckerName
}

func (BoolCompareCasesGenerator) Data() any {
	const (
		reportUse      = "bool-compare: use %s.%s"
		reportSimplify = "bool-compare: need to simplify the check"
	)

	type test struct {
		Name          string
		InvalidChecks []Check
		ValidChecks   []Check
	}

	return struct {
		Pkgs, Objs     []string
		SuiteSelectors []string
		VarSets        [][]string
		Tests          []test
	}{
		Pkgs:           []string{"assert", "require"},
		Objs:           []string{"assObj", "reqObj"},
		SuiteSelectors: []string{"s", "s.Assert()", "assObj", "s.Require()", "reqObj"},
		VarSets: [][]string{
			{"a"}, {"b.b"}, {"c"}, {"d"}, {"*e"}, {"*f"}, {"g.TheyKilledKenny()"}, {"boolOp()"},
		},
		Tests: []test{
			{
				Name: "True",
				InvalidChecks: []Check{
					{Fn: "Equal", Argsf: "%s, true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "true, %s", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "%s, false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "false, %s", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "%s == true", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "true == %s", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "%s == false", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "false == %s", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "%s != true", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "true != %s", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "%s != false", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "false != %s", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "!%s", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "%s"},
				},
				ValidChecks: []Check{
					{Fn: "True", Argsf: "%s"},
				},
			},
			{
				Name: "False",
				InvalidChecks: []Check{
					{Fn: "Equal", Argsf: "%s, false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "false, %s", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "%s, true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "true, %s", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "%s == true", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "true == %s", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "%s == false", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "false == %s", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "%s != true", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "true != %s", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "%s != false", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "false != %s", ReportMsgf: reportSimplify, ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "!%s", ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "%s"},
				},
				ValidChecks: []Check{
					{Fn: "False", Argsf: "%s"},
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

package mostof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestBoolCompare(t *testing.T) {
	{{- block "vars" . }}
	type withBool struct{ b bool }
	boolOp := func() bool { return true }

	var a bool
	var b withBool
	c := true
	const d = false
	e := new(bool)
	var f *bool
	var g withBoolMethod
	{{- end }}

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $ti, $test := $.Tests }}
		// {{ $test.Name }}.
		{
			{{- range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.InvalidChecks }}
				{{ NewCheckerExpander.Expand $check $pkg $vars }}
				{{ end -}}
			}
			{{ end }}
			// Valid.
			{{ range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.ValidChecks }}
				{{ NewCheckerExpander.Expand $check $pkg $vars }}
				{{ end -}}
			}
			{{ end -}}
		}
		{{ end -}}
	})
	{{ end }}

	assObj, reqObj := assert.New(t), require.New(t)

	{{ range $pi, $obj := $.Objs }}
	t.Run("{{ $obj }}", func(t *testing.T) {
		{{- range $ti, $test := $.Tests }}
		// {{ $test.Name }}.
		{
			{{- range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.InvalidChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
				{{ end -}}
			}
			{{ end }}
			// Valid.
			{{ range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.ValidChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
				{{ end -}}
			}
			{{ end -}}
		}
		{{ end -}}
	})
	{{ end -}}
}

type BoolCompareSuite struct {
	suite.Suite
}

func TestBoolCompareSuite(t *testing.T) {
	suite.Run(t, new(BoolCompareSuite))
}

func (s *BoolCompareSuite) TestAll() {
	{{- template "vars" .}}

	assObj, reqObj := s.Assert(), s.Require()

	{{- range $ti, $test := $.Tests }}
	// {{ $test.Name }}.
	{
		{{- range $si, $sel := $.SuiteSelectors }}
		{
			{{- range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.InvalidChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
				{{ end -}}
			}
			{{ end }}
			// Valid.
			{{ range $vi, $vars := $.VarSets }}
			{
				{{- range $ci, $check := $test.ValidChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
				{{ end -}}
			}
			{{ end -}}
		}
		{{ end -}}
	}
	{{ end -}}
}

type withBoolMethod struct{}

func (withBoolMethod) TheyKilledKenny() bool { return false }
`
