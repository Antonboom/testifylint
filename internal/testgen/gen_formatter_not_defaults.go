package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type FormatterNotDefaultsTestsGenerator struct{}

func (g FormatterNotDefaultsTestsGenerator) TemplateData() any {
	var (
		checker                            = checkers.NewFormatter().Name()
		reportUse                          = checker + ": use %s.%s"
		reportRemove                       = checker + ": remove unnecessary fmt.Sprintf"
		reportRemoveAndUse                 = checker + ": remove unnecessary fmt.Sprintf and use %s.%s"
		reportDoNotUseArgsWithNonStringMsg = checker +
			": using msgAndArgs with non-string first element (msg) causes panic"
		reportFailureMsgIsNotFmtString = checker +
			": failure message is not a format string, use msgAndArgs instead"
		reportEmptyMessage = checker +
			": empty message"
	)

	baseAssertions := []Assertion{
		{Fn: "Equal", Argsf: "1, 2"},
		{Fn: "Equal", Argsf: "1, 2, new(time.Time)"},
		{
			Fn:         "Equal",
			Argsf:      "1, 2, new(time.Time), 42",
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
			Fn:    "Fail",
			Argsf: `"test case failed", 1`,
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
			Fn:    "FailNow",
			Argsf: `"test case failed", 1`,
		},
		{Fn: "Equal", Argsf: `1, 2, "msg"`, ReportMsgf: reportUse, ProposedFn: "Equalf"},
		{Fn: "Equal", Argsf: `1, 2, "msg with arg %d", 42`, ReportMsgf: reportUse, ProposedFn: "Equalf"},
		{Fn: "Equal", Argsf: `1, 2, "msg with args %d %s", 42, "42"`, ReportMsgf: reportUse, ProposedFn: "Equalf"},
		// {Fn: "Equalf", Argsf: `1, 2`}, // Not compiled.
		// {Fn: "Equalf", Argsf: `1, 2, new(time.Time)`}, // Not compiled.
		{Fn: "Equalf", Argsf: `1, 2, "msg"`},
		{Fn: "Equalf", Argsf: `1, 2, "msg with arg %d", 42`},
		{Fn: "Equalf", Argsf: `1, 2, "msg with args %d %s", 42, "42"`},
	}

	sprintfAssertions := []Assertion{
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg")`,
			ReportMsgf:    reportRemoveAndUse,
			ProposedFn:    "Equalf",
			ProposedArgsf: `1, 2, "msg"`,
		},
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg with arg %d", 42)`,
			ReportMsgf:    reportRemoveAndUse,
			ProposedFn:    "Equalf",
			ProposedArgsf: `1, 2, "msg with arg %d", 42`,
		},
		{
			Fn:            "Equal",
			Argsf:         `1, 2, fmt.Sprintf("msg with args %d %s", 42, "42")`,
			ReportMsgf:    reportRemoveAndUse,
			ProposedFn:    "Equalf",
			ProposedArgsf: `1, 2, "msg with args %d %s", 42, "42"`,
		},
		{
			Fn:         "Equal",
			Argsf:      `1, 2, fmt.Sprintf("msg"), 42`,
			ReportMsgf: reportUse,
			ProposedFn: "Equalf",
		},
		{
			Fn:         "Equal",
			Argsf:      `1, 2, fmt.Sprintf("msg with arg %d", 42), "42"`,
			ReportMsgf: reportUse,
			ProposedFn: "Equalf",
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

		{Fn: "Equal", Argsf: `want, got, "boom!"`, ReportMsgf: reportUse, ProposedFn: "Equalf"},
		{Fn: "Equal", Argsf: `want, got, "%v != %v", 1, 2`, ReportMsgf: reportUse, ProposedFn: "Equalf"},
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
				Fn:         a.Fn,
				Argsf:      a.Argsf + `, "msg"`,
				ReportMsgf: reportUse,
				ProposedFn: a.Fn + "f",
			},
			Assertion{
				Fn:    a.Fn + "f",
				Argsf: a.Argsf + `, "msg"`,
			},
			Assertion{
				Fn:    a.Fn + "f",
				Argsf: a.Argsf + `, "msg with arg", 42`,
			},
			Assertion{
				Fn:    a.Fn + "f",
				Argsf: a.Argsf + `, "msg with arg %d", 42`,
			},
		)
	}

	return struct {
		CheckerName        CheckerName
		BaseAssertions     []Assertion
		SprintfAssertions  []Assertion
		EmptyMsgAssertions []Assertion
		AllAssertions      []Assertion
	}{
		CheckerName:        CheckerName(checker),
		BaseAssertions:     baseAssertions,
		SprintfAssertions:  sprintfAssertions,
		EmptyMsgAssertions: emptyMsgAssertions,
		AllAssertions:      assertions,
	}
}

func (FormatterNotDefaultsTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("FormatterNotDefaultsTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(formatterNotDefaultsTestTmpl))
}

func (FormatterNotDefaultsTestsGenerator) GoldenTemplate() Executor {
	g := strings.Replace(formatterNotDefaultsTestTmpl, "assert.Error", "assert.Errorf", 1)
	return template.Must(template.New("LenTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(g, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const formatterNotDefaultsTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error
	var args []any
	assert.Error(t, err, "Parse(%v) should fail.", args) // want "formatter: use assert\\.Errorf$"

	{{- range $ai, $assrn := $.BaseAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
	
	{{ range $ai, $assrn := $.SprintfAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_EmptyMessage(t *testing.T) {
	var want, got any
	assertObj := assert.New(t)

	{{- range $ai, $assrn := $.EmptyMsgAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assertObj" "" nil }}
	{{- end }}
}

func {{ .CheckerName.AsTestName }}_AllAssertions(t *testing.T) {
	{{- range $ai, $assrn := $.AllAssertions }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
