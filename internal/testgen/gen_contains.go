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
		CheckerName       CheckerName
		Vars              []string
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
		IgnoredAssertions []Assertion
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
		ValidAssertions: []Assertion{
			{Fn: "Contains", Argsf: `%s, "abc123"`},
			{Fn: "NotContains", Argsf: `%s, "abc123"`},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Contains", Argsf: `errSentinel.Error(), "user"`}, // error-compare case.

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
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
    var (
        s           = "abc123"
        b           = []byte(s)
        errSentinel = errors.New("user not found")
    )

    // Invalid.
    {
        {{- range $ai, $assrn := $.InvalidAssertions }}
            {{- range $vi, $var := $.Vars }}
                {{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
            {{- end }}
        {{- end }}
    }

    // Valid.
    {
        {{- range $ai, $assrn := $.ValidAssertions }}
            {{- range $vi, $var := $.Vars }}
                {{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
            {{- end }}
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
