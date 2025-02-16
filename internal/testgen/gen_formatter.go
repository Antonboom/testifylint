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
	var (
		checker      = g.Checker().Name()
		reportRemove = checker +
			": remove unnecessary fmt.Sprintf"
		reportDoNotUseNonStringMsg = checker +
			": do not use non-string value as first element (msg) of msgAndArgs"
		reportDoNotUseArgsWithNonStringMsg = checker +
			": using msgAndArgs with non-string first element (msg) causes panic"
		reportFailureMsgIsNotFmtString = checker +
			": failure message is not a format string, use msgAndArgs instead"
		reportEmptyMessage = checker +
			": empty message"
	)

	baseAssertion := Assertion{Fn: "Equal", Argsf: "1, 2"}

	nonStringMsgAssertions := []Assertion{
		{
			Fn:            "Equal",
			Argsf:         "1, 2, new(time.Time)",
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `1, 2, "%+v", new(time.Time)`,
		},
		{
			Fn:            "Equal",
			Argsf:         "1, 2, i",
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `1, 2, "%+v", i`,
		},
		{
			Fn:            "Equal",
			Argsf:         "1, 2, tc",
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `1, 2, "%+v", tc`,
		},
		{
			Fn:            "Equal",
			Argsf:         "1, 2, args",
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `1, 2, "%+v", args`,
		},

		{
			Fn:         "Equal",
			Argsf:      "1, 2, new(time.Time), 42",
			ReportMsgf: reportDoNotUseArgsWithNonStringMsg,
		},
		{
			Fn:         "Equal",
			Argsf:      `1, 2, i, 42, "42"`,
			ReportMsgf: reportDoNotUseArgsWithNonStringMsg,
		},
		{
			Fn:         "Equal",
			Argsf:      "1, 2, tc, 0",
			ReportMsgf: reportDoNotUseArgsWithNonStringMsg,
		},
		{
			Fn:         "Fail",
			Argsf:      `"test case [%d] failed.  Expected: %+v, Got: %+v", 1, 2, 3`,
			ReportMsgf: reportFailureMsgIsNotFmtString,
		},
		{
			Fn:         "Fail",
			Argsf:      `"test case [%d] failed", 1`,
			ReportMsgf: reportFailureMsgIsNotFmtString,
		},
		{
			Fn:            "Fail",
			Argsf:         `"test case failed", 1`,
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `"test case failed", "%+v", 1`,
		},
		{
			Fn:         "FailNow",
			Argsf:      `"test case [%d] failed.  Expected: %+v, Got: %+v", 1, 2, 3`,
			ReportMsgf: reportFailureMsgIsNotFmtString,
		},
		{
			Fn:         "FailNow",
			Argsf:      `"test case [%d] failed", 1`,
			ReportMsgf: reportFailureMsgIsNotFmtString,
		},
		{
			Fn:            "FailNow",
			Argsf:         `"test case failed", 1`,
			ReportMsgf:    reportDoNotUseNonStringMsg,
			ProposedArgsf: `"test case failed", "%+v", 1`,
		},

		{Fn: "Equal", Argsf: "1, 2, msg()"},
		{Fn: "Equal", Argsf: "1, 2, new(time.Time).String()"},
		{Fn: "Equal", Argsf: `1, 2, args...`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", new(time.Time)`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", i`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", tc`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", msg()`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", new(time.Time).String()`},
		{Fn: "Equal", Argsf: `1, 2, "%+v", args`},
		{Fn: "Fail", Argsf: `"boom!", "test case [%d] failed", 1`},
		{Fn: "FailNow", Argsf: `"boom!", "test case [%d] failed", 1`},
	}

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

	emptyMsgAssertions := []Assertion{
		{
			Fn: "Equal", Argsf: `want, got, ""`,
			ReportMsgf: reportEmptyMessage, ProposedArgsf: "want, got",
		},
		{
			Fn: "Equalf", Argsf: `want, got, ""`,
			ReportMsgf: reportEmptyMessage + "%.s%.s", ProposedFn: "Equal", ProposedArgsf: "want, got",
		},
		{
			Fn: "Equal", Argsf: `want, got, "", 1, 2`,
			ReportMsgf: reportEmptyMessage,
		},
		{
			Fn: "Equalf", Argsf: `want, got, "", 1, 2`,
			ReportMsgf: reportEmptyMessage,
		},

		{Fn: "Equal", Argsf: `want, got, "boom!"`},
		{Fn: "Equal", Argsf: `want, got, "%v != %v", 1, 2`},
		{Fn: "Equalf", Argsf: `want, got, "%v != %v", 1, 2`},
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
		CheckerName            CheckerName
		BaseAssertion          Assertion
		NonStringMsgAssertions []Assertion
		SprintfAssertions      []Assertion
		EmptyMsgAssertions     []Assertion
		AllAssertions          []Assertion
		IgnoredAssertions      []Assertion
	}{
		CheckerName:            CheckerName(g.Checker().Name()),
		BaseAssertion:          baseAssertion,
		NonStringMsgAssertions: nonStringMsgAssertions,
		SprintfAssertions:      sprintfAssertions,
		EmptyMsgAssertions:     emptyMsgAssertions,
		AllAssertions:          assertions,
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

	assertObj := assert.New(t)

	var i int
	var tc struct {
		strLogLevel        string
		logFunc            func(func() string)
		expectedToBeCalled bool
	}
	{{- range $ai, $assrn := $.NonStringMsgAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assertObj" "" nil }}
	{{- end }}

	{{ range $ai, $assrn := $.SprintfAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assertObj" "" nil }}
	{{- end }}
}

func msg() string {
	return "msg"
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

func {{ .CheckerName.AsTestName }}_EmptyMessage(t *testing.T) {
	var want, got any
	assertObj := assert.New(t)

	{{- range $ai, $assrn := $.EmptyMsgAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assertObj" "" nil }}
	{{- end }}
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

// AssertElementsMatchFunc is similar to the assert.ElementsMatch function, but it allows for a custom comparison function
func AssertElementsMatchFunc[T any](t *testing.T, expected, actual []T, comp func(a, b T) bool, msgAndArgs ...any) bool {
	t.Helper()

	extraA, extraB := diffListsFunc(expected, actual, comp)
	if len(extraA) == 0 && len(extraB) == 0 {
		return true
	}

	return assert.Fail(t, formatListDiff(expected, actual, extraA, extraB), msgAndArgs...)
}

func diffListsFunc[T any](a []T, b []T, comp func(a, b T) bool) ([]T, []T) { return nil, nil }
func formatListDiff[T any](expected, actual, extraExpected, extraActual []T) string { return "" }
`
