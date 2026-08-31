package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ErrorFirstTestsGenerator struct{}

func (ErrorFirstTestsGenerator) Checker() checkers.Checker {
	return checkers.NewErrorFirst()
}

func (g ErrorFirstTestsGenerator) TemplateData() any {
	var (
		checker        = g.Checker().Name()
		reportNotFirst = checker + ": assert error before making other assertions"
	)

	return struct {
		CheckerName    CheckerName
		ReportNotFirst string
	}{
		CheckerName:    CheckerName(checker),
		ReportNotFirst: reportNotFirst,
	}
}

func (ErrorFirstTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ErrorFirstTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(errorFirstTestTmpl))
}

func (ErrorFirstTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const errorFirstTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resultAndErr() (int, error)                { return 0, nil }
func errAndResult() (error, int)                { return nil, 0 }
func resultOnly() int                           { return 0 }
func resultAndErrAndMore() (int, string, error) { return 0, "", nil }
func process(v int) *int                        { return &v }

var (
	cond  bool
	cond2 bool
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	// Invalid: error not checked before asserting result.
	{
		res, err := resultAndErr()
		assert.NotNil(t, res) // want "{{ .ReportNotFirst }}"
		_ = err
	}

	// Invalid: error in first position (not last), not checked.
	{
		err, res := errAndResult()
		assert.NotNil(t, res) // want "{{ .ReportNotFirst }}"
		_ = err
	}

	// Invalid: first call checked, second call not checked.
	{
		res, err := resultAndErr()
		require.NoError(t, err)
		res2, err2 := resultAndErr()
		assert.NotNil(t, res2) // want "{{ .ReportNotFirst }}"
		_ = res
		_ = err2
	}

	// Invalid: three return values (int, string, error); report only first result assertion.
	{
		res, msg, err := resultAndErrAndMore()
		assert.NotNil(t, res)  // want "{{ .ReportNotFirst }}"
		assert.NotEmpty(t, msg)
		_ = err
	}

	// Valid: error checked first with require.NoError.
	{
		res, err := resultAndErr()
		require.NoError(t, err)
		assert.NotNil(t, res)
	}

	// Valid: error checked first with assert.NoError.
	{
		res, err := resultAndErr()
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}

	// Valid: any assertion on err counts as checking the error.
	{
		res, err := resultAndErr()
		assert.Nil(t, err)
		assert.NotNil(t, res)
	}

	// Valid: non-canonical assertion on err (Equal) counts as checking the error.
	{
		res, err := resultAndErr()
		assert.Equal(t, nil, err)
		assert.NotNil(t, res)
	}

	// Invalid: err in message arg does not count as checking the error.
	{
		res, err := resultAndErr()
		assert.Equal(t, 0, res, err) // want "{{ .ReportNotFirst }}"
	}

	// Valid: error checked first with ErrorContains.
	{
		res, err := resultAndErr()
		require.ErrorContains(t, err, "msg")
		assert.Nil(t, res)
	}

	// Valid: error checked first with ErrorIs.
	{
		res, err := resultAndErr()
		require.ErrorIs(t, err, errors.New("msg"))
		assert.Nil(t, res)
	}

	// Valid: error assertion in if condition, result asserted in if body.
	{
		res, err := resultAndErr()
		if assert.NoError(t, err) {
			assert.NotNil(t, res)
		}
	}

	// Valid: no result used — only error asserted.
	{
		_, err := resultAndErr()
		require.NoError(t, err)
	}

	// Valid: function with no error return — not tracked.
	{
		res := resultOnly()
		assert.Equal(t, 0, res)
	}

	// Valid: transitive usage — process(res) returns no error, not tracked.
	{
		res, err := resultAndErr()
		require.NoError(t, err)
		res2 := process(res)
		assert.NotNil(t, res2)
	}

	// Valid: error in first position, checked first.
	{
		err, res := errAndResult()
		require.NoError(t, err)
		assert.NotNil(t, res)
	}

	// Valid: reassignment — first err checked, second err also checked.
	{
		res, err := resultAndErr()
		require.NoError(t, err)
		res, err = resultAndErr()
		require.NoError(t, err)
		assert.NotNil(t, res)
	}

	// Valid: reassignment from non-multi-return expression clears tracking.
	{
		res, err := resultAndErr()
		require.NoError(t, err)
		res = 1
		assert.Equal(t, 1, res)
	}

	// Valid: error checked in both branches of an exhaustive if/else, result used after.
	{
		res, err := resultAndErr()
		if cond {
			require.ErrorContains(t, err, "msg")
		} else {
			require.NoError(t, err)
		}
		assert.Equal(t, 0, res)
	}

	// Valid: error checked in every branch of an exhaustive if/else-if/else, result used after.
	{
		res, err := resultAndErr()
		if cond {
			require.ErrorContains(t, err, "msg")
		} else if cond2 {
			require.ErrorIs(t, err, errors.New("msg"))
		} else {
			require.NoError(t, err)
		}
		assert.Equal(t, 0, res)
	}

	// Invalid: if branch asserts error but there is no else — not exhaustive.
	{
		res, err := resultAndErr()
		if cond {
			require.NoError(t, err)
		}
		assert.Equal(t, 0, res) // want "{{ .ReportNotFirst }}"
	}

	// Valid: error checked in every case of an exhaustive switch (with default), result used after.
	{
		res, err := resultAndErr()
		switch {
		case cond:
			require.ErrorContains(t, err, "msg")
		case cond2:
			require.ErrorIs(t, err, errors.New("msg"))
		default:
			require.NoError(t, err)
		}
		assert.Equal(t, 0, res)
	}

	// Invalid: switch has no default case — not exhaustive.
	{
		res, err := resultAndErr()
		switch {
		case cond:
			require.NoError(t, err)
		}
		assert.Equal(t, 0, res) // want "{{ .ReportNotFirst }}"
	}
}
`
