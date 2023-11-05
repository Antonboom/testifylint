package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type NilCompareTestsGenerator struct{}

func (NilCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewNilCompare()
}

func (g NilCompareTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "value, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "Equal", Argsf: "nil, value", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "Equal", Argsf: `Row["col"], nil`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "Equal", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},

			{Fn: "EqualValues", Argsf: "value, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "EqualValues", Argsf: "nil, value", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "EqualValues", Argsf: `Row["col"], nil`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "EqualValues", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},

			{Fn: "Exactly", Argsf: "value, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "Exactly", Argsf: "nil, value", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "value"},
			{Fn: "Exactly", Argsf: `Row["col"], nil`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "Exactly", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},

			{Fn: "NotEqual", Argsf: "value, nil", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "value"},
			{Fn: "NotEqual", Argsf: "nil, value", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "value"},
			{Fn: "NotEqual", Argsf: `Row["col"], nil`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},
			{Fn: "NotEqual", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},

			{Fn: "NotEqualValues", Argsf: "value, nil", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "value"},
			{Fn: "NotEqualValues", Argsf: "nil, value", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "value"},
			{Fn: "NotEqualValues", Argsf: `Row["col"], nil`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},
			{Fn: "NotEqualValues", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},
		},
		ValidAssertions: []Assertion{
			{Fn: "Nil", Argsf: "value"},
			{Fn: "NotNil", Argsf: "value"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "value, value"},
			{Fn: "Equal", Argsf: "nil, nil"},
			{Fn: "Equal", Argsf: `Row["col"], "foo"`},
			{Fn: "Equal", Argsf: `"foo", Row["col"]`},
			{Fn: "Equal", Argsf: `Row["col"], Row["col"]`},

			{Fn: "EqualValues", Argsf: "value, value"},
			{Fn: "EqualValues", Argsf: "nil, nil"},
			{Fn: "EqualValues", Argsf: `Row["col"], "foo"`},
			{Fn: "EqualValues", Argsf: `"foo", Row["col"]`},
			{Fn: "EqualValues", Argsf: `Row["col"], Row["col"]`},

			{Fn: "Exactly", Argsf: "value, value"},
			{Fn: "Exactly", Argsf: "nil, nil"},
			{Fn: "Exactly", Argsf: `Row["col"], "foo"`},
			{Fn: "Exactly", Argsf: `"foo", Row["col"]`},
			{Fn: "Exactly", Argsf: `Row["col"], Row["col"]`},

			{Fn: "NotEqual", Argsf: "value, value"},
			{Fn: "NotEqual", Argsf: "nil, nil"},
			{Fn: "NotEqual", Argsf: `Row["col"], "foo"`},
			{Fn: "NotEqual", Argsf: `"foo", Row["col"]`},
			{Fn: "NotEqual", Argsf: `Row["col"], Row["col"]`},

			{Fn: "NotEqualValues", Argsf: "value, value"},
			{Fn: "NotEqualValues", Argsf: "nil, nil"},
			{Fn: "NotEqualValues", Argsf: `Row["col"], "foo"`},
			{Fn: "NotEqualValues", Argsf: `"foo", Row["col"]`},
			{Fn: "NotEqualValues", Argsf: `Row["col"], Row["col"]`},
		},
	}
}

func (NilCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("NilCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(nilCompareTestTmpl))
}

func (NilCompareTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("NilCompareTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(nilCompareTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const nilCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var value any
	var Row map[string]any

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
