package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ContainsTestsGenerator struct{}

func (ContainsTestsGenerator) Checker() checkers.Checker {
	return checkers.NewContains()
}

func (g ContainsTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	return struct {
		CheckerName        CheckerName
		Vars               []string
		InvalidAssertions  []Assertion
		InvalidSubsetCases []Assertion
		ValidAssertions    []Assertion
		IgnoredAssertions  []Assertion
	}{
		CheckerName: CheckerName(checker),
		Vars:        []string{"s", "string(b)"},
		InvalidAssertions: []Assertion{
			{
				Fn:            "True",
				Argsf:         `strings.Contains(%s, "abc123")`,
				ReportMsgf:    report,
				ProposedFn:    "Contains",
				ProposedArgsf: `%s, "abc123"`,
			},
			{
				Fn:            "False",
				Argsf:         `!strings.Contains(%s, "abc123")`,
				ReportMsgf:    report,
				ProposedFn:    "Contains",
				ProposedArgsf: `%s, "abc123"`,
			},

			{
				Fn:            "False",
				Argsf:         `strings.Contains(%s, "abc123")`,
				ReportMsgf:    report,
				ProposedFn:    "NotContains",
				ProposedArgsf: `%s, "abc123"`,
			},
			{
				Fn:            "True",
				Argsf:         `!strings.Contains(%s, "abc123")`,
				ReportMsgf:    report,
				ProposedFn:    "NotContains",
				ProposedArgsf: `%s, "abc123"`,
			},
		},
		InvalidSubsetCases: []Assertion{
			{
				Fn:         "Contains",
				Argsf:      `metrics, metric{time: 1}, metric{time: 2}`,
				ReportMsgf: checker + ": invalid usage of assert.Contains, use assert.Subset for multi elements assertion",
			},
			{
				Fn:         "NotContains",
				Argsf:      `metrics, metric{time: 1}, metric{time: 2}`,
				ReportMsgf: checker + ": invalid usage of assert.NotContains, use assert.NotSubset for multi elements assertion",
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Contains", Argsf: `s, "abc123"`},
			{Fn: "NotContains", Argsf: `string(b), "abc123"`},

			{Fn: "Subset", Argsf: `metrics, []metric{{time: 1}, {time: 2}}`},
			{Fn: "NotSubset", Argsf: `metrics, []metric{{time: 1}, {time: 2}}`},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Contains", Argsf: `errSentinel.Error(), "user"`},    // error-compare case.
			{Fn: "NotContains", Argsf: `errSentinel.Error(), "user"`}, // error-compare case.

			{Fn: "Contains", Argsf: `string(b), "MASKED_KEY=[MASKED]"`},
			{Fn: "Contains", Argsf: `metrics, metrics`},
			{Fn: "Contains", Argsf: `metrics, 1, fmt.Sprintf("should contain %d", 1)`},
			{Fn: "NotContains", Argsf: `string(b), "MASKED_KEY=[MASKED]"`},
			{Fn: "NotContains", Argsf: `metrics, metrics`},
			{Fn: "NotContains", Argsf: `metrics, 1, fmt.Sprintf("should contain %d", 1)`},

			// https://github.com/Antonboom/testifylint/issues/154
			{Fn: "True", Argsf: `bytes.Contains(b, []byte("a"))`},
			{Fn: "False", Argsf: `!bytes.Contains(b, []byte("a"))`},
			{Fn: "False", Argsf: `bytes.Contains(b, []byte("a"))`},
			{Fn: "True", Argsf: `!bytes.Contains(b, []byte("a"))`},
		},
	}
}

func (ContainsTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ContainsTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(containsTestTmpl))
}

func (ContainsTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ContainsTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(containsTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const containsTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
    "bytes"
    "errors"
    "fmt"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

type metric struct {
    time int
}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
    var (
        s           = "abc123"
        b           = []byte(s)
        errSentinel = errors.New("user not found")
		metrics     = []metric{}
    )

    // Invalid.
    {
        {{- range $ai, $assrn := $.InvalidAssertions }}
            {{- range $vi, $var := $.Vars }}
                {{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
            {{- end }}
        {{- end }}

        {{- range $ai, $assrn := $.InvalidSubsetCases }}
            {{ NewAssertionExpander.NotFmtSetMode.Expand $assrn "assert" "t" nil }}
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

// ErrorContains returns an assertion to check if the error contains the given string.
func ErrorContains(contains string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
		return assert.Contains(t, err.Error(), contains, msgAndArgs...)
	}
}
`
