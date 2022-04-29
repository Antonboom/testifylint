package main

import "text/template"

type BoolCompareCasesGenerator struct{}

func (g BoolCompareCasesGenerator) Template() *template.Template {
	return boolCmpCasesTmpl
}

func (g BoolCompareCasesGenerator) Data() any {
	type test struct {
		InvalidChecks []Check
		ValidChecks   []Check
	}

	return struct {
		Pkgs    []string
		VarSets [][]any
		True    test
		False   test
	}{
		Pkgs: []string{"assert", "require"},
		VarSets: [][]any{
			{"a"}, {"b.b"}, {"c"}, {"d"}, {"*e"}, {"*f"}, {"g.TheyKilledKenny()"}, {"boolOp()"},
		},
		True: test{
			InvalidChecks: []Check{
				{Fn: "Equal", ArgsTmpl: "t, %s, true", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "Equal", ArgsTmpl: "t, true, %s", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "NotEqual", ArgsTmpl: "t, %s, false", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "NotEqual", ArgsTmpl: "t, false, %s", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "True", ArgsTmpl: "t, %s == true", ReportedMsg: "need to simplify the check"},
				{Fn: "True", ArgsTmpl: "t, true == %s", ReportedMsg: "need to simplify the check"},
				{Fn: "False", ArgsTmpl: "t, %s == false", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "False", ArgsTmpl: "t, false == %s", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "False", ArgsTmpl: "t, %s != true", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "False", ArgsTmpl: "t, true != %s", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
				{Fn: "True", ArgsTmpl: "t, %s != false", ReportedMsg: "need to simplify the check"},
				{Fn: "True", ArgsTmpl: "t, false != %s", ReportedMsg: "need to simplify the check"},
				{Fn: "False", ArgsTmpl: "t, !%s", ReportedMsgf: "use %s.%s", ProposedFn: "True"},
			},
			ValidChecks: []Check{
				{Fn: "True", ArgsTmpl: "t, %s"},
			},
		},
		False: test{
			InvalidChecks: []Check{
				{Fn: "Equal", ArgsTmpl: "t, %s, false", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "Equal", ArgsTmpl: "t, false, %s", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "NotEqual", ArgsTmpl: "t, %s, true", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "NotEqual", ArgsTmpl: "t, true, %s", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "False", ArgsTmpl: "t, %s == true", ReportedMsg: "need to simplify the check"},
				{Fn: "False", ArgsTmpl: "t, true == %s", ReportedMsg: "need to simplify the check"},
				{Fn: "True", ArgsTmpl: "t, %s == false", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "True", ArgsTmpl: "t, false == %s", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "True", ArgsTmpl: "t, %s != true", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "True", ArgsTmpl: "t, true != %s", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
				{Fn: "False", ArgsTmpl: "t, %s != false", ReportedMsg: "need to simplify the check"},
				{Fn: "False", ArgsTmpl: "t, false != %s", ReportedMsg: "need to simplify the check"},
				{Fn: "True", ArgsTmpl: "t, !%s", ReportedMsgf: "use %s.%s", ProposedFn: "False"},
			},
			ValidChecks: []Check{
				{Fn: "False", ArgsTmpl: "t, %s"},
			},
		},
	}
}

var boolCmpCasesTmpl = template.Must(template.New("boolCmpCasesTmpl").
	Funcs(template.FuncMap{
		"ExpandCheck": ExpandCheck,
	}).
	Parse(`// Code generated by testifylint/internal/cmd/testgen. DO NOT EDIT.

package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoolCompare_True(t *testing.T) {
	type withBool struct{ b bool }
	boolOp := func() bool { return true }

	var a bool
	var b withBool
	c := true
	const d = false
	e := new(bool)
	var f *bool
	var g withBoolMethod
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.True.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.True.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
	})
	{{ end }}}

func TestBoolCompare_False(t *testing.T) {
	type withBool struct{ b bool }
	boolOp := func() bool { return true }

	var a bool
	var b withBool
	c := true
	const d = false
	e := new(bool)
	var f *bool
	var g withBoolMethod
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.False.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.False.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
	})
	{{ end }}}

type withBoolMethod struct{}

func (withBoolMethod) TheyKilledKenny() bool { return false }
`))