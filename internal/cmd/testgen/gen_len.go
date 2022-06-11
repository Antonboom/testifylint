package main

import (
	"strings"
	"text/template"
)

type LenCasesGenerator struct{}

func (g LenCasesGenerator) Data() any {
	const (
		report     = "len: use %s.%s"
		proposedFn = "Len"
	)

	return struct {
		Pkgs, Objs     []string
		SuiteSelectors []string
		VarSets        [][]string
		InvalidChecks  []Check
		ValidChecks    []Check
	}{
		Pkgs:           []string{"assert", "require"},
		Objs:           []string{"assObj", "reqObj"},
		SuiteSelectors: []string{"s", "s.Assert()", "assObj", "s.Require()", "reqObj"},
		VarSets: [][]string{
			{"3"}, {"a"}, {"b.i"}, {"c"}, {"d"}, {"*e"}, {"f.Count()"}, {"intOp()"},
		},
		InvalidChecks: []Check{
			{Fn: "Equal", Argsf: "len(arr), %s", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, %s"},
			{Fn: "Equal", Argsf: "%s, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, %s"},
			{Fn: "True", Argsf: "len(arr) == %s", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, %s"},
			{Fn: "True", Argsf: "%s == len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, %s"},
		},
		ValidChecks: []Check{
			{Fn: "Len", Argsf: "arr, %s"},

			{Fn: "NotEqual", Argsf: "%s, len(arr)"},
			{Fn: "Greater", Argsf: "len(arr), %s"},
			{Fn: "Greater", Argsf: "%s, len(arr)"},
			{Fn: "GreaterOrEqual", Argsf: "len(arr), %s"},
			{Fn: "GreaterOrEqual", Argsf: "%s, len(arr)"},
			{Fn: "Less", Argsf: "len(arr), %s"},
			{Fn: "Less", Argsf: "%s, len(arr)"},
			{Fn: "LessOrEqual", Argsf: "len(arr), %s"},
			{Fn: "LessOrEqual", Argsf: "%s, len(arr)"},

			// `ExpectedActual` checker cases.
			// {Fn: "NotEqual", Argsf: "len(arr), %s"},

			// `Compares` checker cases.
			// {Fn: "True", Argsf: "len(arr) != %s"},
			// {Fn: "False", Argsf: "len(arr) != %s"},
			// ...
		},
	}
}

func (g LenCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("LenCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(lenCasesTmplText))
}

func (g LenCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("LenCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(lenCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const lenCasesTmplText = header + `

package mostof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestLen(t *testing.T) {
	{{- block "vars" . }}
	type withInt struct{ i int }
	intOp := func() int { return 42 }

	var a int
	var b withInt
	c := 1
	const d = 2
	e := new(int)
	var f withIntMethod

	arr := [...]int{1, 2, 3}
	{{- end }}

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg $vars }}
			{{ end -}}
		}
		{{ end -}}
	})
	{{ end }}

	assObj, reqObj := assert.New(t), require.New(t)

	{{ range $pi, $obj := $.Objs }}
	t.Run("{{ $obj }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
			{{ end -}}
		}
		{{ end -}}
	})
	{{ end -}}
}

type LenSuite struct {
	suite.Suite
}

func TestLenSuite(t *testing.T) {
	suite.Run(t, new(LenSuite))
}

func (s *LenSuite) TestAll() {
	{{- template "vars" .}}

	assObj, reqObj := s.Assert(), s.Require()

	{{- range $si, $sel := $.SuiteSelectors }}
	{
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
			{{ end -}}
		}
		{{ end -}}
	}
	{{ end -}}
}

type withIntMethod struct{}

func (withIntMethod) Count() int { return 1 }
`
