package main

import (
	"strings"
	"text/template"
)

type ErrorIsCasesGenerator struct{}

func (g ErrorIsCasesGenerator) Data() any {
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
			{
				Fn:         "Error",
				Argsf:      "err, errSentinel",
				ReportMsgf: "error-is: invalid usage of %[1]s.Error, use %[1]s.%[2]s instead",
				ProposedFn: "ErrorIs",
			},
			{
				Fn:         "NoError",
				Argsf:      "err, errSentinel",
				ReportMsgf: "error-is: invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead",
				ProposedFn: "NotErrorIs",
			},
		},
		ValidChecks: []Check{
			{Fn: "Error", Argsf: "err"},
			{Fn: "ErrorIs", Argsf: "err, errSentinel"},
			{Fn: "NoError", Argsf: "err"},
			{Fn: "NotErrorIs", Argsf: "err, errSentinel"},
		},
	}
}

func (g ErrorIsCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("ErrorIsCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorIsCasesTmplText))
}

func (g ErrorIsCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("ErrorIsCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(errorIsCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const errorIsCasesTmplText = header + `

package basic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestErrorInsteadOfErrorIs(t *testing.T) {
	{{- block "vars" . }}
	var errSentinel = errors.New("user not found")
	var err error
	{{- end }}

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutFFuncs.Expand $check $pkg nil }}
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

	{{ range $oi, $obj := $.Objs }}
	t.Run("{{ $obj }}", func(t *testing.T) {
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutFFuncs.WithoutTArg.Expand $check $obj nil }}
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

type ErrorInsteadOfErrorIsSuite struct {
	suite.Suite
}

func TestErrorInsteadOfErrorIsSuite(t *testing.T) {
	suite.Run(t, new(ErrorInsteadOfErrorIsSuite))
}

func (s *ErrorInsteadOfErrorIsSuite) TestAll() {
	{{- template "vars" .}}

	assObj, reqObj := s.Assert(), s.Require()

	{{ range $si, $sel := $.SuiteSelectors }}
	{
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutFFuncs.WithoutTArg.Expand $check $sel nil }}
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
}`
