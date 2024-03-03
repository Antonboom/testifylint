package gorequireissue66

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE(a.telyshev): Neither assert nor require is the best way to test an HTTP handler:
// it leads to redundant stack traces, as well as EOF from the HTTP client.
// Use HTTP mechanisms and place assertions in the main test function.

var mockHTTPFromFile = func(t *testing.T) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open("some file.json")
		assert.NoError(t, err) // want "require-error: for error assertions use require"

		data, err := io.ReadAll(file)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		assert.NoError(t, err)
	}
}

func TestGoRequireAndRequireErrorConflict(t *testing.T) {
	ts := httptest.NewServer(mockHTTPFromFile(t))
	defer ts.Close()

	client := ts.Client()

	req, err := http.NewRequest("GET", ts.URL+"/example", nil)
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
