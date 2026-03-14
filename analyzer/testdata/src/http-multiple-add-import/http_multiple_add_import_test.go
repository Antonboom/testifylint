// Test case for the http-multiple checker's import-addition behaviour.
// This file intentionally does NOT import net/http/httptest;
// the suggested fix must add the import automatically via addImportFix.

package httpmultipleaddimport

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func TestHttpMultipleCheckerAddImport(t *testing.T) {
	// Invalid: httptest is not imported — the fix must add it.
	assert.HTTPStatusCode(t, handler, "GET", "/add-import", nil, http.StatusOK) // want "http-multiple: use httptest\\.NewRecorder\\(\\) instead of multiple HTTP assertions for the same handler call"
	assert.HTTPBodyContains(t, handler, "GET", "/add-import", nil, "body")      // want "http-multiple: use httptest\\.NewRecorder\\(\\) instead of multiple HTTP assertions for the same handler call"
}
