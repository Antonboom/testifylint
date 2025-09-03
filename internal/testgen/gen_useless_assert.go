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
		defaultReport = checker + ": meaningless assertion"
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
		"IsNotType":           "value, value",
		"IsType":              "value, value",
		"JSONEq":              "str, str",
		"Less":                "value, value",
		"LessOrEqual":         "value, value",
		"NotElementsMatch":    "value, value",
		"NotEqual":            "value, value",
		"NotEqualValues":      "value, value",
		"NotErrorAs":          "err, err",
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
			{Fn: "IsNotType", Argsf: "(*testCase)(nil), (*testCase)(nil)", ReportMsgf: sameVarReport},
			{Fn: "IsType", Argsf: "(*testCase)(nil), (*testCase)(nil)", ReportMsgf: sameVarReport},

			{Fn: "Empty", Argsf: `""`, ReportMsgf: defaultReport},
			{Fn: "Empty", Argsf: "``", ReportMsgf: defaultReport},
			{Fn: "Empty", Argsf: `"value"`, ReportMsgf: defaultReport},
			{Fn: "Error", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "False", Argsf: "false", ReportMsgf: defaultReport},
			{Fn: "False", Argsf: "true", ReportMsgf: defaultReport},
			{Fn: "Implements", Argsf: "(*any)(nil), new(testing.T)", ReportMsgf: defaultReport},
			{Fn: "Negative", Argsf: "42", ReportMsgf: defaultReport},
			{Fn: "Negative", Argsf: "0", ReportMsgf: defaultReport},
			{Fn: "Negative", Argsf: "-42", ReportMsgf: defaultReport},
			{Fn: "Nil", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "NoError", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "NotEmpty", Argsf: `""`, ReportMsgf: defaultReport},
			{Fn: "NotEmpty", Argsf: "``", ReportMsgf: defaultReport},
			{Fn: "NotEmpty", Argsf: `"value"`, ReportMsgf: defaultReport},
			{Fn: "NotImplements", Argsf: "(*any)(nil), new(testing.T)", ReportMsgf: defaultReport},
			{Fn: "NotNil", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "42", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "0", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "-42", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: `""`, ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "``", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: `"value"`, ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "false", ReportMsgf: defaultReport},
			{Fn: "NotZero", Argsf: "true", ReportMsgf: defaultReport},
			{Fn: "Positive", Argsf: "42", ReportMsgf: defaultReport},
			{Fn: "Positive", Argsf: "0", ReportMsgf: defaultReport},
			{Fn: "Positive", Argsf: "-42", ReportMsgf: defaultReport},
			{Fn: "True", Argsf: "false", ReportMsgf: defaultReport},
			{Fn: "True", Argsf: "true", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "42", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "0", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "-42", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: `""`, ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "``", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: `"value"`, ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "nil", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "false", ReportMsgf: defaultReport},
			{Fn: "Zero", Argsf: "true", ReportMsgf: defaultReport},

			{Fn: "Negative", Argsf: "len(x)", ReportMsgf: defaultReport},
			{Fn: "Less", Argsf: "len(x), 0", ReportMsgf: defaultReport},
			{Fn: "Greater", Argsf: "0, len(x)", ReportMsgf: defaultReport},
			{Fn: "GreaterOrEqual", Argsf: "len(x), 0", ReportMsgf: defaultReport},
			{Fn: "LessOrEqual", Argsf: "0, len(x)", ReportMsgf: defaultReport},

			{Fn: "Negative", Argsf: "uint(42)", ReportMsgf: defaultReport},
			{Fn: "Less", Argsf: "uint(42), 0", ReportMsgf: defaultReport},
			{Fn: "Greater", Argsf: "0, uint(42)", ReportMsgf: defaultReport},
			{Fn: "GreaterOrEqual", Argsf: "uint(42), 0", ReportMsgf: defaultReport},
			{Fn: "LessOrEqual", Argsf: "0, uint(42)", ReportMsgf: defaultReport},
		},
		InvalidAssertions: twoSideAssertions,
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "value, 42"},
			{Fn: "Equal", Argsf: `value, "value"`},
			{Fn: "Equal", Argsf: `tc.A(), "tc.A()"`},
			{Fn: "Equal", Argsf: "testCase{}.A().B().C(), tc.A().B().C()"},
			{Fn: "IsNotType", Argsf: "tc, testCase{}"},
			{Fn: "IsType", Argsf: "tc, testCase{}"},

			{Fn: "Empty", Argsf: "str"},
			{Fn: "Error", Argsf: "err"},
			{Fn: "False", Argsf: "b"},
			{Fn: "Implements", Argsf: "(*testing.TB)(nil), new(testing.T)"},
			{Fn: "Negative", Argsf: "num"},
			{Fn: "Nil", Argsf: "new(testCase)"},
			{Fn: "NoError", Argsf: "err"},
			{Fn: "NotEmpty", Argsf: "str"},
			{Fn: "NotImplements", Argsf: "(*testing.TB)(nil), new(testing.T)"},
			{Fn: "NotNil", Argsf: "new(testCase)"},
			{Fn: "NotZero", Argsf: "num"},
			{Fn: "NotZero", Argsf: "str"},
			{Fn: "NotZero", Argsf: "new(testCase)"},
			{Fn: "NotZero", Argsf: "b"},
			{Fn: "Positive", Argsf: "num"},
			{Fn: "True", Argsf: "b"},
			{Fn: "Zero", Argsf: "num"},
			{Fn: "Zero", Argsf: "str"},
			{Fn: "Zero", Argsf: "new(testCase)"},
			{Fn: "Zero", Argsf: "b"},

			// NOTE(a.telyshev): An unsigned value can be 0.
			{Fn: "Positive", Argsf: "len(x)"},
			{Fn: "Positive", Argsf: "uint(42)"},
			{Fn: "Greater", Argsf: "len(x), 0"},
			{Fn: "Greater", Argsf: "uint(42), 0"},
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
	var b bool
	var tc testCase
	var x []int

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
