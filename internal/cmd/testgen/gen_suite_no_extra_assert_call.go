package main

import (
	"strings"
	"text/template"
)

type SuiteNoExtraAssertCallCasesGenerator struct{}

func (g SuiteNoExtraAssertCallCasesGenerator) Data() any {
	const (
		report = "suite-no-extra-assert-call: need to simplify the check"
	)

	return struct {
		IgnoredSelectors         []string
		IgnoredWithoutTSelectors []string
		IgnoredCheck             Check
		ExtraCall                string
		ExtraCheck               Check
	}{
		IgnoredSelectors:         []string{"assert", "require"},
		IgnoredWithoutTSelectors: []string{"assObj", "reqObj", "assObjS", "reqObjS", "s", "s.Require()"},
		IgnoredCheck:             Check{Fn: "True", Argsf: "b"},
		ExtraCall:                "s.Assert()",
		ExtraCheck:               Check{Fn: "True", Argsf: "b", ReportMsgf: report, ProposedSelector: "s"},
	}
}

func (g SuiteNoExtraAssertCallCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("SuiteNoExtraAssertCallCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteNoExtraAssertCallCasesTmplText))
}

func (g SuiteNoExtraAssertCallCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("SuiteNoExtraAssertCallCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(suiteNoExtraAssertCallCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

const suiteNoExtraAssertCallCasesTmplText = header + `
package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SuiteNoExtraAssertCallSuite struct {
	suite.Suite
}

func TestSuiteNoExtraAssertCallSuite(t *testing.T) {
	suite.Run(t, new(SuiteNoExtraAssertCallSuite))
}

func (s *SuiteNoExtraAssertCallSuite) TestAll() {
	var b bool

	t := s.T()
	assObj, reqObj := assert.New(t), require.New(t)
	assObjS, reqObjS := s.Assert(), s.Require()

	{{ NewCheckerExpander.WithoutTArg.Expand $.ExtraCheck $.ExtraCall nil }}

	// Valid.
	{{ range $si, $sel := $.IgnoredSelectors }}
	{
		{{ NewCheckerExpander.Expand $.IgnoredCheck $sel nil }}
	}
	{{- end }}
	{{ range $si, $sel := $.IgnoredWithoutTSelectors }}
	{
		{{ NewCheckerExpander.WithoutTArg.Expand $.IgnoredCheck $sel nil }}
	}
	{{- end }}
}
`
