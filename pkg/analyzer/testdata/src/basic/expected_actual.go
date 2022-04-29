package basic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfusedWithExpectedActual(t *testing.T) {
	var result int

	const (
		a = uint(11)
		b = uint8(12)
		c = uint16(13)
		d = uint32(14)
		e = uint64(15)

		f = int(21)
		g = int8(22)
		h = int16(23)
		i = int32(24)
		j = int64(25)

		k = float32(31.)
		l = float64(32.)

		m = complex64(41 - 0.707i)
		n = complex128(42 - 0.707i)

		o = "string"
		p = 'r'
	)

	const (
		Sunday = iota
		Monday
	)

	type Day int
	const (
		DaySunday = iota
		DayMonday
	)

	var expected int
	var tt struct{ expected int }
	ttp := &struct{ expected int }{}

	assert.Equal(t, result, uint(11))
	assert.Equal(t, result, uint8(12))
	assert.Equal(t, result, uint16(13))
	assert.Equal(t, result, uint32(14))
	assert.Equal(t, result, uint64(15))
	assert.Equal(t, result, int(21))
	assert.Equal(t, result, int8(22))
	assert.Equal(t, result, int16(23))
	assert.Equal(t, result, int32(24))
	assert.Equal(t, result, int64(25))
	assert.Equal(t, result, float32(31.))
	assert.Equal(t, result, float64(32.))
	assert.Equal(t, result, complex64(41-0.707i))
	assert.Equal(t, result, complex128(42-0.707i))
	assert.Equal(t, result, "string")
	assert.Equal(t, result, 'r')

	assert.Equal(t, result, a)
	assert.Equal(t, result, b)
	assert.Equal(t, result, c)
	assert.Equal(t, result, d)
	assert.Equal(t, result, e)
	assert.Equal(t, result, f)
	assert.Equal(t, result, g)
	assert.Equal(t, result, h)
	assert.Equal(t, result, i)
	assert.Equal(t, result, j)
	assert.Equal(t, result, k)
	assert.Equal(t, result, l)

	assert.Equal(t, result, Monday)
	assert.Equal(t, result, DayMonday)

	assert.Equal(t, result, expected)
	assert.Equal(t, result, tt.expected)
	assert.Equal(t, result, ttp.expected)
}
