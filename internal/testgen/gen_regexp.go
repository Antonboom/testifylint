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
		checker           = g.Checker().Name()
		reportMustCompile = checker + ": remove unnecessary regexp.MustCompile"
		reportInvalidArg  = checker + ": use string or *regexp.Regexp as the first argument"
	)

	return struct {
		CheckerName           CheckerName
		InvalidMustCompile    []Assertion
		InvalidTypeAssertions []Assertion
		ValidAssertions       []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidMustCompile: []Assertion{
			{
				Fn: "Regexp", Argsf: "regexp.MustCompile(`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`), out",
				ReportMsgf: reportMustCompile, ProposedArgsf: "`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`, out",
			},
			{
				Fn: "NotRegexp", Argsf: "regexp.MustCompile(`\\[.*\\] TRACE message`), out",
				ReportMsgf: reportMustCompile, ProposedArgsf: "`\\[.*\\] TRACE message`, out",
			},
		},
		InvalidTypeAssertions: []Assertion{
			{
				Fn:         "Regexp",
				Argsf:      "[]byte(`\\w+`), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "[]byte(`\\w+`), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "myDefinedString(`\\w+`), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "myDefinedString(`\\w+`), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "rxRegexpStructAlias, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "rxRegexpStructAlias, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "42, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "42, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "int8(42), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "int8(42), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "uint8(42), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "uint8(42), out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "42.0, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "42.0, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "Regexp",
				Argsf:      "testStringer{}, out",
				ReportMsgf: reportInvalidArg,
			},
			{
				Fn:         "NotRegexp",
				Argsf:      "testStringer{}, out",
				ReportMsgf: reportInvalidArg,
			},
		},
		ValidAssertions: []Assertion{
			{Fn: "Regexp", Argsf: "`\\[.*\\] DEBUG \\(.*TestNew.*\\): message`, out"},
			{Fn: "NotRegexp", Argsf: "`\\[.*\\] TRACE message`, out"},
			{Fn: "Regexp", Argsf: "myStrAlias(`\\w+`), out"},
			{Fn: "NotRegexp", Argsf: "myStrAlias(`\\w+`), out"},
			{Fn: "Regexp", Argsf: "compiledRegexp, out"},
			{Fn: "NotRegexp", Argsf: "compiledRegexp, out"},
			{Fn: "Regexp", Argsf: "rxAlias, out"},
			{Fn: "NotRegexp", Argsf: "rxAlias, out"},
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

type myDefinedString string

type myStrAlias = string

type myRegexpStructAlias = regexp.Regexp

type testStringer struct{}

func (testStringer) String() string { return "42" }

// myRxAlias is a type alias for *regexp.Regexp used to verify that alias types are accepted.
type myRxAlias = *regexp.Regexp

// assertRegexpGeneric demonstrates that a type parameter is accepted conservatively.
// This differs from aliases such as myStrAlias or myRxAlias, which are handled separately.
func assertRegexpGeneric[T ~string](t *testing.T, rx T, str string) {
	assert.Regexp(t, rx, str)
	assert.NotRegexp(t, rx, str)
}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var out string
	compiledRegexp := regexp.MustCompile(` + "`" + `\w+` + "`" + `)
	var rxRegexpStructAlias *myRegexpStructAlias = regexp.MustCompile(` + "`" + `\w+` + "`" + `)
	var rxAlias myRxAlias = regexp.MustCompile(` + "`" + `\w+` + "`" + `)

	// Invalid: regexp.MustCompile usage.
	{
		{{- range $ai, $assrn := $.InvalidMustCompile }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Invalid: non-string, non-*regexp.Regexp first argument.
	{
		{{- range $ai, $assrn := $.InvalidTypeAssertions }}
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
