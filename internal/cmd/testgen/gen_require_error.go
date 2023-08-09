package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type RequireErrorCasesGenerator struct{}

func (RequireErrorCasesGenerator) CheckerName() string {
	return checkers.RequireErrorCheckerName
}

func (RequireErrorCasesGenerator) Data() any {
	const (
		report = "require-error: for error assertions use the `require` package"
	)

	type test struct {
		Name           string
		Pkg, Obj       string
		SuiteSelectors []string
		Checks         []Check
	}

	return struct {
		Tests []test
	}{
		Tests: []test{
			{
				Name:           "Asserts",
				Pkg:            "assert",
				Obj:            "assObj",
				SuiteSelectors: []string{"s", "s.Assert()", "assObj"},
				Checks: []Check{
					{Fn: "Error", Argsf: "err", ReportMsgf: report},
					{Fn: "ErrorIs", Argsf: "err, io.EOF", ReportMsgf: report},
					{Fn: "ErrorAs", Argsf: "err, new(os.PathError)", ReportMsgf: report},
					{Fn: "EqualError", Argsf: `err, "end of file"`, ReportMsgf: report},
					{Fn: "ErrorContains", Argsf: `err, "end of file"`, ReportMsgf: report},

					{Fn: "NoError", Argsf: "err", ReportMsgf: report},
					{Fn: "NotErrorIs", Argsf: "err, io.EOF", ReportMsgf: report},
				},
			},
			{
				Name:           "Requires",
				Pkg:            "require",
				Obj:            "reqObj",
				SuiteSelectors: []string{"s.Require()", "reqObj"},
				Checks: []Check{
					{Fn: "Error", Argsf: "err"},
					{Fn: "ErrorIs", Argsf: "err, io.EOF"},
					{Fn: "ErrorAs", Argsf: "err, new(os.PathError)"},
					{Fn: "EqualError", Argsf: `err, "end of file"`},
					{Fn: "ErrorContains", Argsf: `err, "end of file"`},

					{Fn: "NoError", Argsf: "err"},
					{Fn: "NotErrorIs", Argsf: "err, io.EOF"},
				},
			},
		},
	}
}

func (RequireErrorCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("RequireErrorCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(requireErrorCasesTmplText))
}

func (RequireErrorCasesGenerator) GoldenTemplate() *template.Template {
	return nil
}

const requireErrorCasesTmplText = header + `

package requireerror

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRequireError(t *testing.T) {
	var err error

	assObj, reqObj := assert.New(t), require.New(t)
	{{ range $ti, $test := $.Tests }}
		// {{ $test.Name }}.

		{
			{{- range $ci, $check := $test.Checks }}
				{{ NewCheckerExpander.Expand $check $test.Pkg nil }}
			{{ end -}}
		}

		{
			{{- range $ci, $check := $test.Checks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $test.Obj nil }}
			{{ end -}}
		}
	{{ end -}}
}

type RequireErrorSuite struct {
	suite.Suite
}

func TestRequireErrorSuite(t *testing.T) {
	suite.Run(t, new(RequireErrorSuite))
}

func (s *RequireErrorSuite) TestAll() {
	var err error

	assObj, reqObj := s.Assert(), s.Require()

	{{ range $ti, $test := $.Tests }}
		// {{ $test.Name }}.

		{
			{{- range $si, $sel := $test.SuiteSelectors }}
				{{- range $ci, $check := $test.Checks }}
					{{ NewCheckerExpander.WithoutTArg.Expand $check $sel nil }}
				{{ end -}}
			{{ end -}}
		}
	{{ end -}}
}
`
