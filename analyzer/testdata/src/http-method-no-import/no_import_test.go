package httpmethodnoimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHTTPMethodNoImport verifies that the autofix adds the "net/http" import
// when it is absent from the file.
func TestHTTPMethodNoImport(t *testing.T) {
	assert.HTTPStatusCode(t, handleOK, "GET", "/", nil, 200) // want "http-method: use net/http constant instead of string literal"
}
