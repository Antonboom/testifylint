package encodedcompareissue276

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression test for https://github.com/Antonboom/testifylint/issues/276
// Constants whose names contain "json"/"JSON" but whose values are not JSON
// should not be flagged by the encoded-compare checker.
func TestContentTypeHeader(t *testing.T) {
	const contentTypeJSON = "application/json"
	const acceptJSON = "application/json; charset=utf-8"

	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("Content-Type", contentTypeJSON)

	assert.Equal(t, contentTypeJSON, resp.Header.Get("Content-Type"))
	assert.Equal(t, acceptJSON, resp.Header.Get("Accept"))
}
