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

	for _, zeroType := range []string{"", "int", "int8", "int16", "int32", "int64"} {
		v := fmt.Sprintf("%s(a)", zeroType)
		zero := fmt.Sprintf("%s(0)", zeroType)

		if zeroType == "" {
			v, zero = "a", "0"
		}

		invalidAssertions = append(invalidAssertions,
			Assertion{Fn: "Less", Argsf: "a, " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},
			Assertion{Fn: "Greater", Argsf: zero + ", a", ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: "a"},
			Assertion{Fn: "True", Argsf: v + " < " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: zero + " > " + v, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: v + " >= " + zero, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: zero + " <= " + v, ReportMsgf: report, ProposedFn: "Negative", ProposedArgsf: v},

			Assertion{Fn: "Greater", Argsf: "a, " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
			Assertion{Fn: "Less", Argsf: zero + ", a", ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: "a"},
			Assertion{Fn: "True", Argsf: v + " > " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "True", Argsf: zero + " < " + v, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: v + " <= " + zero, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
			Assertion{Fn: "False", Argsf: zero + " >= " + v, ReportMsgf: report, ProposedFn: "Positive", ProposedArgsf: v},
		)
	}

	var ignoredAssertions []Assertion

	for _, fn := range []string{"Equal", "NotEqual", "GreaterOrEqual", "LessOrEqual"} {
		for _, arg := range []string{"-1", "0", "1"} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: arg + ", a"},
				Assertion{Fn: fn, Argsf: "a, " + arg},
			)
		}
	}
	for _, fn := range []string{"Equal", "NotEqual", "GreaterOrEqual", "LessOrEqual"} {
		for _, zeroType := range []string{
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"CustomInt16",
		} {
			v := fmt.Sprintf("%s(a)", zeroType)
			zero := fmt.Sprintf("%s(0)", zeroType)

			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: zero + ", " + v},
				Assertion{Fn: fn, Argsf: v + ", " + zero},
			)
		}
	}

	for _, fn := range []string{"Greater", "Less"} {
		for _, arg := range []string{"-1", "1"} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: arg + ", a"},
				Assertion{Fn: fn, Argsf: "a, " + arg},
			)
		}
	}
	for _, fn := range []string{"Greater", "Less"} {
		for _, zeroType := range []string{
			"uint", "uint8", "uint16", "uint32", "uint64",
			"CustomInt16",
		} {
			v := fmt.Sprintf("%s(a)", zeroType)
			zero := fmt.Sprintf("%s(0)", zeroType)

			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: fn, Argsf: zero + ", " + v},
				Assertion{Fn: fn, Argsf: v + ", " + zero},
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
		v := fmt.Sprintf("%s(a)", zeroType)
		zero := fmt.Sprintf("%s(0)", zeroType)

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
		"uint", "uint8", "uint16", "uint32", "uint64",
		"CustomInt16",
	} {
		v := fmt.Sprintf("%s(a)", zeroType)
		zero := fmt.Sprintf("%s(0)", zeroType)

		for _, op := range []string{">", "<"} {
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
		v := fmt.Sprintf("%s(a)", zeroType)
		zero := fmt.Sprintf("%s(0)", zeroType)

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
	for _, zeroType := range []string{
		"uint", "uint8", "uint16", "uint32", "uint64",
		"CustomInt16",
	} {
		v := fmt.Sprintf("%s(a)", zeroType)
		zero := fmt.Sprintf("%s(0)", zeroType)

		for _, op := range []string{">=", "<="} {
			ignoredAssertions = append(ignoredAssertions,
				Assertion{Fn: "False", Argsf: fmt.Sprintf("%s %s %s", v, op, zero)},
				Assertion{Fn: "False", Argsf: fmt.Sprintf("%s %s %s", zero, op, v)},
			)
		}
	}

	// These one will be reported by useless-assert.
	ignoredAssertions = append(ignoredAssertions,
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
	)

	// These one will be reported by incorrect-assert.
	ignoredAssertions = append(ignoredAssertions,
		Assertion{Fn: "Positive", Argsf: "uint(a)"},
		Assertion{Fn: "Negative", Argsf: "uint(a)"},
		Assertion{Fn: "Greater", Argsf: "uint(a), 0"},
		Assertion{Fn: "Less", Argsf: "uint(a), 0"},
	)

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
		IgnoredAssertions: ignoredAssertions,
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
	var a int

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

type CustomInt16 int16
`
