package main

import "text/template"

type ComparisonsCasesGenerator struct{}

func (g ComparisonsCasesGenerator) Template() *template.Template {
	return comparisonsCasesTmpl
}

func (g ComparisonsCasesGenerator) Data() any {
	return struct {
		Pkgs          []string
		InvalidChecks []Check
		ValidChecks   []Check
	}{
		Pkgs: []string{"assert", "require"},
		InvalidChecks: []Check{
			{Fn: "True", Argsf: "t, a == b", ReportMsgf: "use %s.%s", ProposedFn: "Equal"},
			{Fn: "True", Argsf: "t, a != b", ReportMsgf: "use %s.%s", ProposedFn: "NotEqual"},
			{Fn: "True", Argsf: "t, a > b", ReportMsgf: "use %s.%s", ProposedFn: "Greater"},
			{Fn: "True", Argsf: "t, a >= b", ReportMsgf: "use %s.%s", ProposedFn: "GreaterOrEqual"},
			{Fn: "True", Argsf: "t, a < b", ReportMsgf: "use %s.%s", ProposedFn: "Less"},
			{Fn: "True", Argsf: "t, a <= b", ReportMsgf: "use %s.%s", ProposedFn: "LessOrEqual"},

			{Fn: "False", Argsf: "t, a == b", ReportMsgf: "use %s.%s", ProposedFn: "NotEqual"},
			{Fn: "False", Argsf: "t, a != b", ReportMsgf: "use %s.%s", ProposedFn: "Equal"},
			{Fn: "False", Argsf: "t, a > b", ReportMsgf: "use %s.%s", ProposedFn: "LessOrEqual"},
			{Fn: "False", Argsf: "t, a >= b", ReportMsgf: "use %s.%s", ProposedFn: "Less"},
			{Fn: "False", Argsf: "t, a < b", ReportMsgf: "use %s.%s", ProposedFn: "GreaterOrEqual"},
			{Fn: "False", Argsf: "t, a <= b", ReportMsgf: "use %s.%s", ProposedFn: "Greater"},
		},
		ValidChecks: []Check{
			{Fn: "Equal", Argsf: "t, a, b"},
			{Fn: "NotEqual", Argsf: "t, a, b"},
			{Fn: "Greater", Argsf: "t, a, b"},
			{Fn: "GreaterOrEqual", Argsf: "t, a, b"},
			{Fn: "Less", Argsf: "t, a, b"},
			{Fn: "LessOrEqual", Argsf: "t, a, b"},
		},
	}
}

var comparisonsCasesTmpl = template.Must(template.New("comparisonsCasesTmpl").
	Funcs(template.FuncMap{
		"ExpandCheck": ExpandCheck,
	}).
	Parse(header + `

package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComparisons(t *testing.T) {
	var a, b int
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{
			{{- range $ci, $check := $.InvalidChecks }}
			{{ ExpandCheck $check $pkg nil }}
			{{ end }}}

		// Valid {{ $pkg }}s.

		{
			{{- range $ci, $check := $.ValidChecks }}
			{{ ExpandCheck $check $pkg nil }}
			{{ end }}}
	})
	{{ end }}} 
`))
