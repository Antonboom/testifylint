package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ZeroTestsGenerator struct{}

func (ZeroTestsGenerator) Checker() checkers.Checker {
	return checkers.NewZero()
}

func (g ZeroTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName       CheckerName
		Vars              []string
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		Vars:        []string{"zeroTime", "time.Time{}"},
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "%s, ts", ReportMsgf: report, ProposedFn: "Zero", ProposedArgsf: "ts%.s"},
			{Fn: "EqualValues", Argsf: "%s, ts", ReportMsgf: report, ProposedFn: "Zero", ProposedArgsf: "ts%.s"},
			{Fn: "Exactly", Argsf: "%s, ts", ReportMsgf: report, ProposedFn: "Zero", ProposedArgsf: "ts%.s"},
			{Fn: "True", Argsf: "%s.IsZero()", ReportMsgf: report, ProposedFn: "Zero", ProposedArgsf: "%s"},
			{Fn: "True", Argsf: "ts.Equal(%s)", ReportMsgf: report, ProposedFn: "Zero", ProposedArgsf: "ts%.s"},

			{Fn: "NotEqual", Argsf: "%s, ts", ReportMsgf: report, ProposedFn: "NotZero", ProposedArgsf: "ts%.s"},
			{Fn: "NotEqualValues", Argsf: "%s, ts", ReportMsgf: report, ProposedFn: "NotZero", ProposedArgsf: "ts%.s"},
			{Fn: "False", Argsf: "%s.IsZero()", ReportMsgf: report, ProposedFn: "NotZero", ProposedArgsf: "%s"},
			{Fn: "False", Argsf: "ts.Equal(%s)", ReportMsgf: report, ProposedFn: "NotZero", ProposedArgsf: "ts%.s"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Zero", Argsf: "ts"},
			{Fn: "NotZero", Argsf: "ts"},

			// compares cases.
			{Fn: "Equal", Argsf: "0, ts.Compare(zeroTime)"},
			{Fn: "EqualValues", Argsf: "0, ts.Compare(zeroTime)"},
			{Fn: "Exactly", Argsf: "0, ts.Compare(zeroTime)"},
			{Fn: "NotEqual", Argsf: "0, ts.Compare(time.Time{})"},
			{Fn: "NotEqualValues", Argsf: "0, ts.Compare(time.Time{})"},
		},
	}
}

func (ZeroTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ZeroTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(zeroTestTmpl))
}

func (ZeroTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ZeroTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(zeroTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const zeroTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var ts, zeroTime time.Time

	// Invalid.
	{
        {{- range $ai, $assrn := $.InvalidAssertions }}
            {{- range $vi, $var := $.Vars }}
                {{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
            {{- end }}
        {{- end }}
	}

	// Valid.
	{
		{{- range $ai, $assrn := $.ValidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}
}
`
