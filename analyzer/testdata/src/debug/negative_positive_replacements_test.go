package debug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegativeReplacements(t *testing.T) {
	var tm tMock

	negativeAssertions := []func(v int) bool{
		func(v int) bool { return assert.Less(tm, v, 0) },
		func(v int) bool { return assert.Greater(tm, 0, v) },
		func(v int) bool { return assert.True(tm, v < 0) },
		func(v int) bool { return assert.True(tm, 0 > v) },
		func(v int) bool { return assert.False(tm, v >= 0) },
		func(v int) bool { return assert.False(tm, 0 <= v) },
	}

	t.Run("assert.Negative", func(t *testing.T) {
		for i, original := range negativeAssertions {
			for _, v := range []int{-1, 0, 1} {
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

	positiveAssertions := []func(v int) bool{
		func(v int) bool { return assert.Greater(tm, v, 0) },
		func(v int) bool { return assert.Less(tm, 0, v) },
		func(v int) bool { return assert.True(tm, v > 0) },
		func(v int) bool { return assert.True(tm, 0 < v) },
		func(v int) bool { return assert.False(tm, v <= 0) },
		func(v int) bool { return assert.False(tm, 0 >= v) },
	}

	t.Run("assert.Positive", func(t *testing.T) {
		for i, original := range positiveAssertions {
			for _, v := range []int{-1, 0, 1} {
				t.Run(fmt.Sprintf("%d_%d", i, v), func(t *testing.T) {
					replacement := assert.Positive(tm, v)
					assert.Equal(t, original(v), replacement, "not an equivalent replacement")
				})
			}
		}
	})
}
