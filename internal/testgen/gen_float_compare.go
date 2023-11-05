package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type FloatCompareTestsGenerator struct{}

func (FloatCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewFloatCompare()
}

func (g FloatCompareTestsGenerator) TemplateData() any {
	var (
		name      = g.Checker().Name()
		report    = name + ": use %s.InEpsilon (or %s)"
		reportFmt = name + ": use %s.InEpsilonf (or %s)" // Workaround for two proposed functions.
	)

	type floatDetectionTest struct {
		Vars  []string
		Assrn Assertion
	}

	return struct {
		CheckerName           CheckerName
		FloatDetection        floatDetectionTest
		InvalidAssertions     []Assertion
		InvalidAssertionsF    []Assertion
		ValidAssertions       []Assertion
		UnsupportedAssertions []Assertion
	}{
		CheckerName: CheckerName(name),
		FloatDetection: floatDetectionTest{
			Vars: []string{
				"a",
				"b",
				"cc.value",
				"d",
				"e",
				"(*f).value",
				"*g",
				"h.Calculate()",
				"floatOp()",
				"number(a) + b",
				"cc.value - (*f).value",
				"d * e / a",
				"math.Round(float64(floatOp()))",
			},
			Assrn: Assertion{Fn: "Equal", Argsf: "42.42, %s", ReportMsgf: report, ProposedFn: "InDelta"},
		},
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "%s, %s", ReportMsgf: report, ProposedFn: "InDelta"},
			{Fn: "EqualValues", Argsf: "%s, %s", ReportMsgf: report, ProposedFn: "InDelta"},
			{Fn: "Exactly", Argsf: "%s, %s", ReportMsgf: report, ProposedFn: "InDelta"},
			{Fn: "True", Argsf: "%s == %s", ReportMsgf: report, ProposedFn: "InDelta"},
			{Fn: "False", Argsf: "%s != %s", ReportMsgf: report, ProposedFn: "InDelta"},
		},
		InvalidAssertionsF: []Assertion{
			{Fn: "Equal", Argsf: "%s, %s", ReportMsgf: reportFmt, ProposedFn: "InDelta"},
			{Fn: "EqualValues", Argsf: "%s, %s", ReportMsgf: reportFmt, ProposedFn: "InDelta"},
			{Fn: "Exactly", Argsf: "%s, %s", ReportMsgf: reportFmt, ProposedFn: "InDelta"},
			{Fn: "True", Argsf: "%s == %s", ReportMsgf: reportFmt, ProposedFn: "InDelta"},
			{Fn: "False", Argsf: "%s != %s", ReportMsgf: reportFmt, ProposedFn: "InDelta"},
		},
		ValidAssertions: []Assertion{
			{Fn: "InDelta", Argsf: "%s, %s, 0.0001"},
			{Fn: "InEpsilon", Argsf: "%s, %s, 0.0002"},
		},
		// NOTE(a.telyshev): Waiting for contribution:
		UnsupportedAssertions: []Assertion{
			{Fn: "NotEqual", Argsf: "42.42, resultFl"},
			{Fn: "Greater", Argsf: "42.42, resultFl"},
			{Fn: "GreaterOrEqual", Argsf: "42.42, resultFl"},
			{Fn: "Less", Argsf: "42.42, resultFl"},
			{Fn: "LessOrEqual", Argsf: "42.42, resultFl"},

			{Fn: "True", Argsf: "42.42 != resultFl"},
			{Fn: "True", Argsf: "42.42 > resultFl"},
			{Fn: "True", Argsf: "42.42 >= resultFl"},
			{Fn: "True", Argsf: "42.42 < resultFl"},
			{Fn: "True", Argsf: "42.42 <= resultFl"},

			{Fn: "False", Argsf: "42.42 == resultFl"},
			{Fn: "False", Argsf: "42.42 <= resultFl"},
			{Fn: "False", Argsf: "42.42 < resultFl"},
			{Fn: "False", Argsf: "42.42 <= resultFl"},
			{Fn: "False", Argsf: "42.42 > resultFl"},

			// `any` cases.
			{Fn: "NotEqual", Argsf: "42.42, result"},
			{Fn: "Greater", Argsf: "42.42, result"},
			{Fn: "GreaterOrEqual", Argsf: "42.42, result"},
			{Fn: "Less", Argsf: "42.42, result"},
			{Fn: "LessOrEqual", Argsf: "42.42, result"},
			{Fn: "True", Argsf: "42.42 != result"},
			{Fn: "False", Argsf: "42.42 == result"},
		},
	}
}

func (FloatCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("FloatCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(floatCompareTestTmpl))
}

func (FloatCompareTestsGenerator) GoldenTemplate() Executor {
	// NOTE(a.telyshev): Only the developer knows the needed epsilon / delta.
	return nil
}

const floatCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var result any
	var resultFl float64

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "result" "42.42") }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "42.42" "result" ) }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "resultFl" "42.42") }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "42.42" "resultFl" ) }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "result" "resultFl" ) }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr "resultFl" "result" ) }}
		{{- end }}

		{{ range $ai, $assrn := $.InvalidAssertionsF }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "result" "42.42") }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "42.42" "result" ) }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "resultFl" "42.42") }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "42.42" "resultFl" ) }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "result" "resultFl" ) }}
			{{ NewAssertionExpander.FmtSingleMode.Expand $assrn "assert" "t" (arr "resultFl" "result" ) }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "result" "42.42") }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "42.42" "result" ) }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "resultFl" "42.42") }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "42.42" "resultFl" ) }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "result" "resultFl" ) }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr "resultFl" "result" ) }}
		{{- end }}
	}

	// Unsupported.
	{
		{{- range $ai, $assrn := $.UnsupportedAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}

func {{ .CheckerName.AsTestName }}_NoFloatNoWorries(t *testing.T) {
	var result any
	var resultInt int64

	{{ range $ai, $assrn := $.InvalidAssertions }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "result" "42") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "42" "result") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "resultInt" "42") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "42" "resultInt") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "result" "resultInt") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "resultInt" "result") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "result" "result") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "resultInt" "resultInt") }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" (arr "42" "42") }}
	{{- end }}
}

{{ range $bi, $bits := arr "32" "64" }}
func {{ $.CheckerName.AsTestName }}_Float{{ $bits }}Detection(t *testing.T) {
	type number float{{ $bits }}
	type withFloat{{ $bits }} struct{ value float{{ $bits }} }
	floatOp := func() float{{ $bits }} { return 0. }

	var a float{{ $bits }}
	var b number
	var cc withFloat{{ $bits }}
	d := float{{ $bits }}(1.01)
	const e = float{{ $bits }}(2.02)
	f := new(withFloat{{ $bits }})
	var g *float{{ $bits }}
	var h withFloat{{ $bits }}Method

	{{ range $vi, $var := $.FloatDetection.Vars }}
		{{- NewAssertionExpander.NotFmtSingleMode.Expand $.FloatDetection.Assrn "assert" "t" (arr $var) }}
	{{ end -}}
}

type withFloat{{ $bits }}Method struct{} //
func (withFloat{{ $bits }}Method) Calculate() float{{ $bits }} { return 0. }
{{ end }}
`
