package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type ManualAssertTestsGenerator struct{}

func (ManualAssertTestsGenerator) Checker() checkers.Checker {
	return checkers.NewManualAssert()
}

func (g ManualAssertTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()
	return struct {
		CheckerName CheckerName
		Report      string // sprintf with (pkg, fn)
	}{
		CheckerName: CheckerName(checker),
		Report:      checker + `: replace with %s\\.%s`,
	}
}

func (ManualAssertTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("ManualAssertTestsGenerator.ErroredTemplate").
		Parse(manualAssertTestTmpl))
}

func (ManualAssertTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("ManualAssertTestsGenerator.GoldenTemplate").
		Parse(manualAssertGoldenTmpl))
}

const manualAssertTestHeader = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
`

const manualAssertTestTmpl = manualAssertTestHeader + `
func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error
	var s string
	var p *int
	var b bool
	var sl []int
	var m map[string]int
	target := errors.New("target")

	// Invalid: fatal -> require.
	{
		if err != nil { // want "{{ printf $.Report "require" "NoError" }}"
			t.Fatalf("ValidateToken() error = %v", err)
		}
		if err == nil { // want "{{ printf $.Report "require" "Error" }}"
			t.Fatal("expected error")
		}
		if s == "" { // want "{{ printf $.Report "require" "NotEmpty" }}"
			t.Fatal("GenerateToken() returned empty token")
		}
		if s != "" { // want "{{ printf $.Report "require" "Empty" }}"
			t.Fatal("expected empty")
		}
		if p == nil { // want "{{ printf $.Report "require" "NotNil" }}"
			t.Fatal("nil pointer")
		}
		if p != nil { // want "{{ printf $.Report "require" "Nil" }}"
			t.Fatal("expected nil pointer")
		}
		if b { // want "{{ printf $.Report "require" "False" }}"
			t.Fatal("expected false")
		}
		if !b { // want "{{ printf $.Report "require" "True" }}"
			t.Fatal("expected true")
		}
		if len(sl) == 0 { // want "{{ printf $.Report "require" "NotEmpty" }}"
			t.Fatal("empty slice")
		}
		if len(sl) != 0 { // want "{{ printf $.Report "require" "Empty" }}"
			t.Fatal("expected empty slice")
		}
		if len(sl) != 3 { // want "{{ printf $.Report "require" "Len" }}"
			t.Fatal("wrong len")
		}
		if !errors.Is(err, target) { // want "{{ printf $.Report "require" "ErrorIs" }}"
			t.Fatal("expected target")
		}
		if errors.Is(err, target) { // want "{{ printf $.Report "require" "NotErrorIs" }}"
			t.Fatal("did not expect target")
		}
		if !strings.Contains(s, "sub") { // want "{{ printf $.Report "require" "Contains" }}"
			t.Fatal("missing substring")
		}
		if strings.Contains(s, "sub") { // want "{{ printf $.Report "require" "NotContains" }}"
			t.Fatal("unexpected substring")
		}
	}

	// Invalid: error -> assert.
	{
		var got, want int
		if got != want { // want "{{ printf $.Report "assert" "Equal" }}"
			t.Errorf("got = %v, want %v", got, want)
		}
		if got == want { // want "{{ printf $.Report "assert" "NotEqual" }}"
			t.Errorf("got = %v, want != %v", got, want)
		}
		if err != nil { // want "{{ printf $.Report "assert" "NoError" }}"
			t.Errorf("err = %v", err)
		}
	}

	// Valid: already using testify.
	{
		require.NoError(t, err)
		assert.Equal(t, 1, 2)
		require.NotEmpty(t, s)
	}

	// Ignored: if/else.
	{
		if err != nil {
			t.Fatal("bad")
		} else {
			t.Log("ok")
		}
	}

	// Ignored: multi-statement body.
	{
		if err != nil {
			t.Log("oh no")
			t.Fatal("bad")
		}
	}

	// Ignored: non-Fatal/Error sink.
	{
		if err != nil {
			t.Log("just logging")
		}
	}

	_ = m
}
`

// The golden file is what the test file looks like after the suggested fixes
// are applied. Note: ` + "`strings`" + ` becomes unused after the Contains/NotContains
// rewrites, and analysistest goimports-cleans it from the result.
const manualAssertGoldenTmplBody = `
package {{ .CheckerName.AsPkgName }}

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error
	var s string
	var p *int
	var b bool
	var sl []int
	var m map[string]int
	target := errors.New("target")

	// Invalid: fatal -> require.
	{
		require.NoError(t, err, "ValidateToken() error = %v", err)
		require.Error(t, err, "expected error")
		require.NotEmpty(t, s, "GenerateToken() returned empty token")
		require.Empty(t, s, "expected empty")
		require.NotNil(t, p, "nil pointer")
		require.Nil(t, p, "expected nil pointer")
		require.False(t, b, "expected false")
		require.True(t, b, "expected true")
		require.NotEmpty(t, sl, "empty slice")
		require.Empty(t, sl, "expected empty slice")
		require.Len(t, sl, 3, "wrong len")
		require.ErrorIs(t, err, target, "expected target")
		require.NotErrorIs(t, err, target, "did not expect target")
		require.Contains(t, s, "sub", "missing substring")
		require.NotContains(t, s, "sub", "unexpected substring")
	}

	// Invalid: error -> assert.
	{
		var got, want int
		assert.Equal(t, want, got, "got = %v, want %v", got, want)
		assert.NotEqual(t, want, got, "got = %v, want != %v", got, want)
		assert.NoError(t, err, "err = %v", err)
	}

	// Valid: already using testify.
	{
		require.NoError(t, err)
		assert.Equal(t, 1, 2)
		require.NotEmpty(t, s)
	}

	// Ignored: if/else.
	{
		if err != nil {
			t.Fatal("bad")
		} else {
			t.Log("ok")
		}
	}

	// Ignored: multi-statement body.
	{
		if err != nil {
			t.Log("oh no")
			t.Fatal("bad")
		}
	}

	// Ignored: non-Fatal/Error sink.
	{
		if err != nil {
			t.Log("just logging")
		}
	}

	_ = m
}
`

const manualAssertGoldenTmpl = header + manualAssertGoldenTmplBody
