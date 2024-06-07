package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteDontUsePkgTestsGenerator struct{}

func (SuiteDontUsePkgTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteDontUsePkg()
}

func (g SuiteDontUsePkgTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName  CheckerName
		AssertAssrn  Assertion
		RequireAssrn Assertion
	}{
		CheckerName:  CheckerName(checker),
		AssertAssrn:  Assertion{Fn: "Equal", Argsf: "42, result", ReportMsgf: report, ProposedSelector: "s"},
		RequireAssrn: Assertion{Fn: "Equal", Argsf: "42, result", ReportMsgf: report, ProposedSelector: "s.Require()"},
	}
}

func (SuiteDontUsePkgTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteDontUsePkgTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteDontUsePkgTestTmpl))
}

func (SuiteDontUsePkgTestsGenerator) GoldenTemplate() Executor {
	golden := strings.ReplaceAll(suiteDontUsePkgTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")
	golden = strings.Replace(golden, "s.T()", "", 4)
	return template.Must(template.New("SuiteDontUsePkgTestsGenerator.GoldenTemplate").Funcs(fm).Parse(golden))
}

const suiteDontUsePkgTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"

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
	var result any
	assObj, reqObj := s.Assert(), s.Require()

	{{ range $si, $sel := arr "s" "s.Assert()" "s.Require()" "assObj" "reqObj" }}
		{{ NewAssertionExpander.Expand $.AssertAssrn.WithoutReport $sel "" nil }}
	{{- end }}

	{{ NewAssertionExpander.Expand $.AssertAssrn "assert" "s.T()" nil }}
	{{ NewAssertionExpander.Expand $.AssertAssrn "a" "s.T()" nil }}

	{{ NewAssertionExpander.Expand $.RequireAssrn "require" "s.T()" nil }}
	{{ NewAssertionExpander.Expand $.RequireAssrn "r" "s.T()" nil }}

	s.T().Run("not detected in order to avoid conflict with suite-subtest-run", func(t *testing.T) {
		{{ template "pkg-assertions" . }}
	})
}

func {{ .CheckerName.AsTestName }}_NoSuiteNoProblem(t *testing.T) {
	{{ block "pkg-assertions" . -}}
		var result any
		assObj, reqObj := assert.New(t), require.New(t)

		{{ range $si, $sel := arr "assert" "assObj" "require" "reqObj" }}
			{{- $t := "t" }}{{ if or (eq $sel "assObj") (eq $sel "reqObj") }}{{ $t = "" }}{{ end }}
			{{ NewAssertionExpander.Expand $.AssertAssrn.WithoutReport $sel $t nil }}
		{{- end }}
	{{- end }}
}
`
