package main

import (
	"github.com/Antonboom/testifylint/internal/checkers"
	"strings"
	"text/template"
)

type FormatterTestsGenerator struct{}

func (FormatterTestsGenerator) Checker() checkers.Checker {
	return checkers.NewFormatter()
}

func (g FormatterTestsGenerator) TemplateData() any {

	// TODO
	// 1) Тесты на use *f, подробные
	// 2) По тесту на каждый кейс из современного printf
	//		- тест на другие объекты
	// 3) Тесты на все функции, но по 1 кейсу из 1 и 2
	// 4) Positive tests (ignores)
	// 5) Тест на то, чтобы перезатирать printf и не ругаться на не-testify функции
	// 6) Тесты на fmt.Srpintf
	// 7) Autofix для 6) и 1)
	// 8) Ввести в readme понятие частичного автофикса (+-)
	// 9) тест на несколько варнингов в рамках одного ассерта

	baseAssertions := []Assertion{
		{Fn: "Condition", Argsf: "nil"},
		{Fn: "Contains", Argsf: "nil, nil"},
		{Fn: "DirExists", Argsf: `""`},
		{Fn: "ElementsMatch", Argsf: "nil, nil"},
		{Fn: "Empty", Argsf: "nil"},
		{Fn: "Equal", Argsf: "nil, nil"},
		{Fn: "EqualError", Argsf: `nil, ""`},
		{Fn: "EqualExportedValues", Argsf: "nil, nil"},
		{Fn: "EqualValues", Argsf: "nil, nil"},
		{Fn: "Error", Argsf: "nil"},
		{Fn: "ErrorAs", Argsf: "nil, nil"},
		{Fn: "ErrorContains", Argsf: `nil, ""`},
		{Fn: "ErrorIs", Argsf: "nil, nil"},
		{Fn: "Eventually", Argsf: "nil, 0, 0"},
		{Fn: "EventuallyWithT", Argsf: "nil, 0, 0"},
		{Fn: "Exactly", Argsf: "nil, nil"},
		{Fn: "Fail", Argsf: `""`},
		{Fn: "FailNow", Argsf: `""`},
		{Fn: "False", Argsf: "false"},
		{Fn: "FileExists", Argsf: `""`},
		{Fn: "Implements", Argsf: "nil, nil"},
		{Fn: "InDelta", Argsf: "0., 0., 0."},
		{Fn: "InDeltaMapValues", Argsf: "nil, nil, 0."},
		{Fn: "InDeltaSlice", Argsf: "nil, nil, 0."},
		{Fn: "InEpsilon", Argsf: "nil, nil, 0."},
		{Fn: "InEpsilonSlice", Argsf: "nil, nil, 0."},
		{Fn: "IsType", Argsf: "nil, nil"},
		{Fn: "JSONEq", Argsf: `"", ""`},
		{Fn: "Len", Argsf: "nil, 0"},
		{Fn: "Never", Argsf: "nil, 0, 0"},
		{Fn: "Nil", Argsf: "nil"},
		{Fn: "NoDirExists", Argsf: `""`},
		{Fn: "NoError", Argsf: "nil"},
		{Fn: "NoFileExists", Argsf: `""`},
		{Fn: "NotContains", Argsf: "nil, nil"},
		{Fn: "NotEmpty", Argsf: "nil"},
		{Fn: "NotEqual", Argsf: "nil, nil"},
		{Fn: "NotEqualValues", Argsf: "nil, nil"},
		{Fn: "NotErrorIs", Argsf: "nil, nil"},
		{Fn: "NotNil", Argsf: "nil"},
		{Fn: "NotPanics", Argsf: "nil"},
		{Fn: "NotRegexp", Argsf: `nil, ""`},
		{Fn: "NotSame", Argsf: "nil, nil"},
		{Fn: "NotSubset", Argsf: "nil, nil"},
		{Fn: "NotZero", Argsf: "nil"},
		{Fn: "Panics", Argsf: "nil"},
		{Fn: "PanicsWithError", Argsf: `"", nil`},
		{Fn: "PanicsWithValue", Argsf: "nil, nil"},
		{Fn: "Regexp", Argsf: "nil, nil"},
		{Fn: "Same", Argsf: "nil, nil"},
		{Fn: "Subset", Argsf: "nil, nil"},
		{Fn: "True", Argsf: "true"},
		{Fn: "WithinDuration", Argsf: "time.Time{}, time.Time{}, 0"},
		{Fn: "WithinRange", Argsf: "time.Time{}, time.Time{}, time.Time{}"},
		{Fn: "YAMLEq", Argsf: `"", ""`},
		{Fn: "Zero", Argsf: "nil"},
	}

	assertions := make([]Assertion, 0, len(baseAssertions)*3)
	for _, a := range baseAssertions {
		assertions = append(assertions,
			Assertion{
				Fn:    a.Fn,
				Argsf: a.Argsf,
			},
			Assertion{
				Fn:         a.Fn,
				Argsf:      a.Argsf + `, "simple msg"`,
				ReportMsgf: g.Checker().Name() + ": use %s.%s",
				ProposedFn: a.Fn + "f",
			},
			Assertion{
				Fn:         a.Fn + "f",
				Argsf:      a.Argsf + `, "msg with arg", "arg"`,
				ReportMsgf: g.Checker().Name() + ": %s.%s call has arguments but no formatting directives",
				ProposedFn: a.Fn + "f",
			},
		)
	}

	return struct {
		CheckerName       CheckerName
		Assertions        []Assertion
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(g.Checker().Name()),
		Assertions:  assertions,
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
	return template.Must(template.New("LenTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(formatterTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const formatterTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestFormatterChecker(t *testing.T) {
	var err error
	var args []any
	assert.Error(t, err, "Parse(%v) should fail.", args) // want "formatter: use assert\\.Errorf"

	assert.Equal(t, 1, 2)
	assert.Equal(t, 1, 2, "msg")                           // want "formatter: use assert\\.Equalf"
	assert.Equal(t, 1, 2, "msg with arg %d", 42)           // want "formatter: use assert\\.Equalf"
	assert.Equal(t, 1, 2, "msg with args %d %s", 42, "42") // want "formatter: use assert\\.Equalf"
	assert.Equalf(t, 1, 2, "msg")
	assert.Equalf(t, 1, 2, "msg with arg %d", 42)
	assert.Equalf(t, 1, 2, "msg with args %d %s", 42, "42")
}

func {{ .CheckerName.AsTestName }}_PrintfChecks(t *testing.T) {
	assert.Equalf(t, 1, 2, "msg with arg", "arg")        // want "formatter: assert\\.Equalf call has arguments but no formatting directives"
	assert.Equalf(t, 1, 2, "msg with arg: %w", nil)      // want "formatter: assert\\.Equalf does not support error-wrapping directive %w"
	assert.Equalf(t, 1, 2, "msg with args %d", 42, "42") // want "formatter: assert\\.Equalf call needs 1 arg but has 2 args"
	assert.Equalf(t, 1, 2, "msg with arg %[xd", 42)      // want "formatter: assert\\.Equalf format %\\[xd is missing closing \\]"

	assert.Equalf(t, 1, 2, "msg with arg %[3]*.[2*[1]f", 1, 2, 3)  // want "formatter: assert\\.Equalf format has invalid argument index \\[2\\*\\[1\\]"
	
	assert.Equalf(t, 1, 2, "msg with arg %", 42)           // want "formatter: assert\\.Equalf format % is missing verb at end of string"
	assert.Equalf(t, 1, 2, "msg with arg %r", 42)          // want "formatter: assert\\.Equalf format %r has unknown verb r"
	assert.Equalf(t, 1, 2, "msg with arg %#s", 42)         // want "formatter: assert\\.Equalf format %#s has unrecognized flag #"
	assert.Equalf(t, 1, 2, "msg with arg %.*d x", "42", 3) // want "formatter: assert\\.Equalf format %.*d uses non-int \"42\" as argument of \\*"
	assert.Equalf(t, 1, 2, "msg with arg %d", assertTrue)  // want "formatter: assert\\.Equalf format %d arg assertTrue is a func value, not called"
	assert.Equalf(t, 1, 2, "msg with arg %d", "42")        // want "formatter: assert\\.Equalf format %d has arg \"42\" of wrong type string"

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

func assertTrue(t *testing.T, v bool, arg1 string, arg2 ...interface{}) {
	t.Helper()
	assert.Truef(t, v, arg1, arg2) // want "formatter: missing \\.\\.\\. in args forwarded to printf-like function"
}

func {{ .CheckerName.AsTestName }}_AllAssertions(t *testing.T) {
	{{- range $ai, $assrn := $.Assertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	{{- range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "" nil }}
	{{- end }}
}
`
