package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ExpectedActualTestsGenerator struct{}

func (ExpectedActualTestsGenerator) Checker() checkers.Checker {
	return checkers.NewExpectedActual()
}

func (g ExpectedActualTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": need to reverse actual and expected values"
	)

	type test struct {
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}

	type literal struct {
		TypeStr string
		Value   string
	}

	return struct {
		CheckerName    CheckerName
		Typed          []literal
		Untyped        []string
		ExpVars        []string
		Basic          test
		Strings        test
		NotDetected    []Assertion
		RealLifeJSONEq Assertion
		Ignored        []Assertion
	}{
		CheckerName: CheckerName(checker),
		Typed: []literal{
			{TypeStr: "bool", Value: "true"},

			{TypeStr: "uint", Value: "10"},
			{TypeStr: "uint8", Value: "11"},
			{TypeStr: "uint16", Value: "12"},
			{TypeStr: "uint32", Value: "13"},
			{TypeStr: "uint64", Value: "14"},

			{TypeStr: "int", Value: "20"},
			{TypeStr: "int8", Value: "21"},
			{TypeStr: "int16", Value: "22"},
			{TypeStr: "int32", Value: "23"},
			{TypeStr: "int64", Value: "14"},

			{TypeStr: "float32", Value: "30."},
			{TypeStr: "float64", Value: "31."},

			{TypeStr: "complex64", Value: "40i"},
			{TypeStr: "complex128", Value: "41i"},

			{TypeStr: "string", Value: `"50"`},
			{TypeStr: "uintptr", Value: "60"},
			{TypeStr: "byte", Value: "70"},
			{TypeStr: "rune", Value: `'\x80'`},

			{TypeStr: "OwnInt", Value: "90"},
			{TypeStr: "OwnString", Value: `"91"`},
		},
		Untyped: []string{
			"42",
			"3.14",
			"0.707i",
			`"raw string"`,
			`'\U00101234'`,
			"true",
		},
		ExpVars: []string{
			"expected",
			"tt.expected",
			"tt.exp()",
			"expectedVal()",
			"[]int{1, 2, 3}",
			"[3]int{1, 2, 3}",
			`map[string]int{"0": 1}`,
		},
		Basic: test{
			InvalidAssertions: []Assertion{
				{Fn: "Equal", Argsf: "result, %s", ReportMsgf: report, ProposedArgsf: "%s, result"},
				{Fn: "NotEqual", Argsf: "result, %s", ReportMsgf: report, ProposedArgsf: "%s, result"},
			},
			ValidAssertions: []Assertion{
				{Fn: "Equal", Argsf: "%s, result"},
				{Fn: "NotEqual", Argsf: "%s, result"},
			},
		},
		Strings: test{
			InvalidAssertions: []Assertion{
				{Fn: "JSONEq", Argsf: "result, %s", ReportMsgf: report, ProposedArgsf: "%s, result"},
				{Fn: "YAMLEq", Argsf: "result, %s", ReportMsgf: report, ProposedArgsf: "%s, result"},
			},
			ValidAssertions: []Assertion{
				{Fn: "JSONEq", Argsf: "%s, result"},
				{Fn: "YAMLEq", Argsf: "%s, result"},
			},
		},
		NotDetected: []Assertion{
			{Fn: "Equal", Argsf: "result, %s"},    // Invalid order, but no warning.
			{Fn: "NotEqual", Argsf: "result, %s"}, // Invalid order, but no warning.
		},
		RealLifeJSONEq: Assertion{
			Fn:            "JSONEq",
			Argsf:         "string(body), string(expectedJSON)",
			ReportMsgf:    report,
			ProposedArgsf: "string(expectedJSON), string(body)",
		},
		Ignored: []Assertion{
			{Fn: "Equal", Argsf: `"value", "value"`},
			{Fn: "Equal", Argsf: "expected, expected"},
			{Fn: "Equal", Argsf: "[]int{1, 2}, map[int]int{1: 2}"},
			{Fn: "NotEqual", Argsf: "result, result"},
		},
	}
}

