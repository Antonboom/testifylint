package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorCasesGenerator struct{}

func (ErrorCasesGenerator) CheckerName() string {
	return checkers.ErrorCheckerName
}

func (ErrorCasesGenerator) Data() any {
	return struct {
		Pkgs, Objs     []string
		SuiteSelectors []string
		VarSets        [][]string
		InvalidChecks  []Check
		ValidChecks    []Check
		ValidNilChecks []Check
	}{
		Pkgs:           []string{"assert", "require"},
		Objs:           []string{"assObj", "reqObj"},
		SuiteSelectors: []string{"s", "s.Assert()", "assObj", "s.Require()", "reqObj"},
		VarSets: [][]string{
			{"a"}, {"b.Get1()"}, {"c"}, {"errOp()"},
		},
		InvalidChecks: []Check{
			{Fn: "Nil", Argsf: "%s", ReportMsgf: "error: use %s.%s", ProposedFn: "NoError"},
			{Fn: "NotNil", Argsf: "%s", ReportMsgf: "error: use %s.%s", ProposedFn: "Error"},
		},
		ValidChecks: []Check{
			{Fn: "NoError", Argsf: "%s"},
			{Fn: "Error", Argsf: "%s"},
		},
		ValidNilChecks: []Check{
			{Fn: "Nil", Argsf: "ptr"},
			{Fn: "Nil", Argsf: "iface"},
			{Fn: "Nil", Argsf: "ch"},
			{Fn: "Nil", Argsf: "sl"},
			{Fn: "Nil", Argsf: "fn"},
			{Fn: "Nil", Argsf: "m"},
			{Fn: "Nil", Argsf: "uPtr"},

			{Fn: "NotNil", Argsf: "ptr"},
			{Fn: "NotNil", Argsf: "iface"},
			{Fn: "NotNil", Argsf: "ch"},
			{Fn: "NotNil", Argsf: "sl"},
			{Fn: "NotNil", Argsf: "fn"},
			{Fn: "NotNil", Argsf: "m"},
			{Fn: "NotNil", Argsf: "uPtr"},
		},
	}
}

func (ErrorCasesGenerator) ErroredTemplate() *template.Template {
	return template.Must(template.New("ErrorCasesGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorCasesTmplText))
}

func (ErrorCasesGenerator) GoldenTemplate() *template.Template {
	return template.Must(template.New("ErrorCasesGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(errorCasesTmplText, "NewCheckerExpander", "NewCheckerExpander.AsGolden")))
}

// todo: исправить шаблон
const errorCasesTmplText = header + `

package mostof

import (
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestErrorAsserts(t *testing.T) {
	{{- block "vars" . }}
	errOp := func() error { return io.EOF }
	
	var (
		ptr   *int
		iface any
		ch    chan error
		sl    []error
		fn    func()
		m     map[int]int
		uPtr unsafe.Pointer
	)
	
	var a error
	var b withErroredMethod
	_, c := b.Get2()
	{{- end }}

	{{ range $pi, $pkg := $.Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.Expand $check $pkg $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid nils.

		{
			{{- range $ci, $check := $.ValidNilChecks }}
				{{ NewCheckerExpander.Expand $check $pkg nil }}
			{{ end -}}
		}
	})
	{{ end }}

	assObj, reqObj := assert.New(t), require.New(t)

	{{ range $oi, $obj := $.Objs }}
	t.Run("{{ $obj }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $obj $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid nils.

		{
			{{- range $ci, $check := $.ValidNilChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $obj nil }}
			{{ end -}}
		}
	})
	{{ end -}}
}

type ErrorSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestAll() {
	{{- template "vars" .}}

	assObj, reqObj := s.Assert(), s.Require()

	{{ range $si, $sel := $.SuiteSelectors }}
	{
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ NewCheckerExpander.WithoutTArg.Expand $check $sel $vars }}
			{{ end -}}
		}
		{{ end }}
		// Valid nils.

		{
			{{- range $ci, $check := $.ValidNilChecks }}
				{{ NewCheckerExpander.WithoutTArg.Expand $check $sel nil }}
			{{ end -}}
		}
	}
	{{ end -}}
}

type withErroredMethod struct{}

func (withErroredMethod) Get1() error        { return nil }
func (withErroredMethod) Get2() (int, error) { return 0, nil }
`
