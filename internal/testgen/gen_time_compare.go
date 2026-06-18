package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type TimeCompareTestsGenerator struct{}

func (TimeCompareTestsGenerator) Checker() checkers.Checker {
	return checkers.NewTimeCompare()
}

func (g TimeCompareTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()

		report = checker + ": equality-based assertion on time.Time can be flaky"
	)

	return struct {
		CheckerName       CheckerName
		InvalidAssertions []Assertion
		ValidAssertions   []Assertion
	}{
		CheckerName: CheckerName(checker),
		InvalidAssertions: []Assertion{
			{Fn: "Equal", Argsf: "time.Now(), t1", ReportMsgf: report},
			{Fn: "Equal", Argsf: "alert.Updated, time.Now()", ReportMsgf: report},
			{Fn: "EqualValues", Argsf: "time.Now(), t1", ReportMsgf: report},
			{Fn: "EqualValues", Argsf: "t1, time.Now()", ReportMsgf: report},
			{Fn: "Exactly", Argsf: "time.Now(), cc.firstDetected", ReportMsgf: report},
			{Fn: "Exactly", Argsf: "t1, time.Now()", ReportMsgf: report},

			{Fn: "NotEqual", Argsf: "time.Now(), t1", ReportMsgf: report},
			{Fn: "NotEqual", Argsf: "t1, time.Now()", ReportMsgf: report},
			{Fn: "NotEqualValues", Argsf: "time.Now(), t1", ReportMsgf: report},
			{Fn: "NotEqualValues", Argsf: "t1, time.Now()", ReportMsgf: report},

			// Zero time equality.
			{Fn: "Equal", Argsf: "time.Time{}, t1", ReportMsgf: report},
			{Fn: "Equal", Argsf: "t1, time.Time{}", ReportMsgf: report},
			{Fn: "NotEqual", Argsf: "time.Time{}, t1", ReportMsgf: report},
			{Fn: "NotEqual", Argsf: "t1, time.Time{}", ReportMsgf: report},
		},
		ValidAssertions: []Assertion{
			{Fn: "True", Argsf: "t1.Equal(t2)"},
			{Fn: "False", Argsf: "t1.Equal(t2)"},

			{Fn: "Equal", Argsf: "t1.Unix(), t2.Unix()"},
			{Fn: "Equal", Argsf: "t1.UnixMilli(), t2.UnixMilli()"},
			{Fn: "EqualValues", Argsf: "t1.UnixMicro(), t2.UnixMicro()"},
			{Fn: "Exactly", Argsf: "t1.UnixNano(), t2.UnixNano()"},
			{Fn: "WithinDuration", Argsf: "t1, t2, 0"},

			// Suppressed by default because equality is normalized explicitly.
			{Fn: "Equal", Argsf: "t1.Truncate(0), t2.Round(0)"},
			{Fn: "Equal", Argsf: "t1.UTC(), t2.UTC()"},
			{Fn: "Equal", Argsf: "t1.Local(), t2.Local()"},
			{Fn: "Equal", Argsf: "t1.In(time.UTC), t2.In(time.UTC)"},
			{Fn: "Equal", Argsf: "t1.Add(time.Second), t2.Add(time.Second)"},
			{Fn: "Equal", Argsf: "zeroTime.Add(time.Second), t1"},
			{Fn: "Equal", Argsf: "t1.AddDate(0, 0, 1), t2.AddDate(0, 0, 1)"},
			{Fn: "Equal", Argsf: "time.Date(2017, time.January, 24, 0, 0, 0, 0, time.UTC), d"},

			// Zero time checks.
			{Fn: "True", Argsf: "t1.IsZero()"},
			{Fn: "False", Argsf: "now.IsZero()"},
			{Fn: "Zero", Argsf: "t2"},
			{Fn: "NotZero", Argsf: "now"},

			// Duration comparisons derived from time values.
			{Fn: "Equal", Argsf: "time.Duration(0), expectedTime.Round(time.Second).Sub(expiryTime.Round(time.Second))"},
			{Fn: "Equal", Argsf: "5 * time.Minute, time.Until(token.Expiry).Round(time.Second)"},

			// Compare results are ints; simplification is covered by the compares checker.
			{Fn: "Equal", Argsf: "0, t1.Compare(t2)"},
			{Fn: "True", Argsf: "t1.Compare(t2) == 0"},
			{Fn: "False", Argsf: "t1.Compare(t2) != 0"},

			{Fn: "EqualExportedValues", Argsf: "metric{Value: 100, t: t1}, metric{Value: 100, t: t2}"},
		},
	}
}

func (TimeCompareTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("TimeCompareTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(timeCompareTestTmpl))
}

func (TimeCompareTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const timeCompareTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var t1, t2, zeroTime time.Time
	d := time.Now()
	now := time.Now()
	expectedTime, expiryTime := time.Now(), time.Now()
	token := expiringToken{Expiry: time.Now()}
	alert := alert{Updated: time.Now()}
	cc := container{firstDetected: time.Now()}

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

type alert struct {
	Updated time.Time
}

type container struct {
	firstDetected time.Time
}

type expiringToken struct {
	Expiry time.Time
}

type metric struct {
	Value int
	t     time.Time
}
`
