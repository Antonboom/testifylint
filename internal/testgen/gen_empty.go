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

	type test struct {
		Name              string
		InvalidAssertions []Assertion
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
			Vars:  []string{"elems", "arr", "arrPtr", "sl", "mp", "str", "b", "ch", "[]byte(str)", "string(str)"},
			Assrn: Assertion{Fn: "Equal", Argsf: "0, len(%s)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "%s"},
		},
		Tests: []test{
			{
				Name: "assert.Empty cases",
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n <= 0, n == 0, n <= 0, n < 1
					// 0 >= n, 0 == n, 0 >= n, 1 > n
					{Fn: "LessOrEqual", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "GreaterOrEqual", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Len", Argsf: "elems, 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Len", Argsf: "str, 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Len", Argsf: "string(str), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Len", Argsf: "b, 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "b"},
					{Fn: "Len", Argsf: "string(b), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "Equal", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Equal", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "EqualValues", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "EqualValues", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Exactly", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Exactly", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "LessOrEqual", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "GreaterOrEqual", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Less", Argsf: "len(elems), 1", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "1, len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Zero", Argsf: "len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},

					// Empty string cases.
					{Fn: "Equal", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Equal", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Equal", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "EqualValues", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "EqualValues", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "EqualValues", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "Exactly", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Exactly", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Exactly", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},

					{Fn: "Equal", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Equal", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Equal", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "EqualValues", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "EqualValues", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "EqualValues", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "Exactly", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Exactly", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Exactly", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},

					// Simplification cases.
					{Fn: "Empty", Argsf: "len(elems)", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Empty", Argsf: "len(str)", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Empty", Argsf: "len(string(str))", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Empty", Argsf: `len([]string{"e"})`, ReportMsgf: reportRemoveLen, ProposedArgsf: `[]string{"e"}`},
					{Fn: "Empty", Argsf: "len(b)", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "b"},
					{Fn: "Empty", Argsf: "len(string(b))", ReportMsgf: reportRemoveLen, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "Empty", Argsf: "string(str)", ReportMsgf: reportRemoveConv, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Zero", Argsf: "len(elems)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "elems"},
					{Fn: "Zero", Argsf: "len(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Zero", Argsf: "len(string(str))", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Zero", Argsf: "len(b)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "b"},
					{Fn: "Zero", Argsf: "len(string(b))", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "string(b)"},
					{Fn: "Zero", Argsf: "string(str)", ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: "str"},
					{Fn: "Zero", Argsf: `len([]string{"e"})`, ReportMsgf: reportUse, ProposedFn: "Empty", ProposedArgsf: `[]string{"e"}`},
				},
				ValidAssertions: []Assertion{
					{Fn: "Empty", Argsf: "elems"},
					{Fn: "Empty", Argsf: "string(b)"},
				},
			},
			{
				Name: "assert.NotEmpty cases",
				InvalidAssertions: []Assertion{
					// n := len(elems)
					// n != 0, n > 0
					// 0 != n, 0 < n
					{Fn: "NotEqual", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqual", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqualValues", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEqualValues", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotExactly", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotExactly", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Greater", Argsf: "len(elems), 0", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Less", Argsf: "0, len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "Positive", Argsf: "len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotZero", Argsf: "len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},

					// Empty string cases.
					{Fn: "NotEqual", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqual", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqual", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotEqualValues", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqualValues", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqualValues", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotExactly", Argsf: `"", str`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotExactly", Argsf: `"", string(str)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotExactly", Argsf: `"", string(b)`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},

					{Fn: "NotEqual", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqual", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqual", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotEqualValues", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqualValues", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEqualValues", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotExactly", Argsf: "``, str", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotExactly", Argsf: "``, string(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotExactly", Argsf: "``, string(b)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},

					// Simplification cases.
					{Fn: "NotEmpty", Argsf: "len(elems)", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotEmpty", Argsf: "len(str)", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEmpty", Argsf: "len(string(str))", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotEmpty", Argsf: `len([]string{"e"})`, ReportMsgf: reportRemoveLen, ProposedArgsf: `[]string{"e"}`},
					{Fn: "NotEmpty", Argsf: "len(b)", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "b"},
					{Fn: "NotEmpty", Argsf: "len(string(b))", ReportMsgf: reportRemoveLen, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotEmpty", Argsf: "string(str)", ReportMsgf: reportRemoveConv, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotZero", Argsf: "len(elems)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "elems"},
					{Fn: "NotZero", Argsf: "len(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotZero", Argsf: "len(string(str))", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotZero", Argsf: "len(b)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "b"},
					{Fn: "NotZero", Argsf: "len(string(b))", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "string(b)"},
					{Fn: "NotZero", Argsf: "string(str)", ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: "str"},
					{Fn: "NotZero", Argsf: `len([]string{"e"})`, ReportMsgf: reportUse, ProposedFn: "NotEmpty", ProposedArgsf: `[]string{"e"}`},
				},
				ValidAssertions: []Assertion{
					{Fn: "NotEmpty", Argsf: "elems"},
					{Fn: "NotEmpty", Argsf: "string(b)"},
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
			{Fn: "Zero", Argsf: "err"},

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
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}

			// Valid.
			{{- range $ai, $assrn := $test.ValidAssertions }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
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
		elems []any
		str   string
		b     []byte
		i   int
	)

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
