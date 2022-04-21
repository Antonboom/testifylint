package basic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64Compare(t *testing.T) {
	type number float64

	var a float64
	var b number
	var s struct{ c float64 }
	d := 1.01
	const e = 2.02

	assert.Equal(t, a, 1.01) // use "assert.InDelta"
	assert.NotEqual(t, a, 1.01)
	assert.Greater(t, 1.01, a)
	assert.GreaterOrEqual(t, 1.01, a)
	assert.Less(t, 1.01, a)
	assert.LessOrEqual(t, 1.01, a)

	assert.Equal(t, 1.01, a) // use "assert.InDelta"
	assert.NotEqual(t, 1.01, a)
	assert.Greater(t, a, 1.01)
	assert.GreaterOrEqual(t, a, 1.01)
	assert.Less(t, a, 1.01)
	assert.LessOrEqual(t, a, 1.01)

	assert.True(t, 1.01 == a)
	assert.True(t, 1.01 != a)
	assert.True(t, 1.01 > a)
	assert.True(t, 1.01 >= a)
	assert.True(t, 1.01 < a)
	assert.True(t, 1.01 <= a)

	assert.True(t, a == 1.01)
	assert.True(t, a != 1.01)
	assert.True(t, a > 1.01)
	assert.True(t, a >= 1.01)
	assert.True(t, a < 1.01)
	assert.True(t, a <= 1.01)

	assert.False(t, 1.01 == a)
	assert.False(t, 1.01 != a)
	assert.False(t, 1.01 > a)
	assert.False(t, 1.01 >= a)
	assert.False(t, 1.01 < a)
	assert.False(t, 1.01 <= a)

	assert.False(t, a == 1.01)
	assert.False(t, a != 1.01)
	assert.False(t, a > 1.01)
	assert.False(t, a >= 1.01)
	assert.False(t, a < 1.01)
	assert.False(t, a <= 1.01)

	assert.InDelta(t, a, 1.01, 0.01)
	assert.InDelta(t, 1.01, a, 0.01)

	// a, numerical
	// b, s.c
	// d, e
}

// "don't use Equal for float"