func (ExpectedActualTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ExpectedActualTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(expectedActualTestTmpl))
}

func (ExpectedActualTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ExpectedActualTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(expectedActualTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const expectedActualTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct { expected string } //
func (c testCase) exp() string { return c.expected }

{{ define "var-tests" }}
	{{- $ := index . 0 }}
	{{- $Assertions := index . 1 }}

	{{- range $ai, $assrn := $Assertions }}
		{{- range $vi, $var := $.ExpVars }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
		{{- end }}
	
		{{- range $vi, $var := $.Untyped }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
		{{- end }}
	{{ end -}}
{{ end }}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var expected string
	var tt testCase
	expectedVal := func() any { return nil }

	var result any

	// Invalid.
	{
		{{- template "var-tests" arr . $.Basic.InvalidAssertions -}}
	}

	// Valid.
	{
		{{- template "var-tests" arr . $.Basic.ValidAssertions -}}
	}
}

{{ define "const-tests" }}
	{{- $ := index . 0 }}
	{{- $Assertions := index . 1 }}

	{{- range $ai, $assrn := $Assertions }}
		{{- range $vi, $var := $.Typed }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "tc%d" $vi)) }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "tc%dcasted" $vi)) }}
		{{- end }}
	
		{{- range $vi, $var := $.Untyped }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "uc%d" $vi)) }}
		{{- end }}
	{{- end }}
{{ end }}

func {{ .CheckerName.AsTestName }}_DetectConsts(t *testing.T) {
	type OwnInt int
	type OwnString string

	const (
		{{- range $li, $l := $.Typed }}
			tc{{ $li }} {{ $l.TypeStr }} = {{ $l.Value }}
			tc{{ $li }}casted = {{ $l.TypeStr }}({{ $l.Value }})
		{{- end }}
	)

	const (
		{{- range $li, $l := $.Untyped }}
			uc{{ $li }} = {{ $l }}
		{{- end }}
	)

	var result any

	// Invalid.
	{
		{{- template "const-tests" arr . .Basic.InvalidAssertions -}}
	}

	// Valid.
	{
		{{- template "const-tests" arr . .Basic.ValidAssertions -}}
	}
}

func {{ .CheckerName.AsTestName }}_Strings(t *testing.T) {
	const data = "{}"
	var result string

	// Invalid.
	{
		{{- range $ai, $assrn := $.Strings.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "data") }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.Strings.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "data") }}
		{{- end }}
	}

	// Real life case.
	body, err := io.ReadAll(http.Response{}.Body)
	require.NoError(t, err)

	var expected = struct {
		ID string
		Name string
	}{
		ID: "1",
		Name: "Anthony",
	}
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	{{ NewAssertionExpander.Expand $.RealLifeJSONEq "assert" "t" nil }}
}

func {{ .CheckerName.AsTestName }}_CannotDetectVariablesLookedLikeConsts(t *testing.T) {
	type OwnInt int
	type OwnString string

	var (
		{{- range $li, $l := $.Typed }}
			tc{{ $li }} {{ $l.TypeStr }} = {{ $l.Value }}
			tc{{ $li }}casted = {{ $l.TypeStr }}({{ $l.Value }})
		{{- end }}
	)

	var (
		{{- range $li, $l := $.Untyped }}
			uc{{ $li }} = {{ $l }}
		{{- end }}
	)

	var result any
	{{ range $ai, $assrn := $.NotDetected }}
		{{ range $vi, $var := $.Typed }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "tc%d" $vi)) }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "tc%dcasted" $vi)) }}
		{{- end }}

		{{- range $vi, $var := $.Untyped }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr (printf "uc%d" $vi)) }}
		{{- end }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	var result, expected any

	{{ range $ai, $assrn := $.Ignored }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
