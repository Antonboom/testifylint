package main

import "text/template"

type ErrorCasesGenerator struct{}

func (g ErrorCasesGenerator) Template() *template.Template {
	return errorCasesTmpl
}

func (g ErrorCasesGenerator) Data() any {
	return struct {
		Pkgs           []string
		VarSets        [][]any
		InvalidChecks  []Check
		ValidChecks    []Check
		ValidNilChecks []Check
	}{
		Pkgs: []string{"assert", "require"},
		VarSets: [][]any{
			{"a"}, {"b.Get1()"}, {"c"}, {"errOp()"},
		},
		InvalidChecks: []Check{
			{Fn: "Nil", Argsf: "t, %s", ReportMsgf: "use %s.%s", ProposedFn: "NoError"},
			{Fn: "NotNil", Argsf: "t, %s", ReportMsgf: "use %s.%s", ProposedFn: "Error"},
		},
		ValidChecks: []Check{
			{Fn: "NoError", Argsf: "t, %s"},
			{Fn: "Error", Argsf: "t, %s"},
		},
		ValidNilChecks: []Check{
			{Fn: "Nil", Argsf: "t, ptr"},
			{Fn: "Nil", Argsf: "t, iface"},
			{Fn: "Nil", Argsf: "t, ch"},
			{Fn: "Nil", Argsf: "t, sl"},
			{Fn: "Nil", Argsf: "t, fn"},
			{Fn: "Nil", Argsf: "t, m"},
			{Fn: "Nil", Argsf: "t, uPtr"},

			{Fn: "NotNil", Argsf: "t, ptr"},
			{Fn: "NotNil", Argsf: "t, iface"},
			{Fn: "NotNil", Argsf: "t, ch"},
			{Fn: "NotNil", Argsf: "t, sl"},
			{Fn: "NotNil", Argsf: "t, fn"},
			{Fn: "NotNil", Argsf: "t, m"},
			{Fn: "NotNil", Argsf: "t, uPtr"},
		},
	}
}

var errorCasesTmpl = template.Must(template.New("errorCasesTmpl").
	Funcs(template.FuncMap{
		"ExpandCheck": ExpandCheck,
	}).
	Parse(header + `

package basic

import (
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorAsserts(t *testing.T) {
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
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		{
		{{- range $ci, $check := $.ValidNilChecks }}
			{{ ExpandCheck $check $pkg nil }}
		{{ end }}}
	})
	{{ end }}}

type withErroredMethod struct{}

func (withErroredMethod) Get1() error        { return nil }
func (withErroredMethod) Get2() (int, error) { return 0, nil }
`))
