package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorIsAsTestsGenerator struct{}

func (ErrorIsAsTestsGenerator) Checker() checkers.Checker {
	return checkers.NewErrorIsAs()
}

func (g ErrorIsAsTestsGenerator) TemplateData() any {
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
			{
				Fn:            "True",
				Argsf:         "errors.As(err, &target)",
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "ErrorAs",
				ProposedArgsf: "err, &target",
			},
			{
				Fn:            "False",
				Argsf:         "errors.As(err, &target)",
				ReportMsgf:    checker + ": use %s.%s",
				ProposedFn:    "NotErrorAs",
				ProposedArgsf: "err, &target",
			},

			// Type cases.
			{
				Fn:         "IsType",
				Argsf:      "&http.MaxBytesError{}, err",
				ReportMsgf: checker + ": use %[1]s.ErrorIs or %[1]s.ErrorAs depending on the case",
				ProposedFn: "IsType",
			},
			{
				Fn:         "IsType",
				Argsf:      "(*http.MaxBytesError)(nil), err",
				ReportMsgf: checker + ": use %[1]s.ErrorIs or %[1]s.ErrorAs depending on the case",
				ProposedFn: "IsType",
			},
			{
				Fn:         "IsType",
				Argsf:      "err, os.ErrClosed",
				ReportMsgf: checker + ": use %[1]s.ErrorIs or %[1]s.ErrorAs depending on the case",
				ProposedFn: "IsType",
			},

			{
				Fn:         "IsNotType",
				Argsf:      "&http.MaxBytesError{}, err",
				ReportMsgf: checker + ": use %[1]s.NotErrorIs or %[1]s.NotErrorAs depending on the case",
				ProposedFn: "IsNotType",
			},
			{
				Fn:         "IsNotType",
				Argsf:      "(*http.MaxBytesError)(nil), err",
				ReportMsgf: checker + ": use %[1]s.NotErrorIs or %[1]s.NotErrorAs depending on the case",
				ProposedFn: "IsNotType",
			},
			{
				Fn:         "IsNotType",
				Argsf:      "err, os.ErrClosed",
				ReportMsgf: checker + ": use %[1]s.NotErrorIs or %[1]s.NotErrorAs depending on the case",
				ProposedFn: "IsNotType",
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Error", Argsf: "err"},
			{Fn: "NoError", Argsf: "err"},

			{Fn: "ErrorIs", Argsf: "err, errSentinel"},
			{Fn: "NotErrorIs", Argsf: "err, errSentinel"},

			{Fn: "ErrorAs", Argsf: "err, &target"},
			{Fn: "ErrorAs", Argsf: "err, new(*http.MaxBytesError)"},
			{Fn: "NotErrorAs", Argsf: "err, &target"},
			{Fn: "NotErrorAs", Argsf: "err, new(*http.MaxBytesError)"},
		},
	}
}

func (ErrorIsAsTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ErrorIsAsTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorIsAsTestTmpl))
}

func (ErrorIsAsTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ErrorIsAsTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(errorIsAsTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const errorIsAsTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		errSentinel = errors.New("user not found") 
		err error
		target = new(os.PathError)
	)

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{/* F-assertions don't allow error as the second arg, because string required. */}}
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
