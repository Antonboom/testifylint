package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type EmptyTestsGenerator struct{}

func (EmptyTestsGenerator) Checker() checkers.Checker {
	return checkers.NewEmpty()
}

func (g EmptyTestsGenerator) TemplateData() any {
	var (
		checker          = g.Checker().Name()
		reportUse        = checker + ": use %s.%s"
		reportRemoveLen  = checker + ": remove unnecessary len"
		reportRemoveConv = checker + ": remove unnecessary string conversion"
	)

	vars := []string{"elems", "str", "string(str)", "b", "string(b)", `[]string{"e"}`}

	type test struct {
		Name              string
		Vars              []string
		InvalidAssertions []Assertion
		ExtraStringConv   Assertion
		ValidAssertions   []Assertion
	}

	type lenTest struct {
		Vars  []string
		Assrn Assertion
	}

	return struct {
		CheckerName       CheckerName
		LenTest           lenTest
		Tests             []test
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		LenTest: lenTest{
			Vars:  []string{"elems", "arr", "arrPtr", "sl", "mp", "str", "b", "ch", "[]byte(str)", "string(str)", `[]string{"e"}`},
			Assrn: Assertion{Fn: "Equal", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
		},
		Tests: []test{
			{
				Name: "assert.Empty cases",
				Vars: vars,
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n == 0, n <= 0, n < 1
					// 0 == n, 0 >= n, 1 > n
					{Fn: "Len", Argsf: "%s, 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Zero", Argsf: "%s", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Zero", Argsf: "len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "EqualValues", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "EqualValues", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Exactly", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Exactly", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "LessOrEqual", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "GreaterOrEqual", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Less", Argsf: "len(%s), 1", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Greater", Argsf: "1, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},

					// Empty string cases.
					{Fn: "Equal", Argsf: `"", %s`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "EqualValues", Argsf: `"", %s`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Exactly", Argsf: `"", %s`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Equal", Argsf: "``, %s", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "EqualValues", Argsf: "``, %s", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
					{Fn: "Exactly", Argsf: "``, %s", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},

					// Simplification cases.
					{Fn: "Empty", Argsf: "len(%s)", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "%s"},
				},
				ExtraStringConv: Assertion{
					Fn: "Empty", Argsf: "string(str)", ReportMsgf: reportRemoveConv, ProposedFn: "Empty", ProposedArgsf: "str",
				},
				ValidAssertions: []Assertion{
					{Fn: "Empty", Argsf: "%s"},
				},
			},
			{
				Name: "assert.NotEmpty cases",
				Vars: vars,
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n != 0, n > 0
					// 0 != n, 0 < n
					{Fn: "Positive", Argsf: "len(%s)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotZero", Argsf: "%s", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotZero", Argsf: "len(%s)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqual", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqualValues", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqualValues", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "Greater", Argsf: "len(%s), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "Less", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					// Empty string cases.
					{Fn: "NotEqual", Argsf: `"", %s`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqualValues", Argsf: `"", %s`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					{Fn: "NotEqual", Argsf: "``, %s", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
					{Fn: "NotEqualValues", Argsf: "``, %s", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},

					// Simplification cases.
					{Fn: "NotEmpty", Argsf: "len(%s)", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "%s"},
				},
				ExtraStringConv: Assertion{
					Fn: "NotEmpty", Argsf: "string(str)", ReportMsgf: reportRemoveConv, ProposedFn: "Empty", ProposedArgsf: "str",
				},
				ValidAssertions: []Assertion{
					{Fn: "NotEmpty", Argsf: "%s"},
				},
			},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Len", Argsf: "elems, len(elems)"},
			{Fn: "Len", Argsf: "elems, 1"},

			{Fn: "Equal", Argsf: "len(elems), len(elems)"},
			{Fn: "Equal", Argsf: "len(elems), 1"},
			{Fn: "Equal", Argsf: "1, len(elems)"},
			{Fn: "Equal", Argsf: `nil, elems`},
			{Fn: "Equal", Argsf: `nil, b`},
			{Fn: "Equal", Argsf: `[]byte(nil), b`},

			{Fn: "NotEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "NotEqual", Argsf: "len(elems), 1"},
			{Fn: "NotEqual", Argsf: "1, len(elems)"},

			{Fn: "Greater", Argsf: "len(elems), len(elems)"},
			{Fn: "Greater", Argsf: "len(elems), 2"},
			{Fn: "Greater", Argsf: "2, len(elems)"},

			{Fn: "GreaterOrEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 0"},
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 2"},
			{Fn: "GreaterOrEqual", Argsf: "2, len(elems)"},

			{Fn: "Less", Argsf: "len(elems), len(elems)"},
			{Fn: "Less", Argsf: "len(elems), 2"},
			{Fn: "Less", Argsf: "2, len(elems)"},

			{Fn: "LessOrEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "LessOrEqual", Argsf: "0, len(elems)"},
			{Fn: "LessOrEqual", Argsf: "len(elems), 2"},
			{Fn: "LessOrEqual", Argsf: "2, len(elems)"},

			{Fn: "Equal", Argsf: "0, i"},
			{Fn: "NotEqual", Argsf: "0, i"},
			{Fn: "Empty", Argsf: "err"},
			{Fn: "Zero", Argsf: "arr"},
			{Fn: "Zero", Argsf: "arrPtr"},
			{Fn: "Zero", Argsf: "mp"},
			{Fn: "Zero", Argsf: "i"},
			{Fn: "Zero", Argsf: "ch"},

			// The linter ignores n > 1 case, because it is not exactly equivalent of NotEmpty.
			{Fn: "Greater", Argsf: "len(elems), 1"},
			{Fn: "Less", Argsf: "1, len(elems)"},
			// The linter ignores n >= 1 case, because NotEmpty in such case may impair the readability of the test.
			{Fn: "GreaterOrEqual", Argsf: "len(elems), 1"},
			{Fn: "LessOrEqual", Argsf: "1, len(elems)"},
		},
	}
}

func (EmptyTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("EmptyTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(emptyTestTmpl))
}

func (EmptyTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("EmptyTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(emptyTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const emptyTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		elems []any
		str   string
		b     []byte
	)
	{{ range $ti, $test := $.Tests }}
		// {{ $test.Name }}.
		{
			// Invalid.
			{{- range $ai, $assrn := $test.InvalidAssertions }}
				{{- range $vi, $var := $test.Vars }}
					{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
				{{- end }}
			{{- end }}
			{{ NewAssertionExpander.Expand $test.ExtraStringConv "assert" "t" nil }}

			// Valid.
			{{- range $ai, $assrn := $test.ValidAssertions }}
				{{- range $vi, $var := $test.Vars }}
					{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
				{{- end }}
			{{- end }}
		}
	{{ end -}}
}

func {{ .CheckerName.AsTestName }}_LenVarIndependence(t *testing.T) {
	var (
		elems []any
		arr    [0]int
		arrPtr *[0]int
		sl     []int
		mp     map[int]int
		str    string
		b      []byte
		ch     chan int
	)
	{{ range $vi, $var := $.LenTest.Vars }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $.LenTest.Assrn "assert" "t" (arr $var) }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	var (
		err error
		arr    [0]int
		arrPtr *[0]int
		mp     map[int]int
		i   int
		ch     chan int
		elems []any
		b     []byte
	)

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
