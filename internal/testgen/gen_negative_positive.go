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

	for _, zeroType := range []string{
		"", "int", "int8", "int16", "int32", "int64",
	} {
		v := zeroType + "(a)"
		zero := zeroType + "(0)"

		if zeroType == "" {
			v, zero = "a", "0"
		}

		invalidAssertions = append(invalidAssertions,
			Assertion{Fn: "Less", Argsf: v + ", " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "Greater", Argsf: zero + ", " + v, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: v + " < " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: zero + " > " + v, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: v + " >= " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: zero + " <= " + v, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
		)
	}

	for _, zeroType := range []string{
		"", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
	} {
		v := zeroType + "(a)"
		zero := zeroType + "(0)"

		if zeroType == "" {
			v, zero = "a", "0"
		}

		invalidAssertions = append(invalidAssertions,
			Assertion{Fn: "Greater", Argsf: v + ", " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "Less", Argsf: zero + ", " + v, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: v + " > " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: zero + " < " + v, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: v + " <= " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: zero + " >= " + v, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
		)
	}

	invalidAssertions = append(invalidAssertions,
		Assertion{Fn: "True", Argsf: "len(x) > 0", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "len(x)"},
		Assertion{Fn: "Greater", Argsf: "len(x), 0", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "len(x)"},
		Assertion{Fn: "Less", Argsf: "f(a), 0", ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "f(a)"},
		Assertion{Fn: "False", Argsf: "0 >= f(a)", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "f(a)"},
	)

	var ignoredAssertions []Assertion

	for _, fn := range []string{"Equal", "NotEqual", "GreaterOrEqual", "LessOrEqual"} {
		for _, arg := range []string{"-1", "0", "1"} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: arg + ", a"},
				Assertion{Fn: fn, Argsf: "a, " + arg},
			)
		}

		for _, zeroType := range []string{
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"CustomInt16",
		} {
			v := zeroType + "(a)"
			zero := zeroType + "(0)"

			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: zero + ", " + v},
				Assertion{Fn: fn, Argsf: v + ", " + zero},
			)
		}
	}

	for _, fn := range []string{"Greater", "Less"} {
		for _, arg := range []string{"-1", "CustomInt16(0)", "1"} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: arg + ", a"},
				Assertion{Fn: fn, Argsf: "a, " + arg},
			)
		}
	}

	for _, fn := range []string{"True", "False"} {
		for _, arg := range []string{"-1", "1"} {
			for _, op := range []string{">", "<", ">=", "<=", "==", "!="} {
				ignoredAssertions = append(ignoredAssertions,
					Assertion{Fn: fn, Argsf: fmt.Sprintf("a %s %s", op, arg)},
					Assertion{Fn: fn, Argsf: fmt.Sprintf("%s %s a", arg, op)},
				)
			}
		}
	}
	for _, zeroType := range []string{
		"",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"CustomInt16",
	} {
		v := zeroType + "(a)"
		zero := zeroType + "(0)"

		if zeroType == "" {
			v, zero = "a", "0"
		}

		for _, op := range []string{">=", "<=", "==", "!="} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: "True", Argsf: fmt.Sprintf("%s %s %s", v, op, zero)},
				Assertion{Fn: "True", Argsf: fmt.Sprintf("%s %s %s", zero, op, v)},
			)
		}
	}
	for _, zeroType := range []string{
		"",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"CustomInt16",
	} {
		v := zeroType + "(a)"
		zero := zeroType + "(0)"

		if zeroType == "" {
			v, zero = "a", "0"
		}

		for _, op := range []string{">", "<", "==", "!="} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: "False", Argsf: fmt.Sprintf("%s %s %s", v, op, zero)},
				Assertion{Fn: "False", Argsf: fmt.Sprintf("%s %s %s", zero, op, v)},
			)
		}
	}

	// assert.Negative only cases.
	for _, zeroType := range []string{
		"uint", "uint8", "uint16", "uint32", "uint64",
		"CustomInt16",
	} {
		v := zeroType + "(a)"
		zero := zeroType + "(0)"

		ignoredAssertions = append(ignoredAssertions,
			Assertion{Fn: "Less", Argsf: v + ", " + zero},
			Assertion{Fn: "Greater", Argsf: zero + ", " + v},
			Assertion{Fn: "True", Argsf: v + " < " + zero},
			Assertion{Fn: "True", Argsf: zero + " > " + v},
			Assertion{Fn: "False", Argsf: v + " >= " + zero},
			Assertion{Fn: "False", Argsf: zero + " <= " + v},
		)
	}

	// These one will be reported by useless-assert.
	ignoredAssertions = append(ignoredAssertions, //nolint:gocritic // Worse performance in favor of readability.
		Assertion{Fn: "Equal", Argsf: "0, 0"},
		Assertion{Fn: "Equal", Argsf: "a, a"},
		Assertion{Fn: "NotEqual", Argsf: "0, 0"},
		Assertion{Fn: "NotEqual", Argsf: "a, a"},
		Assertion{Fn: "Greater", Argsf: "0, 0"},
		Assertion{Fn: "Greater", Argsf: "a, a"},
		Assertion{Fn: "GreaterOrEqual", Argsf: "0, 0"},
		Assertion{Fn: "GreaterOrEqual", Argsf: "a, a"},
		Assertion{Fn: "Less", Argsf: "0, 0"},
		Assertion{Fn: "Less", Argsf: "a, a"},
		Assertion{Fn: "LessOrEqual", Argsf: "0, 0"},
		Assertion{Fn: "LessOrEqual", Argsf: "a, a"},
		Assertion{Fn: "True", Argsf: "a > a"},
		Assertion{Fn: "True", Argsf: "a < a"},
		Assertion{Fn: "True", Argsf: "a >= a"},
		Assertion{Fn: "True", Argsf: "a <= a"},
		Assertion{Fn: "True", Argsf: "a == a"},
		Assertion{Fn: "True", Argsf: "a != a"},
		Assertion{Fn: "False", Argsf: "-1 > -1"},
		Assertion{Fn: "False", Argsf: "-1 < -1"},
		Assertion{Fn: "False", Argsf: "-1 >= -1"},
		Assertion{Fn: "False", Argsf: "-1 <= -1"},
		Assertion{Fn: "False", Argsf: "-1 == -1"},
		Assertion{Fn: "False", Argsf: "-1 != -1"},
		Assertion{Fn: "GreaterOrEqual", Argsf: "len(x), 0"},
		Assertion{Fn: "LessOrEqual", Argsf: "0, len(x)"},
	)

	// These one will be reported by incorrect-assert.
	ignoredAssertions = append(ignoredAssertions,
		Assertion{Fn: "Negative", Argsf: "uint(a)"},
		Assertion{Fn: "Less", Argsf: "uint(a), 0"},
		Assertion{Fn: "Greater", Argsf: "0, uint(a)"},
		Assertion{Fn: "True", Argsf: "uint(a) < 0"},
		Assertion{Fn: "True", Argsf: "0 > uint(a)"},
		Assertion{Fn: "False", Argsf: "uint(a) >= 0"},
		Assertion{Fn: "False", Argsf: "0 <= uint(a)"},

		Assertion{Fn: "Negative", Argsf: "len(x)"},
		Assertion{Fn: "Less", Argsf: "len(x), 0"},
		Assertion{Fn: "Greater", Argsf: "0, len(x)"},
		Assertion{Fn: "True", Argsf: "len(x) < 0"},
		Assertion{Fn: "True", Argsf: "0 > len(x)"},
		Assertion{Fn: "False", Argsf: "len(x) >= 0"},
		Assertion{Fn: "False", Argsf: "0 <= len(x)"},
	)

	return struct {
		CheckerName          CheckerName
		InvalidAssertions    []Assertion
		ValidAssertions      []Assertion
		IgnoredAssertions    []Assertion
		RealLifeUintExamples []Assertion
	}{
		CheckerName:       CheckerName(checker),
		InvalidAssertions: invalidAssertions,
		ValidAssertions: []Assertion{
			{Fn: "Negative", Argsf: "a"},
			{Fn: "Negative", Argsf: "f(a)"},
			{Fn: "Positive", Argsf: "a"},
			{Fn: "Positive", Argsf: "len(x)"},
			{Fn: "Positive", Argsf: "f(a)"},
		},
		IgnoredAssertions: ignoredAssertions,
		RealLifeUintExamples: []Assertion{
			{
				Fn: "Less", Argsf: "uint64(0), e.VideoMinutes",
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "e.VideoMinutes",
			},
			{
				Fn: "Less", Argsf: "uint32(0), c.stats.Rto",
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "c.stats.Rto",
			},
			{
				Fn: "Less", Argsf: "uint32(0), c.stats.Ato",
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "c.stats.Ato",
			},
			{
				Fn: "Less", Argsf: "uint64(0), baseLineHeap",
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "baseLineHeap",
			},
			{
				Fn: "Greater", Argsf: "uint64(state.LastUpdatedEpoch), uint64(0)",
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "uint64(state.LastUpdatedEpoch)",
			},
			{
				Fn: "True", Argsf: `uint64(0) < prod["last_claim_time"].(uint64)`,
				ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: `prod["last_claim_time"].(uint64)`,
			},

			{Fn: "Greater", Argsf: "uint64(result.GasUsed), minGasExpected"},
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
		a int
		x []int
		f func(d int) int
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

func {{ .CheckerName.AsTestName }}_RealLifeUintExamples(t *testing.T) {
	var e struct{ VideoMinutes uint64 }
	var c struct{ stats struct{ Rto, Ato uint64 } }
	var baseLineHeap, minGasExpected uint64
	var result struct{ GasUsed int }
	var state struct{ LastUpdatedEpoch ChainEpoch }
	var prod map[string]any

	{{ range $ai, $assrn := $.RealLifeUintExamples }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}

type CustomInt16 int16
type ChainEpoch int64
`
