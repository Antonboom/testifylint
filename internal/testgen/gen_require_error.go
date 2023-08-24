package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type RequireErrorTestsGenerator struct{}

func (RequireErrorTestsGenerator) Checker() checkers.Checker {
	return checkers.NewRequireError()
}

func (g RequireErrorTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": for error assertions use the `require` API"
	)

	return struct {
		CheckerName CheckerName
		Assertions  []Assertion
	}{
		CheckerName: CheckerName(checker),
		Assertions: []Assertion{
			{Fn: "Error", Argsf: "err", ReportMsgf: report},
			{Fn: "ErrorIs", Argsf: "err, io.EOF", ReportMsgf: report},
			{Fn: "ErrorAs", Argsf: "err, new(os.PathError)", ReportMsgf: report},
			{Fn: "EqualError", Argsf: `err, "end of file"`, ReportMsgf: report},
			{Fn: "ErrorContains", Argsf: `err, "end of file"`, ReportMsgf: report},
			{Fn: "NoError", Argsf: "err", ReportMsgf: report},
			{Fn: "NotErrorIs", Argsf: "err, io.EOF", ReportMsgf: report},
		},
	}
}

func (RequireErrorTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("RequireErrorTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(requireErrorTestTmpl))
}

func (RequireErrorTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const requireErrorTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error

	assObj, reqObj := assert.New(t), require.New(t)

	// Invalid.
	{{ range $si, $sel := arr "assert" "assObj" }}
		{{- range $ai, $assrn := $.Assertions }}
			{{- $t := "t" }}{{ if eq $sel "assObj"}}{{ $t = "" }}{{ end }}
			{{- NewAssertionExpander.Expand $assrn $sel $t nil }}
		{{ end -}}
	{{- end }}

	// Valid.
	{{ range $si, $sel := arr "require" "reqObj" }}
		{{- range $ai, $assrn := $.Assertions }}
			{{- $t := "t" }}{{ if eq $sel "reqObj"}}{{ $t = "" }}{{ end }}
			{{- NewAssertionExpander.Expand $assrn.WithoutReport $sel $t nil }}
		{{ end -}}
	{{ end -}}
}

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (s *{{ $suiteName }}) TestAll() {
	var err error

	assObj, reqObj := s.Assert(), s.Require()

	// Invalid.
	{{ range $si, $sel := arr "s" "s.Assert()" "assObj" }}
		{{- range $ai, $assrn := $.Assertions }}
			{{- NewAssertionExpander.Expand $assrn $sel "" nil }}
		{{ end -}}
	{{- end }}

	// Valid.
	{{ range $si, $sel := arr "s.Require()" "reqObj" }}
		{{- range $ai, $assrn := $.Assertions }}
			{{- NewAssertionExpander.Expand $assrn.WithoutReport $sel "" nil }}
		{{ end -}}
	{{ end -}}
}
`
