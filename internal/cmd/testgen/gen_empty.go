package main

import (
	"strings"
	"text/template"
)

type EmptyCasesGenerator struct{}

func (g EmptyCasesGenerator) Data() any {
	const (
		report = "empty: use %s.%s"
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
			{"arr"}, {"arrPtr"}, {"sl"}, {"mp"}, {"str"}, {"ch"},
		},
		Tests: []test{
			{
				Name: "Empty",
				InvalidChecks: []Check{
					{Fn: "Len", Argsf: "%s, 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},

					{Fn: "Equal", Argsf: "len(%s), 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "0, len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},

					{Fn: "Less", Argsf: "len(%s), 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Greater", Argsf: "1, len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},

					{Fn: "True", Argsf: "len(%s) == 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "0 == len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "len(%s) < 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "1 > len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},

					{Fn: "False", Argsf: "len(%s) != 0", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "0 != len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "len(%s) >= 1", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "1 <= len(%s)", ReportMsgf: report, ProposedFn: "Empty", ProposedArgsf: "%s"},
				},
				ValidChecks: []Check{
					{Fn: "Empty", Argsf: "%s"},
				},
			},
			{
				Name: "NotEmpty",
				InvalidChecks: []Check{
					{Fn: "NotEqual", Argsf: "len(%s), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "0, len(%s)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					{Fn: "Greater", Argsf: "len(%s), 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "Less", Argsf: "0, len(%s)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					{Fn: "True", Argsf: "len(%s) != 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "0 != len(%s)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "len(%s) > 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "True", Argsf: "0 < len(%s)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					{Fn: "False", Argsf: "len(%s) == 0", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "False", Argsf: "0 == len(%s)", ReportMsgf: report, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
				},
				ValidChecks: []Check{
					{Fn: "NotEmpty", Argsf: "%s"},
				},
			},
		},
	}
}

func (g EmptyCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("EmptyCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(emptyCasesTmplText))
}

func (g EmptyCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("EmptyCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(emptyCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const emptyCasesTmplText = header + `

package mostof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestEmptyAsserts(t *testing.T) {
	{{- block "vars" . }}
	var (
		arr    [0]int
		arrPtr *[0]int
		sl     []int
		mp     map[int]int
		str    string
		ch     chan int
	)
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

type EmptySuite struct {
	suite.Suite
}

func TestEmptySuite(t *testing.T) {
	suite.Run(t, new(EmptySuite))
}

func (s *EmptySuite) TestAll() {
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
`
