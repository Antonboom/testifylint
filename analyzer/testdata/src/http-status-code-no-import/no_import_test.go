package httpstatuscodenoimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHTTPStatusCodeNoImport verifies that the autofix adds the "net/http" import
// when it is absent from the file.
func TestHTTPStatusCodeNoImport(t *testing.T) {
	assert.HTTPStatusCode(t, handleOK, "GET", "/", nil, 200) // want "http-status-code: use net/http constant instead of integer literal"
}
