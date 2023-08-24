package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteNoExtraAssertCallTestsGenerator struct{}

func (SuiteNoExtraAssertCallTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteNoExtraAssertCall()
}

func (g SuiteNoExtraAssertCallTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": need to simplify the check%.s%.s"
	)

	return struct {
		CheckerName CheckerName
		Assrn       Assertion
	}{
		CheckerName: CheckerName(checker),
		Assrn:       Assertion{Fn: "True", Argsf: "b", ReportMsgf: report, ProposedSelector: "s"},
	}
}

func (SuiteNoExtraAssertCallTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteNoExtraAssertCallTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteNoExtraAssertCallTestTmpl))
}

func (SuiteNoExtraAssertCallTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("SuiteNoExtraAssertCallTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(suiteNoExtraAssertCallTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const suiteNoExtraAssertCallTestTmpl = header + `
package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (s *{{ $suiteName }}) TestAll() {
	var b bool
	{{ NewAssertionExpander.Expand $.Assrn "s.Assert()" "" nil }}
}

func (s *{{ $suiteName }}) TestIgnored() {
	var b bool

	t := s.T()
	assObj, reqObj := assert.New(t), require.New(t)
	assObjS, reqObjS := s.Assert(), s.Require()

	{{ $selectors := arr "assert" "require" "assObj" "reqObj" "assObjS" "reqObjS" "s" "s.Require()" }}
	{{ range $si, $sel := $selectors }}
		{{- $t := "" }}{{ if or (eq $sel "assert") (eq $sel "require") }}{{ $t = "t" }}{{ end }}
		{{ NewAssertionExpander.Expand $.Assrn.WithoutReport $sel $t nil }}
	{{- end }}
}
`
