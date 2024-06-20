package debug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegativeReplacements(t *testing.T) {
	var tm tMock

	negativeAssertionsInt := []func(v int) bool{
		func(v int) bool { return assert.Less(tm, v, 0) },
		func(v int) bool { return assert.Greater(tm, 0, v) },
		func(v int) bool { return assert.True(tm, v < 0) },
		func(v int) bool { return assert.True(tm, 0 > v) },
		func(v int) bool { return assert.False(tm, v >= 0) },
		func(v int) bool { return assert.False(tm, 0 <= v) },
	}

	t.Run("int", func(t *testing.T) {
		for i, original := range negativeAssertionsInt {
			for _, v := range []int{-1, 0, 1} {
				t.Run(fmt.Sprintf("%d_%d", i, v), func(t *testing.T) {
					replacement := assert.Negative(tm, v)
					assert.Equal(t, original(v), replacement, "not an equivalent replacement")
				})
			}
		}
	})

	negativeAssertionsUint := []func(v uint8) bool{
		func(v uint8) bool { return assert.Less(tm, v, uint8(0)) },
		func(v uint8) bool { return assert.Greater(tm, uint8(0), v) },
		func(v uint8) bool { return assert.True(tm, v < uint8(0)) },
		func(v uint8) bool { return assert.True(tm, uint8(0) > v) },
		func(v uint8) bool { return assert.False(tm, v >= uint8(0)) },
		func(v uint8) bool { return assert.False(tm, uint8(0) <= v) },
	}

	t.Run("uint", func(t *testing.T) {
		for i, original := range negativeAssertionsUint {
			for _, v := range []uint8{ /* -1, */ 0, 1} { // constant -1 overflows uint8
				t.Run(fmt.Sprintf("%d_%d", i, v), func(t *testing.T) {
					replacement := assert.Negative(tm, v)
					assert.Equal(t, original(v), replacement, "not an equivalent replacement")
				})
			}
		}
	})
}

func TestPositiveReplacements(t *testing.T) {
	var tm tMock

	positiveAssertionsInt := []func(v int) bool{
		func(v int) bool { return assert.Greater(tm, v, 0) },
		func(v int) bool { return assert.Less(tm, 0, v) },
		func(v int) bool { return assert.True(tm, v > 0) },
		func(v int) bool { return assert.True(tm, 0 < v) },
		func(v int) bool { return assert.False(tm, v <= 0) },
		func(v int) bool { return assert.False(tm, 0 >= v) },
	}

	t.Run("int", func(t *testing.T) {
		for i, original := range positiveAssertionsInt {
			for _, v := range []int{-1, 0, 1} {
				t.Run(fmt.Sprintf("%d_%d", i, v), func(t *testing.T) {
					replacement := assert.Positive(tm, v)
					assert.Equal(t, original(v), replacement, "not an equivalent replacement")
				})
			}
		}
	})

	positiveAssertionsUint := []func(v uint8) bool{
		func(v uint8) bool { return assert.Greater(tm, v, uint8(0)) },
		func(v uint8) bool { return assert.Less(tm, uint8(0), v) },
		func(v uint8) bool { return assert.True(tm, v > uint8(0)) },
		func(v uint8) bool { return assert.True(tm, uint8(0) < v) },
		func(v uint8) bool { return assert.False(tm, v <= uint8(0)) },
		func(v uint8) bool { return assert.False(tm, uint8(0) >= v) },
	}

	t.Run("uint", func(t *testing.T) {
		for i, original := range positiveAssertionsUint {
			for _, v := range []uint8{ /* -1, */ 0, 1} { // constant -1 overflows uint8
				t.Run(fmt.Sprintf("%d_%d", i, v), func(t *testing.T) {
					replacement := assert.Positive(tm, v)
					assert.Equal(t, original(v), replacement, "not an equivalent replacement")
				})
			}
		}
	})
}
