package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

// BaseTestsGenerator implements tests that cover base code of the analyzer (package inspection, resolving testify objects, etc.).
// In addition, it covers some common features for all checkers, such as formatting diagnostic messages or suggested fixes.
//
// These tests should reduce the combinatorial complexity of the checker tests and their number, since in a good way,
// this code should be duplicated in the tests of each checker.
type BaseTestsGenerator struct{}

func (g BaseTestsGenerator) TemplateData() any {
	reportUse := checkers.NewBoolCompare().Name() + ": use %s.%s"

	return struct {
		Pkgs                  []string
		Objs                  []string
		SuiteSelectors        []string
		SuiteDynamicSelectors []string
		Assertions            []Assertion
	}{
		Pkgs: []string{"assert", "require"},
		Objs: []string{"assert.New(t)", "assertObj", "require.New(t)", "requireObj"},
		SuiteSelectors: []string{
			"s", "s.Assert()", "assertObj", "suiteAssertObj",
			"s.Require()", "requireObj", "suiteRequireObj",
		},
		SuiteDynamicSelectors: []string{"assert.New(%s)", "require.New(%s)"},
		Assertions: []Assertion{
			{Fn: "Equal", Argsf: "true, predicate", ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate"},
			{Fn: "True", Argsf: "predicate"}, // Valid assertion.
		},
	}
}

func (BaseTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("BaseTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(baseTestTmpl))
}

func (BaseTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("BaseTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(baseTestTmpl, "NewAssertionExpander", "NewAssertionExpander.FullMode.AsGolden")))
}

const baseTestTmpl = header + `

package basetests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestBase(t *testing.T) {
	var predicate bool
	{{ block "assertions" . }}
		{{- range $si, $sel := $.Pkgs }}
			{{- range $ai, $assrn := $.Assertions }}
				{{ NewAssertionExpander.FullMode.Expand $assrn $sel "t" nil }}
			{{ end -}}
		{{- end }}

		assertObj, requireObj := assert.New(t), require.New(t)
		{{ range $si, $sel := $.Objs }}
			{{- range $ai, $assrn := $.Assertions }}
				{{ NewAssertionExpander.FullMode.Expand $assrn $sel "" nil }}
			{{ end -}}
		{{- end }}
	{{- end }}

	t.Run("subtest1", func(t *testing.T) {
		{{- template "assertions" . }}

		for range []struct{}{} {
			t.Run("nested test", func(t *testing.T) {
				{{- template "assertions" . -}}
			})
		}
	})

	t.Run("subtest2", func(t *testing.T) {
		{{- template "assertions" . -}}
	})
}

type BaseTestsSuite struct {
	suite.Suite
}

func TestBaseTestsSuite(t *testing.T) {
	suite.Run(t, new(BaseTestsSuite))
}

{{ define "suite-assertions" }}
	{{- $ := index . 0 }}
	{{- $tParam := index . 1 }}

	{{- range $si, $sel := $.Pkgs }}
		{{- range $ai, $assrn := $.Assertions }}
			{{ NewAssertionExpander.FullMode.Expand $assrn $sel $tParam nil }}
		{{ end -}}
	{{- end }}

	assertObj, requireObj := assert.New({{ $tParam }}), require.New({{ $tParam }})
	suiteAssertObj, suiteRequireObj := s.Assert(), s.Require()
	
	{{- range $si, $sel := $.SuiteSelectors }}
		{{- range $ai, $assrn := $.Assertions }}
			{{ NewAssertionExpander.FullMode.Expand $assrn $sel "" nil }}
		{{ end -}}
	{{- end }}

	{{- range $si, $sel := $.SuiteDynamicSelectors }}
		{{- range $ai, $assrn := $.Assertions }}
			{{ $sel := printf $sel $tParam }}
			{{ NewAssertionExpander.FullMode.Expand $assrn $sel "" nil }}
		{{ end -}}
	{{- end }}
{{- end }}

func (s *BaseTestsSuite) TestAll() {
	var predicate bool
	{{ template "suite-assertions" (arr . "s.T()") }}

	s.Run("subtest1", func() {
		{{- template "suite-assertions" (arr . "s.T()") }}

		for range []struct{}{} {
			s.Run("nested test", func() {
				{{- template "suite-assertions" (arr . "s.T()") -}}
			})
		}
	})

	s.Run("subtest2", func() {
		{{- template "suite-assertions" (arr . "s.T()") -}}
	})

	s.T().Run("subtest3", func(t *testing.T) {
		{{- template "suite-assertions" (arr . "t") }}

		for range []struct{}{} {
			s.T().Run("nested test", func(t *testing.T) {
				{{- template "suite-assertions" (arr . "t") -}}
			})
		}
	})

	s.T().Run("subtest4", func(t *testing.T) {
		{{- template "suite-assertions" (arr . "t") -}}
	})
}
`
