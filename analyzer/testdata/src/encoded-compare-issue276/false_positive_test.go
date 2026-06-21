package encodedcompareissue276

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Constants whose names contain "json"/"JSON" but whose values are not JSON
// should not be flagged by the `encoded-compare` checker.
func TestContentTypeHeader(t *testing.T) {
	const contentTypeJSON = "application/json"
	const contentTypeJson = "application/json"
	const acceptJSON = "application/json; charset=utf-8"

	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("Content-Type", contentTypeJSON)

	assert.Equal(t, contentTypeJSON, resp.Header.Get("Content-Type"))
	assert.Equal(t, contentTypeJson, resp.Header.Get("Content-Type"))
	assert.Equal(t, acceptJSON, resp.Header.Get("Accept"))
}
