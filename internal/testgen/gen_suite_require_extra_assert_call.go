package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteRequireExtraAssertCallTestsGenerator struct{}

func (g SuiteRequireExtraAssertCallTestsGenerator) TemplateData() any {
	var (
		checker = checkers.NewSuiteExtraAssertCall().Name()
		report  = checker + ": use an explicit %s.%s"
	)

	return struct {
		CheckerName CheckerName
		Assrn       Assertion
	}{
		CheckerName: CheckerName(checker),
		Assrn:       Assertion{Fn: "True", Argsf: "b", ReportMsgf: report, ProposedSelector: "suite.Assert()"},
	}
}

func (SuiteRequireExtraAssertCallTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteRequireExtraAssertCallTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteRequireExtraAssertCallTestTmpl))
}

func (SuiteRequireExtraAssertCallTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("SuiteRequireExtraAssertCallTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(suiteRequireExtraAssertCallTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const suiteRequireExtraAssertCallTestTmpl = header + `

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

func (suite *{{ $suiteName }}) TestAll() {
	var b bool
	{{ NewAssertionExpander.FullMode.Expand $.Assrn "suite" "" nil }}
}

func (s *{{ $suiteName }}) TestIgnored() {
	var b bool

	t := s.T()
	assObj, reqObj := assert.New(t), require.New(t)
	assObjS, reqObjS := s.Assert(), s.Require()

	{{ $selectors := arr "assert" "require" "assObj" "reqObj" "assObjS" "reqObjS" "s.Assert()" "s.Require()" }}
	{{ range $si, $sel := $selectors }}
		{{- $t := "" }}{{ if or (eq $sel "assert") (eq $sel "require") }}{{ $t = "t" }}{{ end }}
		{{ NewAssertionExpander.Expand $.Assrn.WithoutReport $sel $t nil }}
	{{- end }}
}
`
