package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type HTTPMultipleTestsGenerator struct{}

func (HTTPMultipleTestsGenerator) Checker() checkers.Checker {
	return checkers.NewHTTPMultiple()
}

func (g HTTPMultipleTestsGenerator) TemplateData() any {
	checker := g.Checker().Name()
	report := checker + ": use httptest.NewRecorder() instead of multiple HTTP assertions for the same handler call"

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(checker),
		Report:      report,
	}
}

func (HTTPMultipleTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("HTTPMultipleTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(httpMultipleTestTmpl))
}

func (HTTPMultipleTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("HTTPMultipleTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(httpMultipleGoldenTmpl))
}

const httpMultipleTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	// Invalid: multiple HTTP assertions with the same handler and args.
	assert.HTTPStatusCode(t, handler, "GET", "/index", nil, http.StatusOK)    // want {{ QuoteReport .Report }}
	assert.HTTPBodyContains(t, handler, "GET", "/index", nil, "hello")        // want {{ QuoteReport .Report }}

	assert.HTTPError(t, handler, "GET", "/error", nil)               // want {{ QuoteReport .Report }}
	assert.HTTPBodyContains(t, handler, "GET", "/error", nil, "oops") // want {{ QuoteReport .Report }}

	require.HTTPRedirect(t, handler, "GET", "/redirect", nil)                        // want {{ QuoteReport .Report }}
	require.HTTPStatusCode(t, handler, "GET", "/redirect", nil, http.StatusFound)    // want {{ QuoteReport .Report }}

	// Invalid: formatted (*f) variants are also detected.
	assert.HTTPStatusCodef(t, handler, "GET", "/fmt", nil, http.StatusOK, "msg")   // want {{ QuoteReport .Report }}
	assert.HTTPBodyContainsf(t, handler, "GET", "/fmt", nil, "hello", "msg")       // want {{ QuoteReport .Report }}

	// Invalid: non-nil url.Values query parameters.
	assert.HTTPStatusCode(t, handler, "GET", "/search", url.Values{"q": {"go"}}, http.StatusOK) // want {{ QuoteReport .Report }}
	assert.HTTPBodyContains(t, handler, "GET", "/search", url.Values{"q": {"go"}}, "result")    // want {{ QuoteReport .Report }}

	// Valid: single HTTP assertion.
	assert.HTTPStatusCode(t, handler, "GET", "/single", nil, http.StatusOK)

	// Valid: different handlers.
	assert.HTTPSuccess(t, handler, "GET", "/a", nil)
	assert.HTTPSuccess(t, http.NotFound, "GET", "/a", nil)

	// Valid: different methods.
	assert.HTTPSuccess(t, handler, "GET", "/b", nil)
	assert.HTTPSuccess(t, handler, "POST", "/b", nil)

	// Valid: different URLs.
	assert.HTTPSuccess(t, handler, "GET", "/c", nil)
	assert.HTTPSuccess(t, handler, "GET", "/d", nil)

	// Valid: same-key assertions in separate closures are independent scopes.
	t.Run("sub1", func(t *testing.T) {
		assert.HTTPStatusCode(t, handler, "GET", "/subtest", nil, http.StatusOK)
	})
	t.Run("sub2", func(t *testing.T) {
		assert.HTTPBodyContains(t, handler, "GET", "/subtest", nil, "hello")
	})

	// Valid: goroutine closure is an independent scope.
	go func() {
		assert.HTTPStatusCode(t, handler, "GET", "/goroutine", nil, http.StatusOK)
	}()
	go func() {
		assert.HTTPBodyContains(t, handler, "GET", "/goroutine", nil, "hello")
	}()

	// Valid: using httptest directly (the recommended approach).
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/direct", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// {{ .CheckerName.AsTestName }}TT uses a non-"t" TestingT variable to verify the fix
// extracts the actual identifier from the call rather than hard-coding "t".
func {{ .CheckerName.AsTestName }}TT(tt *testing.T) {
	assert.HTTPStatusCode(tt, handler, "GET", "/tt-path", nil, http.StatusOK) // want {{ QuoteReport .Report }}
	assert.HTTPBodyContains(tt, handler, "GET", "/tt-path", nil, "hello")     // want {{ QuoteReport .Report }}
}
`

// httpMultipleGoldenTmpl is the expected state of the test file after all suggested fixes
// have been applied: each group of HTTP assertions is replaced by a scoped httptest block.
const httpMultipleGoldenTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	// Invalid: multiple HTTP assertions with the same handler and args.
	{
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/index", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "hello")
	}

	{
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/error", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.GreaterOrEqual(t, rr.Code, http.StatusBadRequest)
		assert.Contains(t, rr.Body.String(), "oops")
	}

	{
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/redirect", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		require.GreaterOrEqual(t, rr.Code, http.StatusMultipleChoices)
		require.Less(t, rr.Code, http.StatusBadRequest)
		require.Equal(t, http.StatusFound, rr.Code)
	}

	// Invalid: formatted (*f) variants are also detected.
	{
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/fmt", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equalf(t, http.StatusOK, rr.Code, "msg")
		assert.Containsf(t, rr.Body.String(), "hello", "msg")
	}

	// Invalid: non-nil url.Values query parameters.
	{
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/search", nil)
		req.URL.RawQuery = url.Values{"q": {"go"}}.Encode()
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "result")
	}

	// Valid: single HTTP assertion.
	assert.HTTPStatusCode(t, handler, "GET", "/single", nil, http.StatusOK)

	// Valid: different handlers.
	assert.HTTPSuccess(t, handler, "GET", "/a", nil)
	assert.HTTPSuccess(t, http.NotFound, "GET", "/a", nil)

	// Valid: different methods.
	assert.HTTPSuccess(t, handler, "GET", "/b", nil)
	assert.HTTPSuccess(t, handler, "POST", "/b", nil)

	// Valid: different URLs.
	assert.HTTPSuccess(t, handler, "GET", "/c", nil)
	assert.HTTPSuccess(t, handler, "GET", "/d", nil)

	// Valid: same-key assertions in separate closures are independent scopes.
	t.Run("sub1", func(t *testing.T) {
		assert.HTTPStatusCode(t, handler, "GET", "/subtest", nil, http.StatusOK)
	})
	t.Run("sub2", func(t *testing.T) {
		assert.HTTPBodyContains(t, handler, "GET", "/subtest", nil, "hello")
	})

	// Valid: goroutine closure is an independent scope.
	go func() {
		assert.HTTPStatusCode(t, handler, "GET", "/goroutine", nil, http.StatusOK)
	}()
	go func() {
		assert.HTTPBodyContains(t, handler, "GET", "/goroutine", nil, "hello")
	}()

	// Valid: using httptest directly (the recommended approach).
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/direct", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// {{ .CheckerName.AsTestName }}TT uses a non-"t" TestingT variable to verify the fix
// extracts the actual identifier from the call rather than hard-coding "t".
func {{ .CheckerName.AsTestName }}TT(tt *testing.T) {
	{
		req := httptest.NewRequestWithContext(tt.Context(), http.MethodGet, "/tt-path", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(tt, http.StatusOK, rr.Code)
		assert.Contains(tt, rr.Body.String(), "hello")
	}
}
`
