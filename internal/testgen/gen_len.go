package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type LenTestsGenerator struct{}

func (LenTestsGenerator) Checker() checkers.Checker {
	return checkers.NewLen()
}

func (g LenTestsGenerator) TemplateData() any {
	var (
		checker    = g.Checker().Name()
		report     = checker + ": use %s.%s"
		proposedFn = "Len"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "len(arr), 42", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "Equal", Argsf: "42, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "Equal", Argsf: "value, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, value"},
			{Fn: "Equal", Argsf: "len(expArr), len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, len(expArr)"},

			{Fn: "EqualValues", Argsf: "len(arr), 42", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "EqualValues", Argsf: "42, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "EqualValues", Argsf: "value, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, value"},
			{
				Fn: "EqualValues", Argsf: "len(expArr), len(arr)", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "arr, len(expArr)",
			},

			{Fn: "Exactly", Argsf: "len(arr), 42", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "Exactly", Argsf: "42, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "Exactly", Argsf: "value, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, value"},
			{Fn: "Exactly", Argsf: "len(expArr), len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, len(expArr)"},

			{Fn: "True", Argsf: "len(arr) == 42", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "True", Argsf: "42 == len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, 42"},
			{Fn: "True", Argsf: "len(arr) == value", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, value"},
			{Fn: "True", Argsf: "len(arr) == len(expArr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, len(expArr)"},

			// Constant cases.
			{Fn: "Equal", Argsf: "constNum, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, constNum"},
			{Fn: "EqualValues", Argsf: "constNum, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, constNum"},
			{Fn: "Exactly", Argsf: "constNum, len(arr)", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, constNum"},
			{Fn: "True", Argsf: "len(arr) == constNum", ReportMsgf: report, ProposedFn: proposedFn, ProposedArgsf: "arr, constNum"},

			// Type conversions cases.
			{
				Fn: "Equal", Argsf: "42, len(string(resp))", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "string(resp), 42",
			},
			{
				Fn: "Equal", Argsf: "42, len([]byte(resp))", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "[]byte(resp), 42",
			},
			{
				Fn: "Equal", Argsf: "42, len(json.RawMessage(resp))", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "json.RawMessage(resp), 42",
			},
			{
				Fn: "Equal", Argsf: "42, len(string([]byte(json.RawMessage(resp))))", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "string([]byte(json.RawMessage(resp))), 42",
			},
			{
				Fn: "True", Argsf: "len(string(resp)) == 42", ReportMsgf: report,
				ProposedFn: proposedFn, ProposedArgsf: "string(resp), 42",
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Len", Argsf: "arr, 42"},
			{Fn: "Len", Argsf: "arr, value"},
			{Fn: "Len", Argsf: "arr, len(arr)"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "len(arr), value"},
			{Fn: "EqualValues", Argsf: "len(arr), value"},
			{Fn: "Exactly", Argsf: "len(arr), value"},
			{Fn: "True", Argsf: "value == len(arr)"},

			{Fn: "NotEqual", Argsf: "42, len(arr)"},
			{Fn: "NotEqual", Argsf: "len(arr), 42"},
			{Fn: "NotEqualValues", Argsf: "42, len(arr)"},
			{Fn: "NotEqualValues", Argsf: "len(arr), 42"},
			{Fn: "Greater", Argsf: "len(arr), 42"},
			{Fn: "Greater", Argsf: "42, len(arr)"},
			{Fn: "GreaterOrEqual", Argsf: "len(arr), 42"},
			{Fn: "GreaterOrEqual", Argsf: "42, len(arr)"},
			{Fn: "Less", Argsf: "len(arr), 42"},
			{Fn: "Less", Argsf: "42, len(arr)"},
			{Fn: "LessOrEqual", Argsf: "len(arr), 42"},
			{Fn: "LessOrEqual", Argsf: "42, len(arr)"},

			{Fn: "True", Argsf: "42 != len(arr)"},
			{Fn: "True", Argsf: "len(arr) != 42"},
			{Fn: "True", Argsf: "42 > len(arr)"},
			{Fn: "True", Argsf: "len(arr) > 42"},
			{Fn: "True", Argsf: "42 >= len(arr)"},
			{Fn: "True", Argsf: "len(arr) >= 42"},
			{Fn: "True", Argsf: "42 < len(arr)"},
			{Fn: "True", Argsf: "len(arr) < 42"},
			{Fn: "True", Argsf: "42 <= len(arr)"},
			{Fn: "True", Argsf: "len(arr) >= 42"},

			{Fn: "False", Argsf: "42 == len(arr)"},
			{Fn: "False", Argsf: "len(arr) == 42"},
			{Fn: "False", Argsf: "42 != len(arr)"},
			{Fn: "False", Argsf: "len(arr) != 42"},
			{Fn: "False", Argsf: "42 > len(arr)"},
			{Fn: "False", Argsf: "len(arr) > 42"},
			{Fn: "False", Argsf: "42 >= len(arr)"},
			{Fn: "False", Argsf: "len(arr) >= 42"},
			{Fn: "False", Argsf: "42 < len(arr)"},
			{Fn: "False", Argsf: "len(arr) < 42"},
			{Fn: "False", Argsf: "42 <= len(arr)"},
			{Fn: "False", Argsf: "len(arr) <= 42"},
		},
	}
}

func (LenTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("LenTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(lenTestTmpl))
}

func (LenTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("LenTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(lenTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const lenTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const constNum = 10

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var arr, expArr [3]int
	var value int
	var resp string

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
