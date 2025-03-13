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
		checker         = g.Checker().Name()
		reportUse       = checker + ": use %s.%s"
		reportRemoveLen = checker + ": remove unnecessary len"
	)

	vars := []string{"elems", "str", "string(str)", "b", "string(b)", `[]string{"e"}`}

	type test struct {
		Name              string
		Vars              []string
		InvalidAssertions []Assertion
		Special           []Assertion
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
					{Fn: "Empty", Argsf: "len(%s)", ReportMsgf: reportRemoveLen, ProposedArgsf: "%s"},
				},
				Special: []Assertion{
					// Zero is moved to separate cases, because not-string vars are not relevant for it.
					{Fn: "Zero", Argsf: "str", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Zero", Argsf: "string(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(str)"},
					{Fn: "Zero", Argsf: "string(b)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
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
					{Fn: "NotEmpty", Argsf: "len(%s)", ReportMsgf: reportRemoveLen, ProposedArgsf: "%s"},
				},
				Special: []Assertion{
					// NotZero is moved to separate cases, because not-string vars are not relevant for it.
					{Fn: "NotZero", Argsf: "str", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotZero", Argsf: "string(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(str)"},
					{Fn: "NotZero", Argsf: "string(b)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
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
			{Fn: "Equal", Argsf: "nil, elems"},
			{Fn: "Equal", Argsf: "nil, b"},
			{Fn: "Nil", Argsf: "elems"},
			{Fn: "Nil", Argsf: "b"},
			{Fn: "Equal", Argsf: "[]byte(nil), b"},

			{Fn: "NotEqual", Argsf: "len(elems), len(elems)"},
			{Fn: "NotEqual", Argsf: "len(elems), 1"},
			{Fn: "NotEqual", Argsf: "1, len(elems)"},
			{Fn: "NotEqual", Argsf: "nil, elems"},
			{Fn: "NotEqual", Argsf: "nil, b"},
			{Fn: "NotNil", Argsf: "elems"},
			{Fn: "NotNil", Argsf: "b"},
			{Fn: "NotEqual", Argsf: "[]byte(nil), b"},

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
			{Fn: "NotEmpty", Argsf: "err"},

			// Zero and Empty are not equal assertions sometimes, be careful!
			{Fn: "Zero", Argsf: "elems"},
			{Fn: "Zero", Argsf: "arr"},
			{Fn: "Zero", Argsf: "arrPtr"},
			{Fn: "Zero", Argsf: "mp"},
			{Fn: "Zero", Argsf: "b"},
			{Fn: "Zero", Argsf: "i"},
			{Fn: "Zero", Argsf: "ch"},
			{Fn: "Zero", Argsf: `[]string{"e"}`},
			{Fn: "NotZero", Argsf: "elems"},
			{Fn: "NotZero", Argsf: "arr"},
			{Fn: "NotZero", Argsf: "arrPtr"},
			{Fn: "NotZero", Argsf: "mp"},
			{Fn: "NotZero", Argsf: "b"},
			{Fn: "NotZero", Argsf: "i"},
			{Fn: "NotZero", Argsf: "ch"},
			{Fn: "NotZero", Argsf: `[]string{"e"}`},

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
			{{- range $ai, $assrn := $test.Special }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}

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
