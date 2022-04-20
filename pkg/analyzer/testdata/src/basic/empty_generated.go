package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyAsserts(t *testing.T) {
	var (
		a    [0]int
		aPtr *[0]int
		s    []int
		m    map[int]int
		ss   string
		c    chan int
	)

	t.Run("assert", func(t *testing.T) {
		{
			assert.Len(t, a, 0)                         // want "use assert.Empty"
			assert.Len(t, a, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, a, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, a, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, a, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(a), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(a), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(a), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(a), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(a), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(a))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(a), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(a), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(a), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(a), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(a) == 0)                         // want "use assert.Empty"
			assert.True(t, len(a) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(a) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(a) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(a) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(a))                         // want "use assert.Empty"
			assert.True(t, 0 == len(a), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(a), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(a), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(a), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		{
			assert.Len(t, aPtr, 0)                         // want "use assert.Empty"
			assert.Len(t, aPtr, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, aPtr, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, aPtr, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, aPtr, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(aPtr), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(aPtr), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(aPtr), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(aPtr), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(aPtr), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(aPtr))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(aPtr), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(aPtr), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(aPtr), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(aPtr), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(aPtr) == 0)                         // want "use assert.Empty"
			assert.True(t, len(aPtr) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(aPtr) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(aPtr) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(aPtr) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(aPtr))                         // want "use assert.Empty"
			assert.True(t, 0 == len(aPtr), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(aPtr), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(aPtr), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(aPtr), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		{
			assert.Len(t, s, 0)                         // want "use assert.Empty"
			assert.Len(t, s, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, s, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, s, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, s, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(s), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(s), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(s), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(s), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(s), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(s))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(s), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(s), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(s), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(s), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(s) == 0)                         // want "use assert.Empty"
			assert.True(t, len(s) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(s) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(s) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(s) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(s))                         // want "use assert.Empty"
			assert.True(t, 0 == len(s), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(s), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(s), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(s), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		{
			assert.Len(t, m, 0)                         // want "use assert.Empty"
			assert.Len(t, m, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, m, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, m, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, m, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(m), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(m), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(m), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(m), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(m), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(m))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(m), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(m), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(m), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(m), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(m) == 0)                         // want "use assert.Empty"
			assert.True(t, len(m) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(m) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(m) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(m) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(m))                         // want "use assert.Empty"
			assert.True(t, 0 == len(m), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(m), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(m), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(m), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		{
			assert.Len(t, ss, 0)                         // want "use assert.Empty"
			assert.Len(t, ss, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, ss, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, ss, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, ss, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(ss), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(ss), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(ss), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(ss), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(ss), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(ss))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(ss), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(ss), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(ss), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(ss), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(ss) == 0)                         // want "use assert.Empty"
			assert.True(t, len(ss) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(ss) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(ss) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(ss) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(ss))                         // want "use assert.Empty"
			assert.True(t, 0 == len(ss), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(ss), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(ss), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(ss), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		{
			assert.Len(t, c, 0)                         // want "use assert.Empty"
			assert.Len(t, c, 0, "msg")                  // want "use assert.Empty"
			assert.Len(t, c, 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Lenf(t, c, 0, "msg")                 // want "use assert.Emptyf"
			assert.Lenf(t, c, 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, len(c), 0)                         // want "use assert.Empty"
			assert.Equal(t, len(c), 0, "msg")                  // want "use assert.Empty"
			assert.Equal(t, len(c), 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, len(c), 0, "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, len(c), 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.Equal(t, 0, len(c))                         // want "use assert.Empty"
			assert.Equal(t, 0, len(c), "msg")                  // want "use assert.Empty"
			assert.Equal(t, 0, len(c), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Equalf(t, 0, len(c), "msg")                 // want "use assert.Emptyf"
			assert.Equalf(t, 0, len(c), "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, len(c) == 0)                         // want "use assert.Empty"
			assert.True(t, len(c) == 0, "msg")                  // want "use assert.Empty"
			assert.True(t, len(c) == 0, "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, len(c) == 0, "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, len(c) == 0, "msg with arg %d", 42) // want "use assert.Emptyf"

			assert.True(t, 0 == len(c))                         // want "use assert.Empty"
			assert.True(t, 0 == len(c), "msg")                  // want "use assert.Empty"
			assert.True(t, 0 == len(c), "msg with arg %d", 42)  // want "use assert.Empty"
			assert.Truef(t, 0 == len(c), "msg")                 // want "use assert.Emptyf"
			assert.Truef(t, 0 == len(c), "msg with arg %d", 42) // want "use assert.Emptyf"
		}

		// Valid asserts.

		{
			assert.Empty(t, a)
			assert.Empty(t, a, "msg")
			assert.Empty(t, a, "msg with arg %d", 42)
			assert.Emptyf(t, a, "msg")
			assert.Emptyf(t, a, "msg with arg %d", 42)
		}

		{
			assert.Empty(t, aPtr)
			assert.Empty(t, aPtr, "msg")
			assert.Empty(t, aPtr, "msg with arg %d", 42)
			assert.Emptyf(t, aPtr, "msg")
			assert.Emptyf(t, aPtr, "msg with arg %d", 42)
		}

		{
			assert.Empty(t, s)
			assert.Empty(t, s, "msg")
			assert.Empty(t, s, "msg with arg %d", 42)
			assert.Emptyf(t, s, "msg")
			assert.Emptyf(t, s, "msg with arg %d", 42)
		}

		{
			assert.Empty(t, m)
			assert.Empty(t, m, "msg")
			assert.Empty(t, m, "msg with arg %d", 42)
			assert.Emptyf(t, m, "msg")
			assert.Emptyf(t, m, "msg with arg %d", 42)
		}

		{
			assert.Empty(t, ss)
			assert.Empty(t, ss, "msg")
			assert.Empty(t, ss, "msg with arg %d", 42)
			assert.Emptyf(t, ss, "msg")
			assert.Emptyf(t, ss, "msg with arg %d", 42)
		}

		{
			assert.Empty(t, c)
			assert.Empty(t, c, "msg")
			assert.Empty(t, c, "msg with arg %d", 42)
			assert.Emptyf(t, c, "msg")
			assert.Emptyf(t, c, "msg with arg %d", 42)
		}
	})

	t.Run("require", func(t *testing.T) {
		{
			require.Len(t, a, 0)                         // want "use require.Empty"
			require.Len(t, a, 0, "msg")                  // want "use require.Empty"
			require.Len(t, a, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, a, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, a, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(a), 0)                         // want "use require.Empty"
			require.Equal(t, len(a), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(a), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(a), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(a), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(a))                         // want "use require.Empty"
			require.Equal(t, 0, len(a), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(a), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(a), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(a), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(a) == 0)                         // want "use require.Empty"
			require.True(t, len(a) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(a) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(a) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(a) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(a))                         // want "use require.Empty"
			require.True(t, 0 == len(a), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(a), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(a), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(a), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		{
			require.Len(t, aPtr, 0)                         // want "use require.Empty"
			require.Len(t, aPtr, 0, "msg")                  // want "use require.Empty"
			require.Len(t, aPtr, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, aPtr, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, aPtr, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(aPtr), 0)                         // want "use require.Empty"
			require.Equal(t, len(aPtr), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(aPtr), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(aPtr), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(aPtr), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(aPtr))                         // want "use require.Empty"
			require.Equal(t, 0, len(aPtr), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(aPtr), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(aPtr), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(aPtr), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(aPtr) == 0)                         // want "use require.Empty"
			require.True(t, len(aPtr) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(aPtr) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(aPtr) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(aPtr) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(aPtr))                         // want "use require.Empty"
			require.True(t, 0 == len(aPtr), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(aPtr), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(aPtr), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(aPtr), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		{
			require.Len(t, s, 0)                         // want "use require.Empty"
			require.Len(t, s, 0, "msg")                  // want "use require.Empty"
			require.Len(t, s, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, s, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, s, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(s), 0)                         // want "use require.Empty"
			require.Equal(t, len(s), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(s), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(s), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(s), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(s))                         // want "use require.Empty"
			require.Equal(t, 0, len(s), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(s), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(s), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(s), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(s) == 0)                         // want "use require.Empty"
			require.True(t, len(s) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(s) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(s) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(s) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(s))                         // want "use require.Empty"
			require.True(t, 0 == len(s), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(s), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(s), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(s), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		{
			require.Len(t, m, 0)                         // want "use require.Empty"
			require.Len(t, m, 0, "msg")                  // want "use require.Empty"
			require.Len(t, m, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, m, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, m, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(m), 0)                         // want "use require.Empty"
			require.Equal(t, len(m), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(m), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(m), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(m), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(m))                         // want "use require.Empty"
			require.Equal(t, 0, len(m), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(m), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(m), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(m), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(m) == 0)                         // want "use require.Empty"
			require.True(t, len(m) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(m) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(m) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(m) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(m))                         // want "use require.Empty"
			require.True(t, 0 == len(m), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(m), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(m), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(m), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		{
			require.Len(t, ss, 0)                         // want "use require.Empty"
			require.Len(t, ss, 0, "msg")                  // want "use require.Empty"
			require.Len(t, ss, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, ss, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, ss, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(ss), 0)                         // want "use require.Empty"
			require.Equal(t, len(ss), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(ss), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(ss), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(ss), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(ss))                         // want "use require.Empty"
			require.Equal(t, 0, len(ss), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(ss), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(ss), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(ss), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(ss) == 0)                         // want "use require.Empty"
			require.True(t, len(ss) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(ss) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(ss) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(ss) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(ss))                         // want "use require.Empty"
			require.True(t, 0 == len(ss), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(ss), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(ss), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(ss), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		{
			require.Len(t, c, 0)                         // want "use require.Empty"
			require.Len(t, c, 0, "msg")                  // want "use require.Empty"
			require.Len(t, c, 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Lenf(t, c, 0, "msg")                 // want "use require.Emptyf"
			require.Lenf(t, c, 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, len(c), 0)                         // want "use require.Empty"
			require.Equal(t, len(c), 0, "msg")                  // want "use require.Empty"
			require.Equal(t, len(c), 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, len(c), 0, "msg")                 // want "use require.Emptyf"
			require.Equalf(t, len(c), 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.Equal(t, 0, len(c))                         // want "use require.Empty"
			require.Equal(t, 0, len(c), "msg")                  // want "use require.Empty"
			require.Equal(t, 0, len(c), "msg with arg %d", 42)  // want "use require.Empty"
			require.Equalf(t, 0, len(c), "msg")                 // want "use require.Emptyf"
			require.Equalf(t, 0, len(c), "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, len(c) == 0)                         // want "use require.Empty"
			require.True(t, len(c) == 0, "msg")                  // want "use require.Empty"
			require.True(t, len(c) == 0, "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, len(c) == 0, "msg")                 // want "use require.Emptyf"
			require.Truef(t, len(c) == 0, "msg with arg %d", 42) // want "use require.Emptyf"

			require.True(t, 0 == len(c))                         // want "use require.Empty"
			require.True(t, 0 == len(c), "msg")                  // want "use require.Empty"
			require.True(t, 0 == len(c), "msg with arg %d", 42)  // want "use require.Empty"
			require.Truef(t, 0 == len(c), "msg")                 // want "use require.Emptyf"
			require.Truef(t, 0 == len(c), "msg with arg %d", 42) // want "use require.Emptyf"
		}

		// Valid requires.

		{
			require.Empty(t, a)
			require.Empty(t, a, "msg")
			require.Empty(t, a, "msg with arg %d", 42)
			require.Emptyf(t, a, "msg")
			require.Emptyf(t, a, "msg with arg %d", 42)
		}

		{
			require.Empty(t, aPtr)
			require.Empty(t, aPtr, "msg")
			require.Empty(t, aPtr, "msg with arg %d", 42)
			require.Emptyf(t, aPtr, "msg")
			require.Emptyf(t, aPtr, "msg with arg %d", 42)
		}

		{
			require.Empty(t, s)
			require.Empty(t, s, "msg")
			require.Empty(t, s, "msg with arg %d", 42)
			require.Emptyf(t, s, "msg")
			require.Emptyf(t, s, "msg with arg %d", 42)
		}

		{
			require.Empty(t, m)
			require.Empty(t, m, "msg")
			require.Empty(t, m, "msg with arg %d", 42)
			require.Emptyf(t, m, "msg")
			require.Emptyf(t, m, "msg with arg %d", 42)
		}

		{
			require.Empty(t, ss)
			require.Empty(t, ss, "msg")
			require.Empty(t, ss, "msg with arg %d", 42)
			require.Emptyf(t, ss, "msg")
			require.Emptyf(t, ss, "msg with arg %d", 42)
		}

		{
			require.Empty(t, c)
			require.Empty(t, c, "msg")
			require.Empty(t, c, "msg with arg %d", 42)
			require.Emptyf(t, c, "msg")
			require.Emptyf(t, c, "msg with arg %d", 42)
		}
	})
}
