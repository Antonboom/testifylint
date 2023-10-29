package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorAsTargetTestsGenerator struct{}

func (g ErrorAsTargetTestsGenerator) TemplateData() any {
	var (
		checker        = checkers.NewErrorIsAs().Name()
		defaultReport  = checker + ": second argument to %s.%s must be a non-nil pointer to either a type that implements error, or to any interface type" //nolint:lll
		errorPtrReport = checker + ": second argument to %s.%s should not be *error"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "ErrorAs", Argsf: "err, nil", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, pathErrNotPtr", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, pathErrNil", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, err", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, iface", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, &i", ReportMsgf: defaultReport, ProposedFn: "ErrorAs"},
			{Fn: "ErrorAs", Argsf: "err, &err", ReportMsgf: errorPtrReport, ProposedFn: "ErrorAs"},
		},
		ValidAssertions: []Assertion{
			{Fn: "ErrorAs", Argsf: "err, &pathErr"},
			{Fn: "ErrorAs", Argsf: "err, &iface"},
			{Fn: "ErrorAs", Argsf: "err, emptyIface"},
			{Fn: "ErrorAs", Argsf: "err, &emptyIface"},
		},
	}
}

func (ErrorAsTargetTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ErrorAsTargetTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(ErrorAsTargetTestTmpl))
}

func (ErrorAsTargetTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const ErrorAsTargetTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		err error
		pathErrNotPtr os.PathError
		pathErrNil *os.PathError
		pathErr = new(os.PathError)
		i int
		iface interface { m() }
		emptyIface any
	)

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{ NewAssertionExpander.Expand $assrn "require" "t" nil }}
		{{ end -}}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{ NewAssertionExpander.Expand $assrn "require" "t" nil }}
		{{ end -}}
	}
}
`
