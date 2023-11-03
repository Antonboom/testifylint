package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type RequireErrorFnPatternTestsGenerator struct{}

func (g RequireErrorFnPatternTestsGenerator) TemplateData() any {
	var (
		checker = checkers.NewRequireError().Name()
		report  = checker + ": for error assertions use require"
	)

	return struct {
		CheckerName       CheckerName
		ValidAssertions   []Assertion
		InvalidAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		ValidAssertions: []Assertion{
			{Fn: "Error", Argsf: "err"},
			{Fn: "ErrorIs", Argsf: "err, io.EOF"},
			{Fn: "ErrorAs", Argsf: "err, &target"},
			{Fn: "EqualError", Argsf: `err, "end of file"`},
			{Fn: "ErrorContains", Argsf: `err, "end of file"`},
		},
		InvalidAssertions: []Assertion{
			{Fn: "NoError", Argsf: "err", ReportMsgf: report},
			{Fn: "NotErrorIs", Argsf: "err, io.EOF", ReportMsgf: report},
		},
	}
}

func (RequireErrorFnPatternTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("RequireErrorFnPatternTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(requireErrorFnPatternTestTmpl))
}

func (RequireErrorFnPatternTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const requireErrorFnPatternTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error
	var target = new(os.PathError)

	// Invalid.
	{{ range $ai, $assrn := $.InvalidAssertions }}
		{{- if ne $assrn.Fn "NoError" -}}
			{{- NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{ else }}
			{{- NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
			nop()
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" nil }}
		{{ end -}}
	{{ end }}

	// Valid.
	{{ range $ai, $assrn := $.InvalidAssertions }}
		{{- NewAssertionExpander.Expand $assrn.WithoutReport "require" "t" nil }}
	{{ end -}}
	{{ range $si, $sel := arr "assert" "require" }}
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{- NewAssertionExpander.Expand $assrn $sel "t" nil }}
		{{ end -}}
	{{ end -}}
}

func nop() {} // Hack against ignoring of "NoError" sequence.
`
