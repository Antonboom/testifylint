package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type FormatterTestsGenerator struct{}

func (FormatterTestsGenerator) Checker() checkers.Checker {
	return checkers.NewFormatter()
}

func (g FormatterTestsGenerator) TemplateData() any {
	reportRemove := g.Checker().Name() + ": remove unnecessary fmt.Sprintf"

	baseAssertion := Assertion{Fn: "Equal", Argsf: "1, 2"}

	sprintfAssertions := []Assertion{
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg")`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg"`,
		},
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg with arg %d", 42)`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg with arg %d", 42`,
		},
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg with args %d %s", 42, "42")`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg with args %d %s", 42, "42"`,
		},
		{
			Fn:    "Equal",
			Argsf: `1, 2, fmt.Sprintf("msg"), 42`,
		},
		{
			Fn:    "Equal",
			Argsf: `1, 2, fmt.Sprintf("msg with arg %d", 42), "42"`,
		},

		{
			Fn:            "Equalf",
			Argsf:         `1, 2, fmt.Sprintf("msg")`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg"`,
		},
		{
			Fn:            "Equalf",
			Argsf:         `1, 2, fmt.Sprintf("msg with arg %d", 42)`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg with arg %d", 42`,
		},
		{
			Fn:            "Equalf",
			Argsf:         `1, 2, fmt.Sprintf("msg with args %d %s", 42, "42")`,
			ReportMsgf:    reportRemove,
			ProposedArgsf: `1, 2, "msg with args %d %s", 42, "42"`,
		},
		{
			Fn:    "Equalf",
			Argsf: `1, 2, fmt.Sprintf("msg"), 42`,
		},
		{
			Fn:    "Equalf",
			Argsf: `1, 2, fmt.Sprintf("msg with arg %d", 42), "42"`,
		},
	}

	assertions := make([]Assertion, 0, len(allAssertions)*5)
	for _, a := range allAssertions {
		assertions = append(assertions,
			Assertion{
				Fn:    a.Fn,
				Argsf: a.Argsf,
			},
			Assertion{
				Fn:    a.Fn,
				Argsf: a.Argsf + `, "msg"`,
			},
			Assertion{
				Fn:    a.Fn + "f",
				Argsf: a.Argsf + `, "msg"`,
			},
			Assertion{
				Fn:         a.Fn + "f",
				Argsf:      a.Argsf + `, "msg with arg", 42`,
				ReportMsgf: g.Checker().Name() + ": %s.%s call has arguments but no formatting directives",
				ProposedFn: a.Fn + "f",
			},
			Assertion{
				Fn:    a.Fn + "f",
				Argsf: a.Argsf + `, "msg with arg %d", 42`,
			},
		)
	}

	return struct {
		CheckerName       CheckerName
		BaseAssertion     Assertion
		SprintfAssertions []Assertion
		AllAssertions     []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName:       CheckerName(g.Checker().Name()),
		BaseAssertion:     baseAssertion,
		SprintfAssertions: sprintfAssertions,
		AllAssertions:     assertions,
		IgnoredAssertions: []Assertion{
			{Fn: "ObjectsAreEqual", Argsf: "nil, nil"},
			{Fn: "ObjectsAreEqualValues", Argsf: "nil, nil"},
			{Fn: "ObjectsExportedFieldsAreEqual", Argsf: "nil, nil"},
		},
	}
}

func (FormatterTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("FormatterTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(formatterTestTmpl))
}

func (FormatterTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("FormatterTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(formatterTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

// NOTE(a.telyshev): The tests below are not supported by printf fork.
// They are waiting for a wonderful future.
/*
assert.Equalf(t, 1, 2, "msg with arg %.*d x", "42", 3) // want "formatter: assert\\.Equalf format %.*d uses non-int \"42\" as argument of \\*"
assert.Equalf(t, 1, 2, "msg with arg %d", "42")        // want "formatter: assert\\.Equalf format %d has arg \"42\" of wrong type string"

func assertTrue(t *testing.T, v bool, arg1 string, arg2 ...interface{}) {
	t.Helper()
	assert.Truef(t, v, arg1, arg2) // want "formatter: missing \\.\\.\\. in args forwarded to printf-like function"
}
*/

const formatterTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error
	var args []any
	assert.Error(t, err, "Parse(%v) should fail.", args)

	{{ NewAssertionExpander.FullMode.Expand $.BaseAssertion "assert" "t" nil }}

	{{ range $ai, $assrn := $.SprintfAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_PrintfChecks(t *testing.T) {
	assert.Equalf(t, 1, 2, "msg with arg", "arg")        // want "formatter: assert\\.Equalf call has arguments but no formatting directives"
	assert.Equalf(t, 1, 2, "msg with arg: %w", nil)      // want "formatter: assert\\.Equalf does not support error-wrapping directive %w"
	assert.Equalf(t, 1, 2, "msg with args %d", 42, "42") // want "formatter: assert\\.Equalf call needs 1 arg but has 2 args"
	assert.Equalf(t, 1, 2, "msg with arg %[xd", 42)      // want "formatter: assert\\.Equalf format %\\[xd is missing closing \\]"

	assert.Equalf(t, 1, 2, "msg with arg %[3]*.[2*[1]f", 1, 2, 3)  // want "formatter: assert\\.Equalf format has invalid argument index \\[2\\*\\[1\\]"
	
	assert.Equalf(t, 1, 2, "msg with arg %", 42)            // want "formatter: assert\\.Equalf format % is missing verb at end of string"
	assert.Equalf(t, 1, 2, "msg with arg %r", 42)           // want "formatter: assert\\.Equalf format %r has unknown verb r"
	assert.Equalf(t, 1, 2, "msg with arg %#s", 42)          // want "formatter: assert\\.Equalf format %#s has unrecognized flag #"
	assert.Equalf(t, 1, 2, "msg with arg %d", assertFalse)  // want "formatter: assert\\.Equalf format %d arg assertFalse is a func value, not called"

	assert.Equalf(t, 1, 2, "msg with args %s %s", "42") // want "formatter: assert\\.Equalf format %s reads arg #2, but call has 1 arg$"
}

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (suite *{{ $suiteName }}) TestFuncNameInDiagnostic() {
	require.Equalf(suite.T(), 1, 2, "msg with arg", "arg") // want "formatter: require\\.Equalf call has arguments but no formatting directives"

	suite.Require().Equalf(1, 2, "msg with arg", "arg") // want "formatter: suite\\.Require\\(\\)\\.Equalf call has arguments but no formatting directives"
	suite.Equalf(1, 2, "msg with arg", "arg")           // want "formatter: suite\\.Equalf call has arguments but no formatting directives"

	assertObj := assert.New(suite.T())
	assertObj.Equalf(1, 2, "msg with arg", "arg") // want "formatter: assertObj\\.Equalf call has arguments but no formatting directives"

	requireObj := require.New(suite.T())
	requireObj.Equalf(1, 2, "msg with arg", "arg") // want "formatter: requireObj\\.Equalf call has arguments but no formatting directives"
}

func assertFalse(t *testing.T, v bool, arg1 string, arg2 ...interface{}) {
	t.Helper()
	assert.Falsef(t, v, arg1, arg2...)
}

func {{ .CheckerName.AsTestName }}_AllAssertions(t *testing.T) {
	{{- range $ai, $assrn := $.AllAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	{{- range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "" nil }}
	{{- end }}
}
`
