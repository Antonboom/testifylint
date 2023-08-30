package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorIsTestsGenerator struct{}

func (ErrorIsTestsGenerator) Checker() checkers.Checker {
	return checkers.NewErrorIs()
}

func (g ErrorIsTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{
				Fn:         "Error",
				Argsf:      "err, errSentinel",
				ReportMsgf: checker + ": invalid usage of %[1]s.Error, use %[1]s.%[2]s instead",
				ProposedFn: "ErrorIs",
			},
			{
				Fn:         "NoError",
				Argsf:      "err, errSentinel",
				ReportMsgf: checker + ": invalid usage of %[1]s.NoError, use %[1]s.%[2]s instead",
				ProposedFn: "NotErrorIs",
			},
			{
				Fn:            "True",
				Argsf:         "errors.Is(err, errSentinel)",
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "ErrorIs",
				ProposedArgsf: "err, errSentinel",
			},
			{
				Fn:            "False",
				Argsf:         "errors.Is(err, errSentinel)",
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "NotErrorIs",
				ProposedArgsf: "err, errSentinel",
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Error", Argsf: "err"},
			{Fn: "ErrorIs", Argsf: "err, errSentinel"},
			{Fn: "NoError", Argsf: "err"},
			{Fn: "NotErrorIs", Argsf: "err, errSentinel"},
		},
	}
}

func (ErrorIsTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ErrorIsTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorIsTestTmpl))
}

func (ErrorIsTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ErrorIsTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(errorIsTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const errorIsTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)


func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var errSentinel = errors.New("user not found")
	var err error

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{- if or (eq $assrn.Fn "Error") (eq $assrn.Fn "NoError") }}
				{{ NewAssertionExpander.NotFmtSetMode.Expand $assrn "assert" "t" nil }}
			{{ else }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.FullMode.Expand $assrn "assert" "t" nil }}
		{{ end -}}
	}
}
`
