package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ElementsMatchTestsGenerator struct{}

func (ElementsMatchTestsGenerator) Checker() checkers.Checker {
	return checkers.NewElementsMatch()
}

func (g ElementsMatchTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()

	return struct {
		CheckerName   CheckerName
		AssertReport  string
		RequireReport string
		AssertFReport string
	}{
		CheckerName:   CheckerName(checker),
		AssertReport:  QuoteReport(checker + ": use assert.ElementsMatch"),
		RequireReport: QuoteReport(checker + ": use require.ElementsMatch"),
		AssertFReport: QuoteReport(checker + ": use assert.ElementsMatchf"),
	}
}

func (ElementsMatchTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ElementsMatchTestsGenerator.ErroredTemplate").
		Parse(elementsMatchTestTmpl))
}

func (ElementsMatchTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ElementsMatchTestsGenerator.GoldenTemplate").
		Parse(elementsMatchGoldenTmpl))
}

const elementsMatchTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"sort"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var a, b []int

	// Invalid.
	{
		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, slices.Equal(a, b)) // want {{ $.AssertReport }}

		slices.Sort(b)
		slices.Sort(a)
		assert.True(t, slices.Equal(a, b)) // want {{ $.AssertReport }}

		slices.Sort(a)
		slices.Sort(b)
		require.True(t, slices.Equal(a, b)) // want {{ $.RequireReport }}

		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, slices.Equal(a, b), "elements should match") // want {{ $.AssertReport }}

		slices.Sort(a)
		slices.Sort(b)
		assert.Truef(t, slices.Equal(a, b), "elements should match %d and %d", len(a), len(b)) // want {{ $.AssertFReport }}

		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, slices.Equal(b, a)) // want {{ $.AssertReport }}

		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			require.Equal(t, a[i], b[i]) // want {{ $.RequireReport }}
		}

		assert.Equal(t, len(a), len(b))
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i]) // want {{ $.AssertReport }}
		}
	}

	// Valid.
	{
		assert.ElementsMatch(t, a, b)

		// Only one sort call preceding assert.
		slices.Sort(a)
		assert.True(t, slices.Equal(a, b))

		// Sort args don't match Equal args.
		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, slices.Equal(a, a))

		// Not a slices.Equal call.
		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, a[0] == b[0])

		// Not assert.True.
		slices.Sort(a)
		slices.Sort(b)
		assert.Equal(t, a, b)

		// Only one sort call preceding loop.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i])
		}

		// Sort args don't match loop args.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], a[i])
		}

		// Not element-wise loop comparison.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.True(t, a[i] == b[i])
		}

		// Loop body has extra statements.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i])
			assert.NotZero(t, i)
		}
	}
}
`

const elementsMatchGoldenTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"sort"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var a, b []int

	// Invalid.
	{
		assert.ElementsMatch(t, a, b) // want {{ $.AssertReport }}

		assert.ElementsMatch(t, a, b) // want {{ $.AssertReport }}

		require.ElementsMatch(t, a, b) // want {{ $.RequireReport }}

		assert.ElementsMatch(t, a, b, "elements should match") // want {{ $.AssertReport }}

		assert.ElementsMatchf(t, a, b, "elements should match %d and %d", len(a), len(b)) // want {{ $.AssertFReport }}

		assert.ElementsMatch(t, b, a) // want {{ $.AssertReport }}

		require.ElementsMatch(t, a, b)

		assert.Equal(t, len(a), len(b))
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i]) // want {{ $.AssertReport }}
		}
	}

	// Valid.
	{
		assert.ElementsMatch(t, a, b)

		// Only one sort call preceding assert.
		slices.Sort(a)
		assert.True(t, slices.Equal(a, b))

		// Sort args don't match Equal args.
		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, slices.Equal(a, a))

		// Not a slices.Equal call.
		slices.Sort(a)
		slices.Sort(b)
		assert.True(t, a[0] == b[0])

		// Not assert.True.
		slices.Sort(a)
		slices.Sort(b)
		assert.Equal(t, a, b)

		// Only one sort call preceding loop.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i])
		}

		// Sort args don't match loop args.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], a[i])
		}

		// Not element-wise loop comparison.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.True(t, a[i] == b[i])
		}

		// Loop body has extra statements.
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		for i := range b {
			assert.Equal(t, a[i], b[i])
			assert.NotZero(t, i)
		}
	}
}
`
