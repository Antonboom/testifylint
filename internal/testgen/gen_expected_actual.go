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
		CheckerName          CheckerName
		Typed                []literal
		Untyped              []string
		ExpVars              []string
		Basic                test
		OtherExpActFunctions test
		Strings              test
		NotDetected          []Assertion
		RealLifeJSONEq       Assertion
		IgnoredAssertions    []Assertion
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
			"expectedObj.Val",
			"tt.expected",
			"tt.exp()",
			"tt.expPtr()",
			"expectedVal()",
			"[]int{1, 2, 3}",
			"[3]int{1, 2, 3}",
			`map[string]int{"0": 1}`,
			`user{Name: "Rob"}`,
			`struct {Name string}{Name: "Rob"}`,
			`nil`,

			"&expected",
			"&expectedObj.Val",
			"&(expectedObj.Val)",
			"&tt.expected",
			"&(tt.expected)",
			`&user{Name: "Rob"}`,

			"*expectedPtr",
			"expectedObjPtr.Val",
			"*tt.expPtr()",
			"*(tt.expPtr())",
			"ttPtr.expected",

			// NOTE(a.telyshev): Unsupported rare cases:
			// "(*expectedObjPtr).Val",
			// "(*ttPtr).expected",
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
		OtherExpActFunctions: test{
			InvalidAssertions: []Assertion{
				{
					Fn: "EqualExportedValues", Argsf: "resultObj, expectedObj",
					ReportMsgf: report, ProposedArgsf: "expectedObj, resultObj",
				},
				{
					Fn: "EqualExportedValues", Argsf: `resultObj, user{Name: "Rob"}`,
					ReportMsgf: report, ProposedArgsf: `user{Name: "Rob"}, resultObj`,
				},
				{
					Fn: "EqualExportedValues", Argsf: `resultObj, struct {Name string}{Name: "Rob"}`,
					ReportMsgf: report, ProposedArgsf: `struct {Name string}{Name: "Rob"}, resultObj`,
				},

				{Fn: "EqualValues", Argsf: "result, expected", ReportMsgf: report, ProposedArgsf: "expected, result"},
				{Fn: "EqualValues", Argsf: "result, uint32(100)", ReportMsgf: report, ProposedArgsf: "uint32(100), result"},
				{Fn: "NotEqualValues", Argsf: "result, expected", ReportMsgf: report, ProposedArgsf: "expected, result"},
				{Fn: "NotEqualValues", Argsf: "result, uint32(100)", ReportMsgf: report, ProposedArgsf: "uint32(100), result"},

				{Fn: "Exactly", Argsf: "result, expected", ReportMsgf: report, ProposedArgsf: "expected, result"},
				{Fn: "Exactly", Argsf: "result, int64(1)", ReportMsgf: report, ProposedArgsf: "int64(1), result"},

				{Fn: "InDelta", Argsf: "result, expected, 1.0", ReportMsgf: report, ProposedArgsf: "expected, result, 1.0"},
				{Fn: "InDelta", Argsf: "result, 42.42, 1.0", ReportMsgf: report, ProposedArgsf: "42.42, result, 1.0"},

				{
					Fn: "InDeltaMapValues", Argsf: "result, expected, 2.0",
					ReportMsgf: report, ProposedArgsf: "expected, result, 2.0",
				},
				{
					Fn: "InDeltaMapValues", Argsf: `result, map[string]float64{"score": 0.99}, 2.0`,
					ReportMsgf: report, ProposedArgsf: `map[string]float64{"score": 0.99}, result, 2.0`,
				},

				{
					Fn: "InDeltaSlice", Argsf: "result, expected, 1.0",
					ReportMsgf: report, ProposedArgsf: "expected, result, 1.0",
				},
				{
					Fn: "InDeltaSlice", Argsf: `result, []float64{0.98, 0.99}, 1.0`,
					ReportMsgf: report, ProposedArgsf: `[]float64{0.98, 0.99}, result, 1.0`,
				},

				{
					Fn: "InEpsilon", Argsf: "result, expected, 0.0001",
					ReportMsgf: report, ProposedArgsf: "expected, result, 0.0001",
				},
				{
					Fn: "InEpsilon", Argsf: "result, 42.42, 0.0001",
					ReportMsgf: report, ProposedArgsf: "42.42, result, 0.0001",
				},

				{
					Fn: "InEpsilonSlice", Argsf: "result, expected, 0.0001",
					ReportMsgf: report, ProposedArgsf: "expected, result, 0.0001",
				},
				{
					Fn: "InEpsilonSlice", Argsf: `result, []float64{0.9801, 0.9902}, 0.0001`,
					ReportMsgf: report, ProposedArgsf: `[]float64{0.9801, 0.9902}, result, 0.0001`,
				},

				{Fn: "IsType", Argsf: "result, expected", ReportMsgf: report, ProposedArgsf: "expected, result"},
				{Fn: "IsType", Argsf: "result, user{}", ReportMsgf: report, ProposedArgsf: "user{}, result"},
				{Fn: "IsType", Argsf: "result, (*user)(nil)", ReportMsgf: report, ProposedArgsf: "(*user)(nil), result"},

				{Fn: "Same", Argsf: "resultPtr, expectedPtr", ReportMsgf: report, ProposedArgsf: "expectedPtr, resultPtr"},
				{Fn: "Same", Argsf: "&value, expectedPtr", ReportMsgf: report, ProposedArgsf: "expectedPtr, &value"},
				{Fn: "Same", Argsf: "value, &expected", ReportMsgf: report, ProposedArgsf: "&expected, value"},
				{Fn: "NotSame", Argsf: "resultPtr, expectedPtr", ReportMsgf: report, ProposedArgsf: "expectedPtr, resultPtr"},
				{Fn: "NotSame", Argsf: "&value, expectedPtr", ReportMsgf: report, ProposedArgsf: "expectedPtr, &value"},
				{Fn: "NotSame", Argsf: "value, &expected", ReportMsgf: report, ProposedArgsf: "&expected, value"},

				{
					Fn: "WithinDuration", Argsf: "resultTime, expectedTime, time.Second",
					ReportMsgf: report, ProposedArgsf: "expectedTime, resultTime, time.Second",
				},
				{
					Fn: "WithinDuration", Argsf: "resultTime, time.Date(2023, 01, 12, 11, 46, 33, 0, nil), 100*time.Millisecond",
					ReportMsgf: report, ProposedArgsf: "time.Date(2023, 01, 12, 11, 46, 33, 0, nil), resultTime, 100*time.Millisecond",
				},
			},
			ValidAssertions: []Assertion{
				{Fn: "EqualExportedValues", Argsf: "expectedObj, resultObj"},
				{Fn: "EqualExportedValues", Argsf: `user{Name: "Rob"}, resultObj`},
				{Fn: "EqualExportedValues", Argsf: `struct {Name string}{Name: "Rob"}, resultObj`},

				{Fn: "EqualValues", Argsf: "expected, result"},
				{Fn: "EqualValues", Argsf: "uint32(100), result"},
				{Fn: "NotEqualValues", Argsf: "expected, result"},
				{Fn: "NotEqualValues", Argsf: "uint32(100), result"},

				{Fn: "Exactly", Argsf: "expected, result"},
				{Fn: "Exactly", Argsf: "int64(1), result"},

				{Fn: "InDelta", Argsf: "expected, result, 1.0"},
				{Fn: "InDelta", Argsf: "42.42, result, 1.0"},

				{Fn: "InDeltaMapValues", Argsf: "expected, result, 2.0"},
				{Fn: "InDeltaMapValues", Argsf: `map[string]float64{"score": 0.99}, result, 2.0`},

				{Fn: "InEpsilon", Argsf: "expected, result, 0.0001"},
				{Fn: "InEpsilon", Argsf: "42.42, result, 0.0001"},

				{Fn: "IsType", Argsf: "expected, result"},
				{Fn: "IsType", Argsf: "user{}, result"},
				{Fn: "IsType", Argsf: "(*user)(nil), result"},

				{Fn: "Same", Argsf: "expectedPtr, resultPtr"},
				{Fn: "Same", Argsf: "&value, resultPtr"},
				{Fn: "NotSame", Argsf: "expectedPtr, resultPtr"},
				{Fn: "NotSame", Argsf: "&value, resultPtr"},

				{Fn: "WithinDuration", Argsf: "expectedTime, resultTime, time.Second"},
				{Fn: "WithinDuration", Argsf: "time.Date(2023, 01, 12, 11, 46, 33, 0, nil), resultTime, 100*time.Millisecond"},
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
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "nil, nil"},
			{Fn: "Equal", Argsf: `"value", "value"`},
			{Fn: "Equal", Argsf: "expected, expected"},
			{Fn: "Equal", Argsf: "value, &resultPtr"},
			{Fn: "Equal", Argsf: "[]int{1, 2}, map[int]int{1: 2}"},
			{Fn: "NotEqual", Argsf: "result, result"},
			{Fn: "NotEqual", Argsf: "value, &resultPtr"},
			{Fn: "EqualExportedValues", Argsf: `user{Name: "Rob"}, struct {Name string}{Name: "Rob"}`},
			{Fn: "EqualValues", Argsf: "uint32(100), int32(100)"},
			{Fn: "Exactly", Argsf: "int32(200), int64(200)"},
			{Fn: "NotEqualValues", Argsf: "int32(100), uint32(100)"},
			{Fn: "InDelta", Argsf: "42.42, expected, 1.0"},
			{Fn: "InDeltaMapValues", Argsf: `map[string]float64{"score": 0.99}, nil, 2.0`},
			{Fn: "InDeltaSlice", Argsf: `[]float64{0.98, 0.99}, []float64{0.97, 0.99}, 1.0`},
			{Fn: "InEpsilon", Argsf: "42.42, 0.0001, 0.0001"},
			{Fn: "IsType", Argsf: "(*user)(nil), user{}"},
			{Fn: "Same", Argsf: "&value, &value"},
			{Fn: "Same", Argsf: "resultPtr, &value"},
			{Fn: "Same", Argsf: "expectedPtr, &value"},
			{Fn: "Same", Argsf: "&expected, value"},
			{Fn: "NotSame", Argsf: "&value, &value"},
			{Fn: "NotSame", Argsf: "resultPtr, &value"},
			{Fn: "NotSame", Argsf: "expectedPtr, &value"},
			{Fn: "NotSame", Argsf: "&expected, value"},
			{Fn: "WithinDuration", Argsf: "expectedTime, time.Now(), time.Second"},
			{Fn: "WithinDuration", Argsf: "time.Date(2023, 01, 12, 11, 46, 33, 0, nil), " +
				"time.Date(2023, 01, 12, 11, 46, 33, 0, nil), time.Millisecond"},
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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type user struct {
	Name string
	id   uint64
}

type testCase struct { expected string } //
func (c testCase) exp() string { return c.expected }
func (c testCase) expPtr() *string { return &c.expected }

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
	var expectedPtr *string
	var tt testCase
	var ttPtr *testCase
	expectedVal := func() any { return nil }
	var expectedObj struct { Val int }
	var expectedObjPtr = &expectedObj

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

func {{ .CheckerName.AsTestName }}_Other(t *testing.T) {
	var (
		result, expected any
		resultPtr, expectedPtr *int
		resultObj, expectedObj user
		resultTime, expectedTime time.Time
		value int
	)

	// Invalid.
	{
		{{- range $ai, $assrn := $.OtherExpActFunctions.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.OtherExpActFunctions.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
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
	var (
		result, expected any
		resultPtr, expectedPtr *int
		value int
		expectedTime time.Time
	)

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
