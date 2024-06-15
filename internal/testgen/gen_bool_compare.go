package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type BoolCompareTestsGenerator struct{}

func (BoolCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewBoolCompare()
}

func (g BoolCompareTestsGenerator) TemplateData() any {
	var (
		checker        = g.Checker().Name()
		reportUse      = checker + ": use %s.%s"
		reportSimplify = checker + ": need to simplify the assertion"
	)

	type test struct {
		Name              string
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}

	var ignoredAssertions []Assertion
	for _, fn := range []string{
		"Equal",
		"NotEqual",
	} {
		ignoredAssertions = append(ignoredAssertions,
			Assertion{Fn: fn, Argsf: "true, true"},
			Assertion{Fn: fn, Argsf: "false, false"},
			Assertion{Fn: fn, Argsf: "predicate, predicate"},
			Assertion{Fn: fn, Argsf: "false, false"},

			// https://go.dev/ref/spec#Comparison_operators
			// A value x of non-interface type X and a value t of interface type T can be compared
			// if type X is comparable and X implements T.
			Assertion{Fn: fn, Argsf: `true, result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"], true`},
			Assertion{Fn: fn, Argsf: `false, result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"], false`},
		)
	}

	for _, fn := range []string{
		"True",
		"False",
	} {
		ignoredAssertions = append(ignoredAssertions,
			Assertion{Fn: fn, Argsf: "true == true"},
			Assertion{Fn: fn, Argsf: "false == false"},
			Assertion{Fn: fn, Argsf: "true != true"},
			Assertion{Fn: fn, Argsf: "false != false"},
			Assertion{Fn: fn, Argsf: "predicate == predicate"},
			Assertion{Fn: fn, Argsf: "predicate != predicate"},

			// https://go.dev/ref/spec#Comparison_operators
			// A value x of non-interface type X and a value t of interface type T can be compared
			// if type X is comparable and X implements T.

			Assertion{Fn: fn, Argsf: `true == result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"] == true`},
			Assertion{Fn: fn, Argsf: `false == result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"] == false`},
			Assertion{Fn: fn, Argsf: `true != result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"] != true`},
			Assertion{Fn: fn, Argsf: `false != result["flag"]`},
			Assertion{Fn: fn, Argsf: `result["flag"] != false`},
		)
	}

	// `any` cases.
	ignoredAssertions = append(ignoredAssertions,
		Assertion{Fn: "Equal", Argsf: "foo, foo"},
		Assertion{Fn: "NotEqual", Argsf: "foo, foo"},
		Assertion{Fn: "True", Argsf: "foo == foo"},
		Assertion{Fn: "False", Argsf: "foo == foo"},
		Assertion{Fn: "True", Argsf: "foo != foo"},
		Assertion{Fn: "False", Argsf: "foo != foo"},
	)

	var invalidAssertionsForTrue []Assertion
	var invalidAssertionsForFalse []Assertion
	for _, fn := range []string{
		"Equal",
		"EqualValues",
		"Exactly",
	} {
		invalidAssertionsForTrue = append(invalidAssertionsForTrue,
			Assertion{
				Fn: fn, Argsf: "predicate, true", ReportMsgf: reportUse,
				ProposedFn: "True", ProposedArgsf: "predicate",
			},
			Assertion{
				Fn: fn, Argsf: "true, predicate", ReportMsgf: reportUse,
				ProposedFn: "True", ProposedArgsf: "predicate",
			},
		)

		invalidAssertionsForFalse = append(invalidAssertionsForFalse,
			Assertion{
				Fn: fn, Argsf: "predicate, false",
				ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate",
			},
			Assertion{
				Fn: fn, Argsf: "false, predicate",
				ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate",
			},
		)
	}

	for _, fn := range []string{
		"NotEqual",
		"NotEqualValues",
	} {
		invalidAssertionsForTrue = append(invalidAssertionsForTrue,
			Assertion{
				Fn: fn, Argsf: "predicate, false",
				ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate",
			},
			Assertion{
				Fn: fn, Argsf: "false, predicate",
				ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate",
			},
		)

		invalidAssertionsForFalse = append(invalidAssertionsForFalse,
			Assertion{
				Fn: fn, Argsf: "predicate, true",
				ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate",
			},
			Assertion{
				Fn: fn, Argsf: "true, predicate",
				ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate",
			},
		)
	}

	for _, argsf := range []string{
		"predicate == false",
		"false == predicate",
		"predicate != true",
		"true != predicate",
		"!predicate",
	} {
		invalidAssertionsForTrue = append(invalidAssertionsForTrue,
			Assertion{
				Fn: "False", Argsf: argsf,
				ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: "predicate",
			},
		)

		invalidAssertionsForFalse = append(invalidAssertionsForFalse,
			Assertion{
				Fn: "True", Argsf: argsf,
				ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: "predicate",
			},
		)
	}

	for _, argsf := range []string{
		"predicate == true",
		"true == predicate",
		"predicate != false",
		"false != predicate",
	} {
		invalidAssertionsForTrue = append(invalidAssertionsForTrue,
			Assertion{Fn: "True", Argsf: argsf, ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
		)

		invalidAssertionsForFalse = append(invalidAssertionsForFalse,
			Assertion{Fn: "False", Argsf: argsf, ReportMsgf: reportSimplify, ProposedArgsf: "predicate"},
		)
	}

	invalidAssertionsForTrue = append(invalidAssertionsForTrue,
		Assertion{
			Fn: "False", Argsf: `!result["flag"].(bool)`,
			ReportMsgf: reportUse, ProposedFn: "True", ProposedArgsf: `result["flag"].(bool)`,
		},
	)

	invalidAssertionsForFalse = append(invalidAssertionsForFalse,
		Assertion{
			Fn: "True", Argsf: `!result["flag"].(bool)`,
			ReportMsgf: reportUse, ProposedFn: "False", ProposedArgsf: `result["flag"].(bool)`,
		},
	)

	return struct {
		CheckerName       CheckerName
		Tests             []test
		IgnoredAssertions []Assertion
	}{
		CheckerName: CheckerName(checker),
		Tests: []test{
			{
				Name:              "assert.True cases",
				InvalidAssertions: invalidAssertionsForTrue,
				ValidAssertions: []Assertion{
					{Fn: "True", Argsf: "predicate"},
				},
			},
			{
				Name:              "assert.False cases",
				InvalidAssertions: invalidAssertionsForFalse,
				ValidAssertions: []Assertion{
					{Fn: "False", Argsf: "predicate"},
				},
			},
		},
		IgnoredAssertions: ignoredAssertions,
	}
}

func (BoolCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("BoolCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(boolCompareTestTmpl))
}

func (BoolCompareTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("BoolCompareTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(boolCompareTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const boolCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var predicate bool
	result := map[string]any{}

	{{ range $ti, $test := $.Tests }}
		// {{ $test.Name }}.
		{
			// Invalid.
			{{- range $ai, $assrn := $test.InvalidAssertions }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}

			// Valid.
			{{- range $ai, $assrn := $test.ValidAssertions }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
			{{- end }}
		}
	{{ end -}}
}

func {{ .CheckerName.AsTestName }}_Ignored(t *testing.T) {
	var predicate bool
	var foo any
	result := map[string]any{}

	foo = true
	assert.Equal(t, true, foo)

	{{ range $ai, $assrn := $.IgnoredAssertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
