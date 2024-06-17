package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type NegativePositiveTestsGenerator struct{}

func (NegativePositiveTestsGenerator) Checker() checkers.Checker {
	return checkers.NewNegativePositive()
}

func (g NegativePositiveTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	var invalidAssertions []Assertion

	for _, zero := range []string{
		"0",
		"int(0)", "int8(0)", "int16(0)", "int32(0)", "int64(0)",
		"uint(0)", "uint8(0)", "uint16(0)", "uint32(0)", "uint64(0)",
	} {
		invalidAssertions = append(invalidAssertions,
			Assertion{Fn: "Less", Argsf: fmt.Sprintf("a, %s", zero), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},
			Assertion{Fn: "Greater", Argsf: fmt.Sprintf("%s, a", zero), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},

			Assertion{Fn: "Greater", Argsf: fmt.Sprintf("a, %s", zero), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
			Assertion{Fn: "Less", Argsf: fmt.Sprintf("%s, a", zero), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
		)
	}

	for zero, a := range map[string]string{
		"0":         "a",
		"int(0)":    "a",
		"int8(0)":   "int8A",
		"int16(0)":  "int16A",
		"int32(0)":  "int32A",
		"int64(0)":  "int64A",
		"uint(0)":   "uintA",
		"uint8(0)":  "uint8A",
		"uint16(0)": "uint16A",
		"uint32(0)": "uint32A",
		"uint64(0)": "uint64A",
	} {
		invalidAssertions = append(invalidAssertions,
			Assertion{Fn: "True", Argsf: fmt.Sprintf("%s < %s", a, zero), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: a},
			Assertion{Fn: "True", Argsf: fmt.Sprintf("%s > %s", zero, a), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: a},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("%s >= %s", a, zero), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: a},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("%s <= %s", zero, a), ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: a},

			Assertion{Fn: "True", Argsf: fmt.Sprintf("%s > %s", a, zero), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: a},
			Assertion{Fn: "True", Argsf: fmt.Sprintf("%s < %s", zero, a), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: a},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("%s <= %s", a, zero), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: a},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("%s >= %s", zero, a), ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: a},
		)
	}

	sortAssertions(invalidAssertions)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName:       CheckerName(checker),
		InvalidAssertions: invalidAssertions,
		ValidAssertions: []Assertion{
			{Fn: "Negative", Argsf: "a"},
			{Fn: "Positive", Argsf: "a"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "-1, a"},
			{Fn: "Equal", Argsf: "a, -1"},
			{Fn: "Equal", Argsf: "0, a"},
			{Fn: "Equal", Argsf: "a, 0"},
			{Fn: "Equal", Argsf: "1, a"},
			{Fn: "Equal", Argsf: "a, 1"},

			{Fn: "NotEqual", Argsf: "-1, a"},
			{Fn: "NotEqual", Argsf: "a, -1"},
			{Fn: "NotEqual", Argsf: "0, a"},
			{Fn: "NotEqual", Argsf: "a, 0"},
			{Fn: "NotEqual", Argsf: "1, a"},
			{Fn: "NotEqual", Argsf: "a, 1"},

			{Fn: "Greater", Argsf: "-1, a"},
			{Fn: "Greater", Argsf: "a, -1"},
			{Fn: "Greater", Argsf: "a, 1"},
			{Fn: "Greater", Argsf: "1, a"},

			{Fn: "GreaterOrEqual", Argsf: "-1, a"},
			{Fn: "GreaterOrEqual", Argsf: "a, -1"},
			{Fn: "GreaterOrEqual", Argsf: "0, a"},
			{Fn: "GreaterOrEqual", Argsf: "a, 0"},
			{Fn: "GreaterOrEqual", Argsf: "1, a"},
			{Fn: "GreaterOrEqual", Argsf: "a, 1"},

			{Fn: "Less", Argsf: "-1, a"},
			{Fn: "Less", Argsf: "a, -1"},
			{Fn: "Less", Argsf: "1, a"},
			{Fn: "Less", Argsf: "a, 1"},

			{Fn: "LessOrEqual", Argsf: "-1, a"},
			{Fn: "LessOrEqual", Argsf: "a, -1"},
			{Fn: "LessOrEqual", Argsf: "0, a"},
			{Fn: "LessOrEqual", Argsf: "a, 0"},
			{Fn: "LessOrEqual", Argsf: "1, a"},
			{Fn: "LessOrEqual", Argsf: "a, 1"},

			{Fn: "True", Argsf: "a > -1"},
			{Fn: "True", Argsf: "a < -1"},
			{Fn: "True", Argsf: "a >= -1"},
			{Fn: "True", Argsf: "a <= -1"},
			{Fn: "True", Argsf: "a == -1"},
			{Fn: "True", Argsf: "a != -1"},
			{Fn: "True", Argsf: "-1 > a"},
			{Fn: "True", Argsf: "-1 < a"},
			{Fn: "True", Argsf: "-1 >= a"},
			{Fn: "True", Argsf: "-1 <= a"},
			{Fn: "True", Argsf: "-1 == a"},
			{Fn: "True", Argsf: "-1 != a"},

			{Fn: "True", Argsf: "a >= 0"},
			{Fn: "True", Argsf: "a <= 0"},
			{Fn: "True", Argsf: "a == 0"},
			{Fn: "True", Argsf: "a != 0"},
			{Fn: "True", Argsf: "0 >= a"},
			{Fn: "True", Argsf: "0 <= a"},
			{Fn: "True", Argsf: "0 == a"},
			{Fn: "True", Argsf: "0 != a"},

			{Fn: "True", Argsf: "a > 1"},
			{Fn: "True", Argsf: "a < 1"},
			{Fn: "True", Argsf: "a >= 1"},
			{Fn: "True", Argsf: "a <= 1"},
			{Fn: "True", Argsf: "a == 1"},
			{Fn: "True", Argsf: "a != 1"},
			{Fn: "True", Argsf: "1 > a"},
			{Fn: "True", Argsf: "1 < a"},
			{Fn: "True", Argsf: "1 >= a"},
			{Fn: "True", Argsf: "1 <= a"},
			{Fn: "True", Argsf: "1 == a"},
			{Fn: "True", Argsf: "1 != a"},

			{Fn: "False", Argsf: "a > -1"},
			{Fn: "False", Argsf: "a < -1"},
			{Fn: "False", Argsf: "a >= -1"},
			{Fn: "False", Argsf: "a <= -1"},
			{Fn: "False", Argsf: "a == -1"},
			{Fn: "False", Argsf: "a != -1"},
			{Fn: "False", Argsf: "-1 > a"},
			{Fn: "False", Argsf: "-1 < a"},
			{Fn: "False", Argsf: "-1 >= a"},
			{Fn: "False", Argsf: "-1 <= a"},
			{Fn: "False", Argsf: "-1 == a"},
			{Fn: "False", Argsf: "-1 != a"},

			{Fn: "False", Argsf: "a > 0"},
			{Fn: "False", Argsf: "a < 0"},
			{Fn: "False", Argsf: "a == 0"},
			{Fn: "False", Argsf: "a != 0"},
			{Fn: "False", Argsf: "0 > a"},
			{Fn: "False", Argsf: "0 < a"},
			{Fn: "False", Argsf: "0 == a"},
			{Fn: "False", Argsf: "0 != a"},

			{Fn: "False", Argsf: "a > 1"},
			{Fn: "False", Argsf: "a < 1"},
			{Fn: "False", Argsf: "a >= 1"},
			{Fn: "False", Argsf: "a <= 1"},
			{Fn: "False", Argsf: "a == 1"},
			{Fn: "False", Argsf: "a != 1"},
			{Fn: "False", Argsf: "1 > a"},
			{Fn: "False", Argsf: "1 < a"},
			{Fn: "False", Argsf: "1 >= a"},
			{Fn: "False", Argsf: "1 <= a"},
			{Fn: "False", Argsf: "1 == a"},
			{Fn: "False", Argsf: "1 != a"},

			// These one will be reported by useless-assert.
			{Fn: "Equal", Argsf: "0, 0"},
			{Fn: "Equal", Argsf: "a, a"},
			{Fn: "NotEqual", Argsf: "0, 0"},
			{Fn: "NotEqual", Argsf: "a, a"},
			{Fn: "Greater", Argsf: "0, 0"},
			{Fn: "Greater", Argsf: "a, a"},
			{Fn: "GreaterOrEqual", Argsf: "0, 0"},
			{Fn: "GreaterOrEqual", Argsf: "a, a"},
			{Fn: "Less", Argsf: "0, 0"},
			{Fn: "Less", Argsf: "a, a"},
			{Fn: "LessOrEqual", Argsf: "0, 0"},
			{Fn: "LessOrEqual", Argsf: "a, a"},
			{Fn: "True", Argsf: "a > a"},
			{Fn: "True", Argsf: "a < a"},
			{Fn: "True", Argsf: "a >= a"},
			{Fn: "True", Argsf: "a <= a"},
			{Fn: "True", Argsf: "a == a"},
			{Fn: "True", Argsf: "a != a"},
			{Fn: "False", Argsf: "-1 > -1"},
			{Fn: "False", Argsf: "-1 < -1"},
			{Fn: "False", Argsf: "-1 >= -1"},
			{Fn: "False", Argsf: "-1 <= -1"},
			{Fn: "False", Argsf: "-1 == -1"},
			{Fn: "False", Argsf: "-1 != -1"},
		},
	}
}

func (NegativePositiveTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("NegativePositiveTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(negativePositiveTestTmpl))
}

func (NegativePositiveTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("NegativePositiveTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(negativePositiveTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const negativePositiveTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		a       int
		int8A   int8
		int16A  int16
		int32A  int32
		int64A  int64
		uintA   uint
		uint8A  uint8
		uint16A uint16
		uint32A uint32
		uint64A uint64
	)

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Ignored.
	{
		{{- range $ai, $assrn := $.IgnoredAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}
`
