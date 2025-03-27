package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUselessAsserts(t *testing.T) {
	assert.Empty(t, "")
	assert.False(t, false)
	assert.Implements(t, (*any)(nil), new(testing.T))
	assert.Negative(t, -42)
	assert.Nil(t, nil)
	assert.NoError(t, nil)
	assert.NotEmpty(t, "value")
	assert.NotZero(t, 42)
	assert.NotZero(t, "value")
	assert.Zero(t, nil)
	assert.Positive(t, 42)
	assert.True(t, true)
	assert.Zero(t, 0)
	assert.Zero(t, "")
	assert.Zero(t, nil)
}

func TestPositiveZero(t *testing.T) {
	assert.Positive(t, uint64(0))
}
