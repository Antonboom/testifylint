package encodedcompareissue274

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression test: variables whose names contain a negative word (e.g. "invalid")
// followed by "JSON" should not be flagged by the encoded-compare checker,
// because asserting JSONEq on invalid JSON would cause the test to fail at runtime.
func TestInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{invalid json`)
	result := []byte(`{invalid json`)

	assert.Equal(t, string(invalidJSON), string(result))

	var badJSON, wrongJSON, malformedJSON, brokenJSON, corruptJSON string
	var other string

	assert.Equal(t, badJSON, other)
	assert.Equal(t, wrongJSON, other)
	assert.Equal(t, malformedJSON, other)
	assert.Equal(t, brokenJSON, other)
	assert.Equal(t, corruptJSON, other)
}
