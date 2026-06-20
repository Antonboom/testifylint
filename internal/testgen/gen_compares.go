package main

import (
	"fmt"
	"go/token"
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ComparesTestsGenerator struct{}

func (ComparesTestsGenerator) Checker() checkers.Checker {
	return checkers.NewCompares()
}

func (g ComparesTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()

		report           = checker + ": use %s.%s"
		reportUseT1Equal = checker + ": use t1.Equal%.s%.s"
	)

	boolOps := []token.Token{token.LAND, token.LOR}
	ignored := make([]Assertion, 0, len(boolOps)*2)
	for _, tok := range boolOps {
		ignored = append(ignored,
			Assertion{Fn: "True", Argsf: fmt.Sprintf("c %s d", tok)},
			Assertion{Fn: "False", Argsf: fmt.Sprintf("d %s c", tok)},
		)
	}

	for _, n := range []int{-3, -2, 2, 3} {
		ignored = append(ignored,
			Assertion{Fn: "Equal", Argsf: fmt.Sprintf("%d, t1.Compare(t2)", n)},
			Assertion{Fn: "Equal", Argsf: fmt.Sprintf("t1.Compare(t2), %d", n)},
		)
	}
	ignored = append(ignored,
		Assertion{Fn: "Equal", Argsf: "wantCmp, t1.Compare(t2)"},
		Assertion{Fn: "Equal", Argsf: "t1.Compare(t2), wantCmp"},
		Assertion{Fn: "Equal", Argsf: "t1.Compare(t2), t2.Compare(t1)"},
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "True", Argsf: "a == b", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a != b", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a > b", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a >= b", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a < b", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "a, b"},
			{Fn: "True", Argsf: "a <= b", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "a, b"},

			{Fn: "False", Argsf: "a == b", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a != b", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a > b", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a >= b", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a < b", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "a, b"},
			{Fn: "False", Argsf: "a <= b", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "a, b"},

			{Fn: "True", Argsf: "ptrA == ptrB", ReportMsgf: report, ProposedFn: "Same", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "True", Argsf: "ptrA != ptrB", ReportMsgf: report, ProposedFn: "NotSame", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "False", Argsf: "ptrA == ptrB", ReportMsgf: report, ProposedFn: "NotSame", ProposedArgsf: "ptrA, ptrB"},
			{Fn: "False", Argsf: "ptrA != ptrB", ReportMsgf: report, ProposedFn: "Same", ProposedArgsf: "ptrA, ptrB"},

			// `time.Time` cases.

			{Fn: "True", Argsf: "t1 == t2", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "t1, t2"},
			{Fn: "True", Argsf: "t1 != t2", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "t1, t2"},
			{Fn: "True", Argsf: "t1.After(t2)", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "t1, t2"},
			{Fn: "True", Argsf: "t1.Before(t2)", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "t1, t2"},

			{Fn: "False", Argsf: "t1 == t2", ReportMsgf: report, ProposedFn: "NotEqual", ProposedArgsf: "t1, t2"},
			{Fn: "False", Argsf: "t1 != t2", ReportMsgf: report, ProposedFn: "Equal", ProposedArgsf: "t1, t2"},
			{Fn: "False", Argsf: "t1.After(t2)", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "False", Argsf: "t1.Before(t2)", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "t1, t2"},

			// Be careful, not assert.Equal(t, t1, t2) or assert.NotEqual(t, t1, t2)!
			// Because we cannot suggest invalid flaky assertions (see `time-compare`).
			{Fn: "Equal", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "WithinDuration", ProposedArgsf: "t2, t1, 0"},
			{Fn: "EqualValues", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "WithinDuration", ProposedArgsf: "t2, t1, 0"},
			{Fn: "Exactly", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "WithinDuration", ProposedArgsf: "t2, t1, 0"},
			{Fn: "NotEqual", Argsf: "0, t1.Compare(t2)", ReportMsgf: reportUseT1Equal, ProposedFn: "False", ProposedArgsf: "t1.Equal(t2)"},
			{Fn: "NotEqualValues", Argsf: "0, t1.Compare(t2)", ReportMsgf: reportUseT1Equal, ProposedFn: "False", ProposedArgsf: "t1.Equal(t2)"},

			{Fn: "Greater", Argsf: "t1.Compare(t2), 0", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "t1, t2"},
			{Fn: "Less", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "t1, t2"},
			{Fn: "GreaterOrEqual", Argsf: "t1.Compare(t2), 0", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "LessOrEqual", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "Less", Argsf: "t1.Compare(t2), 0", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "t1, t2"},
			{Fn: "Greater", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "t1, t2"},
			{Fn: "LessOrEqual", Argsf: "t1.Compare(t2), 0", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "GreaterOrEqual", Argsf: "0, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "t1, t2"},

			{Fn: "Equal", Argsf: "1, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "t1, t2"},
			{Fn: "Equal", Argsf: "t1.Compare(t2), 1", ReportMsgf: report, ProposedFn: "Greater", ProposedArgsf: "t1, t2"},
			{Fn: "NotEqual", Argsf: "-1, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "NotEqual", Argsf: "t1.Compare(t2), -1", ReportMsgf: report, ProposedFn: "GreaterOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "Equal", Argsf: "-1, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "t1, t2"},
			{Fn: "Equal", Argsf: "t1.Compare(t2), -1", ReportMsgf: report, ProposedFn: "Less", ProposedArgsf: "t1, t2"},
			{Fn: "NotEqual", Argsf: "1, t1.Compare(t2)", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "t1, t2"},
			{Fn: "NotEqual", Argsf: "t1.Compare(t2), 1", ReportMsgf: report, ProposedFn: "LessOrEqual", ProposedArgsf: "t1, t2"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "a, b"},
			{Fn: "NotEqual", Argsf: "a, b"},
			{Fn: "Greater", Argsf: "a, b"},
			{Fn: "GreaterOrEqual", Argsf: "a, b"},
			{Fn: "Less", Argsf: "a, b"},
			{Fn: "LessOrEqual", Argsf: "a, b"},

			{Fn: "Same", Argsf: "ptrA, ptrB"},
			{Fn: "NotSame", Argsf: "ptrA, ptrB"},

			{Fn: "True", Argsf: "t1.Equal(t2)"},
			{Fn: "False", Argsf: "t1.Equal(t2)"},
			{Fn: "Greater", Argsf: "t1, t2"},
			{Fn: "Less", Argsf: "t1, t2"},
			{Fn: "LessOrEqual", Argsf: "t1, t2"},

			// time-compare cases.
			{Fn: "Equal", Argsf: "t1, t2"},
			{Fn: "EqualValues", Argsf: "t1, t2"},
			{Fn: "Exactly", Argsf: "t1, t2"},
			{Fn: "NotEqual", Argsf: "t1, t2"},
			{Fn: "NotEqualValues", Argsf: "t1, t2"},
		},
		IgnoredAssertions: ignored,
	}
}

func (ComparesTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ComparesTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(comparesTestTmpl))
}

func (ComparesTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ComparesTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(comparesTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const comparesTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var a, b int
	var c, d bool
	var wantCmp int
	var ptrA, ptrB *int
	var t1, t2 time.Time

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
