package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type NegatedAssertTestsGenerator struct{}

func (NegatedAssertTestsGenerator) Checker() checkers.Checker {
	return checkers.NewNegatedAssert()
}

func (g NegatedAssertTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": " + checkers.NegatedAssertReport
	)

	return struct {
		CheckerName CheckerName
		Assertions  []Assertion
	}{
		CheckerName: CheckerName(checker),
		Assertions: []Assertion{
			{Fn: "NoError", Argsf: "err", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "NoErrorf", Argsf: `err, "msg"`, ReportMsgf: report, ProposedFn: "NoErrorf", ProposedArgsf: `err, "msg"`},
			{Fn: "Error", Argsf: "err", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
			{Fn: "Equal", Argsf: "42, 42", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "42, 42"},
			{Fn: "True", Argsf: "true", ReportMsgf: report, ProposedFn: "True", ProposedArgsf: "true"},
		},
	}
}

func (NegatedAssertTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("NegatedAssertTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(negatedAssertTestTmpl))
}

func (NegatedAssertTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("NegatedAssertTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(negatedAssertGoldenTmpl))
}

const negatedAssertTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error

	// Invalid.
	{
		{{ range $.Assertions -}}
		if !assert.{{ .Fn }}(t, {{ .Argsf }}) { // want {{ QuoteReport .ReportMsgf }}
			return
		}
		{{ end -}}
	}

	// Valid.
	{
		{{ range $.Assertions -}}
		require.{{ .Fn }}(t, {{ .Argsf }})
		{{ end -}}

		// Positive condition (not negated) — not fixable.
		if assert.NoError(t, err) {
			_ = err
		}
		// Complex body — not fixable.
		if !assert.NoError(t, err) {
			_ = "something"
			return
		}
		// Has else — not fixable.
		if !assert.NoError(t, err) {
			return
		} else {
			_ = err
		}
		// Else-if chain — outer has else, not fixable.
		if !assert.NoError(t, err) {
			return
		} else if !assert.Error(t, err) {
			return
		}
	}
}
`

const negatedAssertGoldenTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error

	// Invalid.
	{
		{{ range $.Assertions -}}
		require.{{ .Fn }}(t, {{ .Argsf }})
		{{ end -}}
	}

	// Valid.
	{
		{{ range $.Assertions -}}
		require.{{ .Fn }}(t, {{ .Argsf }})
		{{ end -}}

		// Positive condition (not negated) — not fixable.
		if assert.NoError(t, err) {
			_ = err
		}
		// Complex body — not fixable.
		if !assert.NoError(t, err) {
			_ = "something"
			return
		}
		// Has else — not fixable.
		if !assert.NoError(t, err) {
			return
		} else {
			_ = err
		}
		// Else-if chain — outer has else, not fixable.
		if !assert.NoError(t, err) {
			return
		} else if !assert.Error(t, err) {
			return
		}
	}
}
`
