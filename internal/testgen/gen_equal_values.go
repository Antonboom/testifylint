package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type EqualValuesTestsGenerator struct{}

func (EqualValuesTestsGenerator) Checker() checkers.Checker {
	return checkers.NewEqualValues()
}

func (g EqualValuesTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "EqualValues", Argsf: "42, i", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "int8(42), i8", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "int16(42), i16", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "int32(42), i32", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "int64(42), i64", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "uint(42), ui", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "uint8(42), ui8", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "uint16(42), ui16", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "uint32(42), ui32", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: "uint64(42), ui64", ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `"42", str`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `req, req`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `struct{ int }{}, ss`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `map[any]any{}, mm`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `[]byte(nil), b`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `map[string]string(nil), m`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "EqualValues", Argsf: `(*tls.Config)(nil), tlsConf`, ReportMsgf: report, ProposedFn: "Equal"},
			{Fn: "NotEqualValues", Argsf: "42, i", ReportMsgf: report, ProposedFn: "NotEqual"},
		},
		ValidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "42, i"},
			{Fn: "Equal", Argsf: "int8(42), i8"},
			{Fn: "Equal", Argsf: "int16(42), i16"},
			{Fn: "Equal", Argsf: "int32(42), i32"},
			{Fn: "Equal", Argsf: "int64(42), i64"},
			{Fn: "Equal", Argsf: "uint(42), ui"},
			{Fn: "Equal", Argsf: "uint8(42), ui8"},
			{Fn: "Equal", Argsf: "uint16(42), ui16"},
			{Fn: "Equal", Argsf: "uint32(42), ui32"},
			{Fn: "Equal", Argsf: "uint64(42), ui64"},
			{Fn: "Equal", Argsf: `"42", str`},
			{Fn: "Equal", Argsf: `req, req`},
			{Fn: "Equal", Argsf: `struct{ int }{}, ss`},
			{Fn: "Equal", Argsf: `map[any]any{}, mm`},
			{Fn: "Equal", Argsf: `[]byte(nil), b`},
			{Fn: "Equal", Argsf: `map[string]string(nil), m`},
			{Fn: "Equal", Argsf: `(*tls.Config)(nil), tlsConf`},
			{Fn: "NotEqual", Argsf: "42, i"},

			{Fn: "EqualValues", Argsf: `2048, mm["Etype"]`},
			{Fn: "EqualValues", Argsf: `req, dto`},
			{Fn: "EqualValues", Argsf: `42, fortyTwoAny`},
			{Fn: "EqualValues", Argsf: `42, fortyTwoAny`},
			{Fn: "EqualValues", Argsf: `42, fortyTwoInterface`},
			{Fn: "EqualValues", Argsf: `42, fortyTwoAnyAlias`},
			{Fn: "EqualValues", Argsf: `fortyTwoAny, fortyTwoInterface`},
			{Fn: "EqualValues", Argsf: `req, reqWithTags`},
			{Fn: "EqualValues", Argsf: `S{"1"}, []string{"1"}`},
			{Fn: "EqualValues", Argsf: `f, (func())(nil)`},
			{Fn: "EqualValues", Argsf: `(func())(nil), f`},
			{Fn: "NotEqualValues", Argsf: `2048, mm["Etype"]`},
		},
	}
}

func (EqualValuesTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("EqualValuesTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(equalValuesTestTmpl))
}

func (EqualValuesTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("EqualValuesTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(equalValuesTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const equalValuesTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

type S []string

type Request struct { ID string }
type RequestWithTags struct { ID string ` + "`json:\"name\"` }" + `
type Arg struct { ID string }
type customAnyAlias any

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64
	)
	var (
		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64
	)
	var (
		str string
		req Request
		reqWithTags RequestWithTags
		dto Arg
		ss  struct{ int }
		m   map[string]string
		mm  map[any]any
		b   []byte
		f   func() bool
	)

	var (
		fortyTwoAny       any            = 42
		fortyTwoInterface interface{}    = float32(42.0)
		fortyTwoAnyAlias  customAnyAlias = uint8(42)
	)

	tlsConf := new(tls.Config)

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
}
`
