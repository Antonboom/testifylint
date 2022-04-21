package basic

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLen(t *testing.T) {
	a := [...]int{1, 2, 3}
	aPtr := &a
	s := []int{1, 2, 3}
	m := map[int]int{1: 1, 2: 2, 3: 3}
	ss := "go!"
	c := make(chan int, 3)

	t.Run("assert", func(t *testing.T) {
		{
			assert.Equal(t, len(a), 3)                         // want "use assert.Len"
			assert.Equal(t, len(a), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(a), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(a), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(a), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(a))                         // want "use assert.Len"
			assert.Equal(t, 3, len(a), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(a), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(a), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(a), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(a) == 3)                         // want "use assert.Len"
			assert.True(t, len(a) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(a) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(a) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(a) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(a))                         // want "use assert.Len"
			assert.True(t, 3 == len(a), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(a), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(a), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(a), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		{
			assert.Equal(t, len(aPtr), 3)                         // want "use assert.Len"
			assert.Equal(t, len(aPtr), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(aPtr), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(aPtr), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(aPtr), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(aPtr))                         // want "use assert.Len"
			assert.Equal(t, 3, len(aPtr), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(aPtr), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(aPtr), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(aPtr), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(aPtr) == 3)                         // want "use assert.Len"
			assert.True(t, len(aPtr) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(aPtr) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(aPtr) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(aPtr) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(aPtr))                         // want "use assert.Len"
			assert.True(t, 3 == len(aPtr), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(aPtr), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(aPtr), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(aPtr), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		{
			assert.Equal(t, len(s), 3)                         // want "use assert.Len"
			assert.Equal(t, len(s), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(s), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(s), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(s), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(s))                         // want "use assert.Len"
			assert.Equal(t, 3, len(s), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(s), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(s), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(s), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(s) == 3)                         // want "use assert.Len"
			assert.True(t, len(s) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(s) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(s) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(s) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(s))                         // want "use assert.Len"
			assert.True(t, 3 == len(s), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(s), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(s), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(s), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		{
			assert.Equal(t, len(m), 3)                         // want "use assert.Len"
			assert.Equal(t, len(m), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(m), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(m), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(m), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(m))                         // want "use assert.Len"
			assert.Equal(t, 3, len(m), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(m), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(m), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(m), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(m) == 3)                         // want "use assert.Len"
			assert.True(t, len(m) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(m) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(m) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(m) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(m))                         // want "use assert.Len"
			assert.True(t, 3 == len(m), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(m), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(m), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(m), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		{
			assert.Equal(t, len(ss), 3)                         // want "use assert.Len"
			assert.Equal(t, len(ss), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(ss), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(ss), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(ss), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(ss))                         // want "use assert.Len"
			assert.Equal(t, 3, len(ss), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(ss), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(ss), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(ss), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(ss) == 3)                         // want "use assert.Len"
			assert.True(t, len(ss) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(ss) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(ss) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(ss) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(ss))                         // want "use assert.Len"
			assert.True(t, 3 == len(ss), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(ss), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(ss), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(ss), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		{
			assert.Equal(t, len(c), 3)                         // want "use assert.Len"
			assert.Equal(t, len(c), 3, "msg")                  // want "use assert.Len"
			assert.Equal(t, len(c), 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, len(c), 3, "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, len(c), 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.Equal(t, 3, len(c))                         // want "use assert.Len"
			assert.Equal(t, 3, len(c), "msg")                  // want "use assert.Len"
			assert.Equal(t, 3, len(c), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Equalf(t, 3, len(c), "msg")                 // want "use assert.Lenf"
			assert.Equalf(t, 3, len(c), "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, len(c) == 3)                         // want "use assert.Len"
			assert.True(t, len(c) == 3, "msg")                  // want "use assert.Len"
			assert.True(t, len(c) == 3, "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, len(c) == 3, "msg")                 // want "use assert.Lenf"
			assert.Truef(t, len(c) == 3, "msg with arg %d", 42) // want "use assert.Lenf"

			assert.True(t, 3 == len(c))                         // want "use assert.Len"
			assert.True(t, 3 == len(c), "msg")                  // want "use assert.Len"
			assert.True(t, 3 == len(c), "msg with arg %d", 42)  // want "use assert.Len"
			assert.Truef(t, 3 == len(c), "msg")                 // want "use assert.Lenf"
			assert.Truef(t, 3 == len(c), "msg with arg %d", 42) // want "use assert.Lenf"
		}

		// Valid asserts.

		{
			assert.Len(t, a, 3)
			assert.Len(t, a, 3, "msg")
			assert.Len(t, a, 3, "msg with arg %d", 42)
			assert.Lenf(t, a, 3, "msg")
			assert.Lenf(t, a, 3, "msg with arg %d", 42)

			assert.Len(t, aPtr, 3)
			assert.Len(t, aPtr, 3, "msg")
			assert.Len(t, aPtr, 3, "msg with arg %d", 42)
			assert.Lenf(t, aPtr, 3, "msg")
			assert.Lenf(t, aPtr, 3, "msg with arg %d", 42)

			assert.Len(t, s, 3)
			assert.Len(t, s, 3, "msg")
			assert.Len(t, s, 3, "msg with arg %d", 42)
			assert.Lenf(t, s, 3, "msg")
			assert.Lenf(t, s, 3, "msg with arg %d", 42)

			assert.Len(t, m, 3)
			assert.Len(t, m, 3, "msg")
			assert.Len(t, m, 3, "msg with arg %d", 42)
			assert.Lenf(t, m, 3, "msg")
			assert.Lenf(t, m, 3, "msg with arg %d", 42)

			assert.Len(t, ss, 3)
			assert.Len(t, ss, 3, "msg")
			assert.Len(t, ss, 3, "msg with arg %d", 42)
			assert.Lenf(t, ss, 3, "msg")
			assert.Lenf(t, ss, 3, "msg with arg %d", 42)

			assert.Len(t, c, 3)
			assert.Len(t, c, 3, "msg")
			assert.Len(t, c, 3, "msg with arg %d", 42)
			assert.Lenf(t, c, 3, "msg")
			assert.Lenf(t, c, 3, "msg with arg %d", 42)
		}

		{
			assert.Equal(t, len(a), len(s))
			assert.Equal(t, len(s), len(a), "msg")
			assert.Equal(t, len(a), len(aPtr), "msg with arg %d", 42)
			assert.Equalf(t, len(m), len(ss), "msg")
			assert.Equalf(t, len(c), len(aPtr), "msg with arg %d", 42)

			assert.NotEqual(t, len(a), len(s))
			assert.NotEqual(t, len(s), len(a), "msg")
			assert.NotEqual(t, len(a), len(aPtr), "msg with arg %d", 42)
			assert.NotEqualf(t, len(m), len(ss), "msg")
			assert.NotEqualf(t, len(c), len(aPtr), "msg with arg %d", 42)

			assert.True(t, len(a) == len(s))
			assert.True(t, len(a) == len(s), "msg")
			assert.True(t, len(ss) == len(s), "msg with arg %d", 42)
			assert.Truef(t, len(aPtr) == len(m), "msg")
			assert.Truef(t, len(c) == len(c), "msg with arg %d", 42)
		}
	})

	t.Run("require", func(t *testing.T) {
		{
			require.Equal(t, len(a), 3)                         // want "use require.Len"
			require.Equal(t, len(a), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(a), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(a), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(a), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(a))                         // want "use require.Len"
			require.Equal(t, 3, len(a), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(a), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(a), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(a), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(a) == 3)                         // want "use require.Len"
			require.True(t, len(a) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(a) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(a) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(a) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(a))                         // want "use require.Len"
			require.True(t, 3 == len(a), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(a), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(a), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(a), "msg with arg %d", 42) // want "use require.Lenf"
		}

		{
			require.Equal(t, len(aPtr), 3)                         // want "use require.Len"
			require.Equal(t, len(aPtr), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(aPtr), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(aPtr), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(aPtr), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(aPtr))                         // want "use require.Len"
			require.Equal(t, 3, len(aPtr), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(aPtr), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(aPtr), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(aPtr), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(aPtr) == 3)                         // want "use require.Len"
			require.True(t, len(aPtr) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(aPtr) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(aPtr) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(aPtr) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(aPtr))                         // want "use require.Len"
			require.True(t, 3 == len(aPtr), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(aPtr), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(aPtr), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(aPtr), "msg with arg %d", 42) // want "use require.Lenf"
		}

		{
			require.Equal(t, len(s), 3)                         // want "use require.Len"
			require.Equal(t, len(s), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(s), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(s), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(s), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(s))                         // want "use require.Len"
			require.Equal(t, 3, len(s), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(s), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(s), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(s), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(s) == 3)                         // want "use require.Len"
			require.True(t, len(s) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(s) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(s) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(s) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(s))                         // want "use require.Len"
			require.True(t, 3 == len(s), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(s), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(s), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(s), "msg with arg %d", 42) // want "use require.Lenf"
		}

		{
			require.Equal(t, len(m), 3)                         // want "use require.Len"
			require.Equal(t, len(m), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(m), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(m), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(m), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(m))                         // want "use require.Len"
			require.Equal(t, 3, len(m), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(m), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(m), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(m), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(m) == 3)                         // want "use require.Len"
			require.True(t, len(m) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(m) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(m) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(m) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(m))                         // want "use require.Len"
			require.True(t, 3 == len(m), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(m), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(m), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(m), "msg with arg %d", 42) // want "use require.Lenf"
		}

		{
			require.Equal(t, len(ss), 3)                         // want "use require.Len"
			require.Equal(t, len(ss), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(ss), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(ss), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(ss), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(ss))                         // want "use require.Len"
			require.Equal(t, 3, len(ss), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(ss), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(ss), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(ss), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(ss) == 3)                         // want "use require.Len"
			require.True(t, len(ss) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(ss) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(ss) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(ss) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(ss))                         // want "use require.Len"
			require.True(t, 3 == len(ss), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(ss), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(ss), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(ss), "msg with arg %d", 42) // want "use require.Lenf"
		}

		{
			require.Equal(t, len(c), 3)                         // want "use require.Len"
			require.Equal(t, len(c), 3, "msg")                  // want "use require.Len"
			require.Equal(t, len(c), 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, len(c), 3, "msg")                 // want "use require.Lenf"
			require.Equalf(t, len(c), 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.Equal(t, 3, len(c))                         // want "use require.Len"
			require.Equal(t, 3, len(c), "msg")                  // want "use require.Len"
			require.Equal(t, 3, len(c), "msg with arg %d", 42)  // want "use require.Len"
			require.Equalf(t, 3, len(c), "msg")                 // want "use require.Lenf"
			require.Equalf(t, 3, len(c), "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, len(c) == 3)                         // want "use require.Len"
			require.True(t, len(c) == 3, "msg")                  // want "use require.Len"
			require.True(t, len(c) == 3, "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, len(c) == 3, "msg")                 // want "use require.Lenf"
			require.Truef(t, len(c) == 3, "msg with arg %d", 42) // want "use require.Lenf"

			require.True(t, 3 == len(c))                         // want "use require.Len"
			require.True(t, 3 == len(c), "msg")                  // want "use require.Len"
			require.True(t, 3 == len(c), "msg with arg %d", 42)  // want "use require.Len"
			require.Truef(t, 3 == len(c), "msg")                 // want "use require.Lenf"
			require.Truef(t, 3 == len(c), "msg with arg %d", 42) // want "use require.Lenf"
		}

		// Valid requires.

		{
			require.Len(t, a, 3)
			require.Len(t, a, 3, "msg")
			require.Len(t, a, 3, "msg with arg %d", 42)
			require.Lenf(t, a, 3, "msg")
			require.Lenf(t, a, 3, "msg with arg %d", 42)

			require.Len(t, aPtr, 3)
			require.Len(t, aPtr, 3, "msg")
			require.Len(t, aPtr, 3, "msg with arg %d", 42)
			require.Lenf(t, aPtr, 3, "msg")
			require.Lenf(t, aPtr, 3, "msg with arg %d", 42)

			require.Len(t, s, 3)
			require.Len(t, s, 3, "msg")
			require.Len(t, s, 3, "msg with arg %d", 42)
			require.Lenf(t, s, 3, "msg")
			require.Lenf(t, s, 3, "msg with arg %d", 42)

			require.Len(t, m, 3)
			require.Len(t, m, 3, "msg")
			require.Len(t, m, 3, "msg with arg %d", 42)
			require.Lenf(t, m, 3, "msg")
			require.Lenf(t, m, 3, "msg with arg %d", 42)

			require.Len(t, ss, 3)
			require.Len(t, ss, 3, "msg")
			require.Len(t, ss, 3, "msg with arg %d", 42)
			require.Lenf(t, ss, 3, "msg")
			require.Lenf(t, ss, 3, "msg with arg %d", 42)

			require.Len(t, c, 3)
			require.Len(t, c, 3, "msg")
			require.Len(t, c, 3, "msg with arg %d", 42)
			require.Lenf(t, c, 3, "msg")
			require.Lenf(t, c, 3, "msg with arg %d", 42)
		}

		{
			require.Equal(t, len(a), len(s))
			require.Equal(t, len(s), len(a), "msg")
			require.Equal(t, len(a), len(aPtr), "msg with arg %d", 42)
			require.Equalf(t, len(m), len(ss), "msg")
			require.Equalf(t, len(c), len(aPtr), "msg with arg %d", 42)

			require.NotEqual(t, len(a), len(s))
			require.NotEqual(t, len(s), len(a), "msg")
			require.NotEqual(t, len(a), len(aPtr), "msg with arg %d", 42)
			require.NotEqualf(t, len(m), len(ss), "msg")
			require.NotEqualf(t, len(c), len(aPtr), "msg with arg %d", 42)

			require.True(t, len(a) == len(s))
			require.True(t, len(a) == len(s), "msg")
			require.True(t, len(ss) == len(s), "msg with arg %d", 42)
			require.Truef(t, len(aPtr) == len(m), "msg")
			require.Truef(t, len(c) == len(c), "msg with arg %d", 42)

			//TODO:
			// >= len(s)
			// <, >, !=, NotEqual
		}
	})
}
