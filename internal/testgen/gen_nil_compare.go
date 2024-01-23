package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type NilCompareTestsGenerator struct{}

func (NilCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewNilCompare()
}

func (g NilCompareTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": use %s.%s"
	)

	var unsupportedAssertions []Assertion
	for _, fn := range []string{"Equal", "EqualValues", "Exactly", "NotEqual", "NotEqualValues"} {
		unsupportedAssertions = append(unsupportedAssertions, []Assertion{
			{Fn: fn, Argsf: "(chan struct{})(nil), ch"},
			{Fn: fn, Argsf: "(func())(nil), fn"},
			{Fn: fn, Argsf: "(any)(nil), iface"},
			{Fn: fn, Argsf: "(map[int]int)(nil), mp"},
			{Fn: fn, Argsf: "(*int)(nil), ptr"},
			{Fn: fn, Argsf: "[]int(nil), slice"},
			{Fn: fn, Argsf: "(unsafe.Pointer)(nil), unsafePtr"},
		}...)
	}

	return struct {
		CheckerName           CheckerName
		BullshitAssertions    []Assertion
		InvalidAssertions     []Assertion
		UnsupportedAssertions []Assertion
		ValidAssertions       []Assertion
		IgnoredAssertions     []Assertion
		UntypedNilBug         []Assertion
	}{
		CheckerName: CheckerName(checker),
		BullshitAssertions: []Assertion{
			{Fn: "Equal", Argsf: "nil, %s", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},
			{Fn: "Equal", Argsf: "%s, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},
			{Fn: "EqualValues", Argsf: "nil, %s", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},
			{Fn: "EqualValues", Argsf: "%s, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},
			{Fn: "Exactly", Argsf: "nil, %s", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},
			{Fn: "Exactly", Argsf: "%s, nil", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "%s"},

			{Fn: "NotEqual", Argsf: "nil, %s", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "%s"},
			{Fn: "NotEqual", Argsf: "%s, nil", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "%s"},
			{Fn: "NotEqualValues", Argsf: "nil, %s", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "%s"},
			{Fn: "NotEqualValues", Argsf: "%s, nil", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "%s"},
		},
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "nil, iface", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "iface"},
			{Fn: "EqualValues", Argsf: "nil, iface", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "iface"},
			{Fn: "Exactly", Argsf: "nil, iface", ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: "iface"},
			{Fn: "NotEqual", Argsf: "nil, iface", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "iface"},
			{Fn: "NotEqualValues", Argsf: "nil, iface", ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: "iface"},
		},
		UnsupportedAssertions: unsupportedAssertions,
		ValidAssertions: []Assertion{
			{Fn: "Nil", Argsf: "iface"},
			{Fn: "NotNil", Argsf: "iface"},
		},
		IgnoredAssertions: []Assertion{
			{Fn: "Equal", Argsf: "iface, iface"},
			{Fn: "Equal", Argsf: "nil, nil"},

			{Fn: "EqualValues", Argsf: "iface, iface"},
			{Fn: "EqualValues", Argsf: "nil, nil"},

			{Fn: "Exactly", Argsf: "iface, iface"},
			{Fn: "Exactly", Argsf: "nil, nil"},

			{Fn: "NotEqual", Argsf: "iface, iface"},
			{Fn: "NotEqual", Argsf: "nil, nil"},

			{Fn: "NotEqualValues", Argsf: "iface, iface"},
			{Fn: "NotEqualValues", Argsf: "nil, nil"},
		},
		// Some special cases that covers bug from the past. The reason is
		// pass.TypesInfo.TypeOf(interface).(*types.Basic).Kind() == types.UntypedNil
		// https://github.com/golang/go/blob/3e54329cbe41e887f309c733a619206bf6351768/src/go/types/expr.go#L428
		UntypedNilBug: []Assertion{
			{Fn: "Equal", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "EqualValues", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "Exactly", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "Nil", ProposedArgsf: `Row["col"]`},
			{Fn: "NotEqual", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},
			{Fn: "NotEqualValues", Argsf: `nil, Row["col"]`, ReportMsgf: report, ProposedFn: "NotNil", ProposedArgsf: `Row["col"]`},

			{Fn: "Equal", Argsf: `Row["col"], "foo"`},
			{Fn: "Equal", Argsf: `"foo", Row["col"]`},
			{Fn: "Equal", Argsf: `Row["col"], Row["col"]`},

			{Fn: "EqualValues", Argsf: `Row["col"], "foo"`},
			{Fn: "EqualValues", Argsf: `"foo", Row["col"]`},
			{Fn: "EqualValues", Argsf: `Row["col"], Row["col"]`},

			{Fn: "Exactly", Argsf: `Row["col"], "foo"`},
			{Fn: "Exactly", Argsf: `"foo", Row["col"]`},
			{Fn: "Exactly", Argsf: `Row["col"], Row["col"]`},

			{Fn: "NotEqual", Argsf: `Row["col"], "foo"`},
			{Fn: "NotEqual", Argsf: `"foo", Row["col"]`},
			{Fn: "NotEqual", Argsf: `Row["col"], Row["col"]`},

			{Fn: "NotEqualValues", Argsf: `Row["col"], "foo"`},
			{Fn: "NotEqualValues", Argsf: `"foo", Row["col"]`},
			{Fn: "NotEqualValues", Argsf: `Row["col"], Row["col"]`},
		},
	}
}

func (NilCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("NilCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(nilCompareTestTmpl))
}

func (NilCompareTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("NilCompareTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(nilCompareTestTmpl, "NewAssertionExpander", "NewAssertionExpander.AsGolden")))
}

const nilCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var (
		ch        chan struct{}
		fn        func()
		iface     any
		mp        map[int]int
		ptr       *int
		slice     []int
		unsafePtr unsafe.Pointer
	)

	// Not working as expected!
	{
		{{- range $ai, $assrn := $.BullshitAssertions }}
			{{- range $vi, $var := (arr "ch" "fn" "mp" "ptr" "slice" "unsafePtr") }}
				{{ NewAssertionExpander.Expand $assrn "assert" "t" (arr $var) }}
			{{- end }}
		{{- end }}
	}

	// Invalid.
	{
		{{- range $ai, $assrn := $.InvalidAssertions }}
			{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
		{{- end }}
	}

	// Invalid, but unsupported.
	{
		{{- range $ai, $assrn := $.UnsupportedAssertions }}
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

func {{ .CheckerName.AsTestName }}_UntypedNilBug(t *testing.T) {
	var Row map[string]any
	{{ range $ai, $assrn := $.UntypedNilBug }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
}
`
