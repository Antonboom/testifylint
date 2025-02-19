package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type RegexpTestsGenerator struct{}

func (RegexpTestsGenerator) Checker() checkers.Checker {
	return checkers.NewRegexp()
}

func (g RegexpTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": remove unnecessary regexp.MustCompile"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{
				Fn: "Regexp", Argsf: "regexp.MustCompile(`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`), out",
				ReportMsgf: report, ProposedArgsf: "`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`, out",
			},

			{
				Fn: "NotRegexp", Argsf: "regexp.MustCompile(`\\[.*\\] TRACE message`), out",
				ReportMsgf: report, ProposedArgsf: "`\\[.*\\] TRACE message`, out",
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Regexp", Argsf: "`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`, out"},
			{Fn: "NotRegexp", Argsf: "`\\[.*\\] TRACE message`, out"},
		},
	}
}

func (RegexpTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("RegexpTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(regexpTestTmpl))
}

func (RegexpTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("RegexpTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(regexpTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const regexpTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var out string

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
