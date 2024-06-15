package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type UselessAssertTestsGenerator struct{}

func (UselessAssertTestsGenerator) Checker() checkers.Checker {
	return checkers.NewUselessAssert()
}

func (g UselessAssertTestsGenerator) TemplateData() any {
	var (
		checker       = g.Checker().Name()
		sameVarReport = checker + ": asserting of the same variable"
	)

	var twoSideAssertions []Assertion
	for fn, args := range map[string]string{
		"Contains":            "value, value",
		"ElementsMatch":       "value, value",
		"Equal":               "value, value",
		"EqualExportedValues": "value, value",
		"EqualValues":         "value, value",
		"ErrorAs":             "err, err",
		"ErrorIs":             "err, err",
		"Exactly":             "value, value",
		"Greater":             "value, value",
		"GreaterOrEqual":      "value, value",
		"Implements":          "value, value",
		"InDelta":             "value, value, 0.01",
		"InDeltaMapValues":    "value, value, 0.01",
		"InDeltaSlice":        "value, value, 0.01",
		"InEpsilon":           "value, value, 0.0001",
		"InEpsilonSlice":      "value, value, 0.0001",
		"IsType":              "value, value",
		"JSONEq":              "str, str",
		"Less":                "value, value",
		"LessOrEqual":         "value, value",
		"NotEqual":            "value, value",
		"NotEqualValues":      "value, value",
		"NotErrorIs":          "err, err",
		"NotRegexp":           "value, value",
		"NotSame":             "value, value",
		"NotSubset":           "value, value",
		"Regexp":              "value, value",
		"Same":                "value, value",
		"Subset":              "value, value",
		"WithinDuration":      "elapsed, elapsed, time.Second",
		"YAMLEq":              "str, str",
	} {
		twoSideAssertions = append(twoSideAssertions,
			Assertion{Fn: fn, Argsf: args, ReportMsgf: sameVarReport})
	}

	for _, args := range []string{
		"num > num",
		"num < num",
		"num >= num",
		"num <= num",
		"num == num",
		"num != num",
	} {
		for _, fn := range []string{"True", "False"} {
			twoSideAssertions = append(twoSideAssertions,
				Assertion{Fn: fn, Argsf: args, ReportMsgf: sameVarReport})
		}
	}

	sortAssertions(twoSideAssertions)

	return struct {
		CheckerName            CheckerName
		InvalidAssertionsSmoke []Assertion
		InvalidAssertions      []Assertion
		ValidAssertions        []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertionsSmoke: []Assertion{
			{Fn: "Equal", Argsf: "42, 42", ReportMsgf: sameVarReport},
			{Fn: "Equal", Argsf: `"value", "value"`, ReportMsgf: sameVarReport},
			{Fn: "Equal", Argsf: "value, value", ReportMsgf: sameVarReport},
			{Fn: "Equal", Argsf: "tc.A(), tc.A()", ReportMsgf: sameVarReport},
			{Fn: "Equal", Argsf: "testCase{}.A().B().C(), testCase{}.A().B().C()", ReportMsgf: sameVarReport},
			{Fn: "IsType", Argsf: "(*testCase)(nil), (*testCase)(nil)", ReportMsgf: sameVarReport},
		},
		InvalidAssertions: twoSideAssertions,
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "value, 42"},
			{Fn: "Equal", Argsf: `value, "value"`},
			{Fn: "Equal", Argsf: `tc.A(), "tc.A()"`},
			{Fn: "Equal", Argsf: "testCase{}.A().B().C(), tc.A().B().C()"},
			{Fn: "IsType", Argsf: "tc, testCase{}"},
		},
	}
}

func (UselessAssertTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("UselessAssertTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(uselessAssertTestTmpl))
}

func (UselessAssertTestsGenerator) GoldenTemplate() Executor {
	// NOTE(a.telyshev): Only the developer understands the correct picture.
	return nil
}

const uselessAssertTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var value any
	var err error
	var elapsed time.Time
	var str string
	var num int
	var tc testCase

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertionsSmoke }}
			{{ NewAssertionExpander.FullMode.Expand $assrn "assert" "t" nil }}
		{{- end }}

		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.FullMode.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}

type testCase struct{}

func (testCase) A() testCase { return testCase{} }
func (testCase) B() testCase { return testCase{} }
func (testCase) C() testCase { return testCase{} }
`
