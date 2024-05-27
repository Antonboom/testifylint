package debug

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameNotSameReplacements(t *testing.T) {
	var tm tMock

	a := new(os.PathError)
	b := new(os.PathError)
	c := &os.PathError{Path: "/tmp"}

	// Invalid.
	assert.Equal(t, assert.True(tm, a == b), assert.Equal(tm, a, b))
	assert.Equal(t, assert.True(tm, a != b), assert.NotEqual(tm, a, b))
	assert.Equal(t, assert.True(tm, a == c), assert.Equal(tm, a, c))
	assert.Equal(t, assert.True(tm, a != c), assert.NotEqual(tm, a, c))

	// Valid.
	assert.Equal(t, assert.True(tm, a == b), assert.Same(tm, a, b))
	assert.Equal(t, assert.True(tm, a != b), assert.NotSame(tm, a, b))
	assert.Equal(t, assert.True(tm, a == c), assert.Same(tm, a, c))
	assert.Equal(t, assert.True(tm, a != c), assert.NotSame(tm, a, c))
}
