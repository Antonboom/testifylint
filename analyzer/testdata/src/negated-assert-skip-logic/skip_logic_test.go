package negatedassertskiplogic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompoundOr verifies that compound `||` patterns are detected and fixed.
func TestCompoundOr(t *testing.T) {
	var (
		err error
		a   int
	)

	// Single negated assertion — reported.
	if !assert.NoError(t, err) { // want "negated-assert: use require instead"
		return
	}

	// Compound || — reported (first assertion gets the diagnostic).
	if !assert.NoError(t, err) || !assert.Equal(t, a, a) { // want "negated-assert: use require instead"
		return
	}

	// Three-way || — reported.
	if !assert.NoError(t, err) || !assert.Equal(t, a, a) || !assert.True(t, true) { // want "negated-assert: use require instead"
		return
	}

	// Mixed: non-assert condition in || — not fixable, skipped.
	cond := true
	if !assert.NoError(t, err) || cond {
		return
	}

	// && semantics are different — skipped.
	if !assert.NoError(t, err) && !assert.Equal(t, a, a) {
		return
	}

	// Already uses require — not reported.
	require.NoError(t, err)
}

// TestGoroutine verifies that patterns inside goroutines are skipped.
func TestGoroutine(t *testing.T) {
	var err error

	go func() {
		// Skipped: inside a goroutine.
		if !assert.NoError(t, err) {
			return
		}
	}()
}

// TestHTTPHandler verifies that patterns inside HTTP handlers are skipped.
func TestHTTPHandler(t *testing.T) {
	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct{ ID string }
		err := json.NewDecoder(r.Body).Decode(&req)
		// Skipped: inside an HTTP handler.
		if !assert.NoError(t, err) {
			return
		}
	}))
}

// TestCleanup verifies that patterns inside t.Cleanup are skipped.
func TestCleanup(t *testing.T) {
	var err error

	t.Cleanup(func() {
		// Skipped: inside a test cleanup.
		if !assert.NoError(t, err) {
			return
		}
	})
}

// TestContinue verifies that `if !assert.Xxx { continue }` patterns in loops are fixed.
func TestContinue(t *testing.T) {
	var err error

	for i := 0; i < 10; i++ {
		if !assert.NoError(t, err) { // want "negated-assert: use require instead"
			continue
		}
	}
}
