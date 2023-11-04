package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorNilTestsGenerator struct{}

func (ErrorNilTestsGenerator) Checker() checkers.Checker {
	return checkers.NewErrorNil()
}

func (g ErrorNilTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	type errorDetectionTest struct {
		Vars  []string
		Assrn Assertion
	}

	type validNilsTest struct {
		Vars       []string
		Assertions []Assertion
	}

	return struct {
		CheckerName       CheckerName
		ErrDetection      errorDetectionTest
		ValidNils         validNilsTest
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		ErrDetection: errorDetectionTest{
			Vars:  []string{"a", "b.Get1()", "c", "errOp()"},
			Assrn: Assertion{Fn: "Nil", Argsf: "%s", ReportMsgf: report, ProposedFn: "NoError"},
		},
		ValidNils: validNilsTest{
			Vars: []string{"ptr", "iface", "ch", "sl", "fn", "m", "uPtr"},
			Assertions: []Assertion{
				{Fn: "Nil", Argsf: "%s"},
				{Fn: "NotNil", Argsf: "%s"},
			},
		},
		InvalidAssertions: []Assertion{
			{Fn: "Nil", Argsf: "err", ReportMsgf: report, ProposedFn: "NoError"},
			{Fn: "NotNil", Argsf: "err", ReportMsgf: report, ProposedFn: "Error"},
			{Fn: "Equal", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "Equal", Argsf: "nil, err", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "EqualValues", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "EqualValues", Argsf: "nil, err", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "Exactly", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "Exactly", Argsf: "nil, err", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "NotEqual", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
			{Fn: "NotEqual", Argsf: "nil, err", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
			{Fn: "NotEqualValues", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
			{Fn: "NotEqualValues", Argsf: "nil, err", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
			{Fn: "ErrorIs", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "NoError", ProposedArgsf: "err"},
			{Fn: "NotErrorIs", Argsf: "err, nil", ReportMsgf: report, ProposedFn: "Error", ProposedArgsf: "err"},
		},
		ValidAssertions: []Assertion{
			{Fn: "NoError", Argsf: "err"},
			{Fn: "Error", Argsf: "err"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Nil", Argsf: "nil"},
			{Fn: "NotNil", Argsf: "nil"},
			{Fn: "Equal", Argsf: "err, err"},
			{Fn: "Equal", Argsf: "nil, nil"},
			{Fn: "NotEqual", Argsf: "err, err"},
			{Fn: "NotEqual", Argsf: "nil, nil"},
		},
	}
}

func (ErrorNilTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ErrorNilTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorNilTestTmpl))
}

func (ErrorNilTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ErrorNilTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(errorNilTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const errorNilTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error

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

func {{ .CheckerName.AsTestName }}_ErrorDetection(t *testing.T) {
	errOp := func() error { return io.EOF }
	var a error
	var b withErroredMethod
	_, c := b.Get2()
	{{ range $vi, $var := $.ErrDetection.Vars }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $.ErrDetection.Assrn "assert" "t" (arr $var) }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_ValidNils(t *testing.T) {
	var (
		ptr   *int
		iface any
		ch    chan error
		sl    []error
		fn    func()
		m     map[int]int
		uPtr unsafe.Pointer
	)
	{{ range $vi, $var := $.ValidNils.Vars }}
		{{- range $ai, $assrn := $.ValidNils.Assertions }}
			{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" (arr $var) }}
		{{- end }}
	{{- end }}	
}

type withErroredMethod struct{}

func (withErroredMethod) Get1() error        { return nil }
func (withErroredMethod) Get2() (int, error) { return 0, nil }
`
