package debug

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	assert.False(t, true)
	assert.False(t, true, new(http.Response))
	// assert.False(t, true, new(http.Response), 1, 2, 3) // panic
	assert.False(t, true, "hello")
	assert.False(t, true, "hello", 1, 2)
	assert.False(t, true, "hello_%v_%d", 3, 4)
	assert.Falsef(t, true, "world")
	assert.Falsef(t, true, "world_%d_%v", 5, 6)
}
