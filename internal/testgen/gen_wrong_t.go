package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type WrongTTestsGenerator struct{}

func (WrongTTestsGenerator) Checker() checkers.Checker {
	return checkers.NewWrongT()
}

func (g WrongTTestsGenerator) TemplateData() any {
	var (
		checker     = g.Checker().Name()
		nilReport   = QuoteReport(checker + ": do not pass nil as testing.T")
		freshReport = QuoteReport(checker + ": do not pass freshly created testing.T")
	)

	return struct {
		CheckerName CheckerName
		NilReport   string
		FreshReport string
		Checker     string
	}{
		CheckerName: CheckerName(checker),
		NilReport:   nilReport,
		FreshReport: freshReport,
		Checker:     checker,
	}
}

func (WrongTTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("WrongTTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(wrongTTestTmpl))
}

func (WrongTTestsGenerator) GoldenTemplate() Executor {
	// No autofix available for these issues.
	return nil
}

const wrongTTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Case 1: nil testing.T passed to package-level assertions.
func {{ .CheckerName.AsTestName }}_NilT(t *testing.T) {
	var err error
	var result int

	assert.Equal(nil, 0, result)         // want {{ $.NilReport }}
	assert.NoError(nil, err)             // want {{ $.NilReport }}
	require.Equal(nil, 0, result)        // want {{ $.NilReport }}
	require.NoError(nil, err)            // want {{ $.NilReport }}

	// Valid – using the test's t.
	assert.Equal(t, 0, result)
	assert.NoError(t, err)
	require.Equal(t, 0, result)
	require.NoError(t, err)

	t.Run("sub", func(t *testing.T) {
		assert.Equal(t, 0, result)
		require.NoError(t, err)
	})
}

// Case 2: freshly created testing.T passed to package-level assertions.
func {{ .CheckerName.AsTestName }}_FreshT(t *testing.T) {
	var err error
	var result int

	assert.Equal(&testing.T{}, 0, result)    // want {{ $.FreshReport }}
	assert.NoError(new(testing.T), err)       // want {{ $.FreshReport }}
	require.Equal(&testing.T{}, 0, result)   // want {{ $.FreshReport }}
	require.NoError(new(testing.T), err)      // want {{ $.FreshReport }}

	u := &testing.T{}
	assert.Equal(u, 0, result) // want {{ $.FreshReport }}
	require.Equal(u, 0, result) // want {{ $.FreshReport }}

	v := new(testing.T)
	assert.NoError(v, err) // want {{ $.FreshReport }}
	require.NoError(v, err) // want {{ $.FreshReport }}

	// Valid – fresh variable reassigned to a real t before use.
	w := &testing.T{}
	w = t
	assert.Equal(w, 0, result)
	require.NoError(w, err)
}

// Case 3: assertion object created with outer t used inside t.Run subtest.
func {{ .CheckerName.AsTestName }}_Scope(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	var err error
	var result int

	// Valid – used in same scope.
	a.Equal(0, result)
	r.NoError(err)

	cases := []struct{ Name string }{}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			a.Equal(0, result)  // want {{ QuoteReport (printf "%s: assertion object a was created with outer scope's testing.T" $.Checker) }}
			r.NoError(err)      // want {{ QuoteReport (printf "%s: assertion object r was created with outer scope's testing.T" $.Checker) }}

			// Valid – created inside this subtest.
			subA := assert.New(t)
			subA.Equal(0, result)
		})
	}

	t.Run("single", func(t *testing.T) {
		a.Equal(0, result) // want {{ QuoteReport (printf "%s: assertion object a was created with outer scope's testing.T" $.Checker) }}

		// Nested subtest – a is still from outer scope.
		t.Run("nested", func(t *testing.T) {
			a.Equal(0, result) // want {{ QuoteReport (printf "%s: assertion object a was created with outer scope's testing.T" $.Checker) }}

			// Valid – b is created inside this nested subtest.
			b := assert.New(t)
			b.Equal(0, result)
		})
	})

	// Rebound to outer t (not the callback's t): still an error because the
	// rebinding does not use the subtest's own t.
	outerT := t
	t.Run("rebind-wrong-t", func(t *testing.T) {
		a = assert.New(outerT)
		a.Equal(0, result) // want {{ QuoteReport (printf "%s: assertion object a was created with outer scope's testing.T" $.Checker) }}
	})
}

// Valid uses: object-style calls (no t argument needed), or t from same scope.
func {{ .CheckerName.AsTestName }}_Valid(t *testing.T) {
	var err error
	var result int

	a := assert.New(t)
	r := require.New(t)

	a.Equal(0, result)
	r.NoError(err)

	t.Run("sub", func(t *testing.T) {
		assert.Equal(t, 0, result)
		require.NoError(t, err)

		inner := assert.New(t)
		inner.Equal(0, result)
	})

	// Valid – outer-scope assertion object rebound to subtest's t before use.
	t.Run("rebound", func(t *testing.T) {
		a = assert.New(t)
		r = require.New(t)
		a.Equal(0, result)
		r.NoError(err)
	})
}
`
