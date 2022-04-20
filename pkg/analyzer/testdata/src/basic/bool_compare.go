package basic

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoolAsserts(t *testing.T) {
	var b bool

	t.Run("assert", func(t *testing.T) {
		{
			assert.Equal(t, b, true)                         // want "use assert.True"
			assert.Equal(t, b, true, "msg")                  // want "use assert.True"
			assert.Equal(t, b, true, "msg with arg %d", 42)  // want "use assert.True"
			assert.Equalf(t, b, true, "msg")                 // want "use assert.Truef"
			assert.Equalf(t, b, true, "msg with arg %d", 42) // want "use assert.Truef"

			assert.Equal(t, true, b)                         // want "use assert.True"
			assert.Equal(t, true, b, "msg")                  // want "use assert.True"
			assert.Equal(t, true, b, "msg with arg %d", 42)  // want "use assert.True"
			assert.Equalf(t, true, b, "msg")                 // want "use assert.Truef"
			assert.Equalf(t, true, b, "msg with arg %d", 42) // want "use assert.Truef"

			assert.NotEqual(t, b, false)                         // want "use assert.True"
			assert.NotEqual(t, b, false, "msg")                  // want "use assert.True"
			assert.NotEqual(t, b, false, "msg with arg %d", 42)  // want "use assert.True"
			assert.NotEqualf(t, b, false, "msg")                 // want "use assert.Truef"
			assert.NotEqualf(t, b, false, "msg with arg %d", 42) // want "use assert.Truef"

			assert.NotEqual(t, false, b)                         // want "use assert.True"
			assert.NotEqual(t, false, b, "msg")                  // want "use assert.True"
			assert.NotEqual(t, false, b, "msg with arg %d", 42)  // want "use assert.True"
			assert.NotEqualf(t, false, b, "msg")                 // want "use assert.Truef"
			assert.NotEqualf(t, false, b, "msg with arg %d", 42) // want "use assert.Truef"

			assert.True(t, b == true)                         // want "need to simplify the check"
			assert.True(t, b == true, "msg")                  // want "need to simplify the check"
			assert.True(t, b == true, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Truef(t, b == true, "msg")                 // want "need to simplify the check"
			assert.Truef(t, b == true, "msg with arg %d", 42) // want "need to simplify the check"

			assert.True(t, true == b)                         // want "need to simplify the check"
			assert.True(t, true == b, "msg")                  // want "need to simplify the check"
			assert.True(t, true == b, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Truef(t, true == b, "msg")                 // want "need to simplify the check"
			assert.Truef(t, true == b, "msg with arg %d", 42) // want "need to simplify the check"

			assert.False(t, b == false)                         // want "use assert.True"
			assert.False(t, b == false, "msg")                  // want "use assert.True"
			assert.False(t, b == false, "msg with arg %d", 42)  // want "use assert.True"
			assert.Falsef(t, b == false, "msg")                 // want "use assert.Truef"
			assert.Falsef(t, b == false, "msg with arg %d", 42) // want "use assert.Truef"

			assert.False(t, false == b)                         // want "use assert.True"
			assert.False(t, false == b, "msg")                  // want "use assert.True"
			assert.False(t, false == b, "msg with arg %d", 42)  // want "use assert.True"
			assert.Falsef(t, false == b, "msg")                 // want "use assert.Truef"
			assert.Falsef(t, false == b, "msg with arg %d", 42) // want "use assert.Truef"

			assert.False(t, b != true)                         // want "use assert.True"
			assert.False(t, b != true, "msg")                  // want "use assert.True"
			assert.False(t, b != true, "msg with arg %d", 42)  // want "use assert.True"
			assert.Falsef(t, b != true, "msg")                 // want "use assert.Truef"
			assert.Falsef(t, b != true, "msg with arg %d", 42) // want "use assert.Truef"

			assert.False(t, true != b)                         // want "use assert.True"
			assert.False(t, true != b, "msg")                  // want "use assert.True"
			assert.False(t, true != b, "msg with arg %d", 42)  // want "use assert.True"
			assert.Falsef(t, true != b, "msg")                 // want "use assert.Truef"
			assert.Falsef(t, true != b, "msg with arg %d", 42) // want "use assert.Truef"

			assert.True(t, b != false)                         // want "need to simplify the check"
			assert.True(t, b != false, "msg")                  // want "need to simplify the check"
			assert.True(t, b != false, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Truef(t, b != false, "msg")                 // want "need to simplify the check"
			assert.Truef(t, b != false, "msg with arg %d", 42) // want "need to simplify the check"

			assert.True(t, false != b)                         // want "need to simplify the check"
			assert.True(t, false != b, "msg")                  // want "need to simplify the check"
			assert.True(t, false != b, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Truef(t, false != b, "msg")                 // want "need to simplify the check"
			assert.Truef(t, false != b, "msg with arg %d", 42) // want "need to simplify the check"

			assert.False(t, !b)                         // want "use assert.True"
			assert.False(t, !b, "msg")                  // want "use assert.True"
			assert.False(t, !b, "msg with arg %d", 42)  // want "use assert.True"
			assert.Falsef(t, !b, "msg")                 // want "use assert.Truef"
			assert.Falsef(t, !b, "msg with arg %d", 42) // want "use assert.Truef"
		}

		{
			assert.Equal(t, b, false)                         // want "use assert.False"
			assert.Equal(t, b, false, "msg")                  // want "use assert.False"
			assert.Equal(t, b, false, "msg with arg %d", 42)  // want "use assert.False"
			assert.Equalf(t, b, false, "msg")                 // want "use assert.Falsef"
			assert.Equalf(t, b, false, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.Equal(t, false, b)                         // want "use assert.False"
			assert.Equal(t, false, b, "msg")                  // want "use assert.False"
			assert.Equal(t, false, b, "msg with arg %d", 42)  // want "use assert.False"
			assert.Equalf(t, false, b, "msg")                 // want "use assert.Falsef"
			assert.Equalf(t, false, b, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.NotEqual(t, b, true)                         // want "use assert.False"
			assert.NotEqual(t, b, true, "msg")                  // want "use assert.False"
			assert.NotEqual(t, b, true, "msg with arg %d", 42)  // want "use assert.False"
			assert.NotEqualf(t, b, true, "msg")                 // want "use assert.Falsef"
			assert.NotEqualf(t, b, true, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.NotEqual(t, true, b)                         // want "use assert.False"
			assert.NotEqual(t, true, b, "msg")                  // want "use assert.False"
			assert.NotEqual(t, true, b, "msg with arg %d", 42)  // want "use assert.False"
			assert.NotEqualf(t, true, b, "msg")                 // want "use assert.Falsef"
			assert.NotEqualf(t, true, b, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.False(t, b == true)                         // want "need to simplify the check"
			assert.False(t, b == true, "msg")                  // want "need to simplify the check"
			assert.False(t, b == true, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Falsef(t, b == true, "msg")                 // want "need to simplify the check"
			assert.Falsef(t, b == true, "msg with arg %d", 42) // want "need to simplify the check"

			assert.False(t, true == b)                         // want "need to simplify the check"
			assert.False(t, true == b, "msg")                  // want "need to simplify the check"
			assert.False(t, true == b, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Falsef(t, true == b, "msg")                 // want "need to simplify the check"
			assert.Falsef(t, true == b, "msg with arg %d", 42) // want "need to simplify the check"

			assert.True(t, b == false)                         // want "use assert.False"
			assert.True(t, b == false, "msg")                  // want "use assert.False"
			assert.True(t, b == false, "msg with arg %d", 42)  // want "use assert.False"
			assert.Truef(t, b == false, "msg")                 // want "use assert.Falsef"
			assert.Truef(t, b == false, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.True(t, false == b)                         // want "use assert.False"
			assert.True(t, false == b, "msg")                  // want "use assert.False"
			assert.True(t, false == b, "msg with arg %d", 42)  // want "use assert.False"
			assert.Truef(t, false == b, "msg")                 // want "use assert.Falsef"
			assert.Truef(t, false == b, "msg with arg %d", 42) // want "use assert.Falsef"

			assert.True(t, b != true)                         // want "use assert.False"
			assert.True(t, b != true, "msg")                  // want "use assert.False"
			assert.True(t, b != true, "msg with arg %d", 42)  // want "use assert.False"
			assert.Truef(t, b != true, "msg")                 // want "use assert.False"
			assert.Truef(t, b != true, "msg with arg %d", 42) // want "use assert.False"

			assert.True(t, true != b)                         // want "use assert.False"
			assert.True(t, true != b, "msg")                  // want "use assert.False"
			assert.True(t, true != b, "msg with arg %d", 42)  // want "use assert.False"
			assert.Truef(t, true != b, "msg")                 // want "use assert.False"
			assert.Truef(t, true != b, "msg with arg %d", 42) // want "use assert.False"

			assert.False(t, b != false)                         // want "need to simplify the check"
			assert.False(t, b != false, "msg")                  // want "need to simplify the check"
			assert.False(t, b != false, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Falsef(t, b != false, "msg")                 // want "need to simplify the check"
			assert.Falsef(t, b != false, "msg with arg %d", 42) // want "need to simplify the check"

			assert.False(t, false != b)                         // want "need to simplify the check"
			assert.False(t, false != b, "msg")                  // want "need to simplify the check"
			assert.False(t, false != b, "msg with arg %d", 42)  // want "need to simplify the check"
			assert.Falsef(t, false != b, "msg")                 // want "need to simplify the check"
			assert.Falsef(t, false != b, "msg with arg %d", 42) // want "need to simplify the check"

			assert.True(t, !b)                         // want "use assert.False"
			assert.True(t, !b, "msg")                  // want "use assert.False"
			assert.True(t, !b, "msg with arg %d", 42)  // want "use assert.False"
			assert.Truef(t, !b, "msg")                 // want "use assert.Falsef"
			assert.Truef(t, !b, "msg with arg %d", 42) // want "use assert.Falsef"
		}

		// Valid asserts.

		assert.True(t, b)
		assert.True(t, b, "msg")
		assert.True(t, b, "msg with arg %d", 42)
		assert.Truef(t, b, "msg")
		assert.Truef(t, b, "msg with arg %d", 42)

		assert.False(t, b)
		assert.False(t, b, "msg")
		assert.False(t, b, "msg with arg %d", 42)
		assert.Falsef(t, b, "msg")
		assert.Falsef(t, b, "msg with arg %d", 42)
	})

	t.Run("require", func(t *testing.T) {
		{
			require.Equal(t, b, true)                         // want "use require.True"
			require.Equal(t, b, true, "msg")                  // want "use require.True"
			require.Equal(t, b, true, "msg with arg %d", 42)  // want "use require.True"
			require.Equalf(t, b, true, "msg")                 // want "use require.Truef"
			require.Equalf(t, b, true, "msg with arg %d", 42) // want "use require.Truef"

			require.Equal(t, true, b)                         // want "use require.True"
			require.Equal(t, true, b, "msg")                  // want "use require.True"
			require.Equal(t, true, b, "msg with arg %d", 42)  // want "use require.True"
			require.Equalf(t, true, b, "msg")                 // want "use require.Truef"
			require.Equalf(t, true, b, "msg with arg %d", 42) // want "use require.Truef"

			require.NotEqual(t, b, false)                         // want "use require.True"
			require.NotEqual(t, b, false, "msg")                  // want "use require.True"
			require.NotEqual(t, b, false, "msg with arg %d", 42)  // want "use require.True"
			require.NotEqualf(t, b, false, "msg")                 // want "use require.Truef"
			require.NotEqualf(t, b, false, "msg with arg %d", 42) // want "use require.Truef"

			require.NotEqual(t, false, b)                         // want "use require.True"
			require.NotEqual(t, false, b, "msg")                  // want "use require.True"
			require.NotEqual(t, false, b, "msg with arg %d", 42)  // want "use require.True"
			require.NotEqualf(t, false, b, "msg")                 // want "use require.Truef"
			require.NotEqualf(t, false, b, "msg with arg %d", 42) // want "use require.Truef"

			require.True(t, b == true)                         // want "need to simplify the check"
			require.True(t, b == true, "msg")                  // want "need to simplify the check"
			require.True(t, b == true, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Truef(t, b == true, "msg")                 // want "need to simplify the check"
			require.Truef(t, b == true, "msg with arg %d", 42) // want "need to simplify the check"

			require.True(t, true == b)                         // want "need to simplify the check"
			require.True(t, true == b, "msg")                  // want "need to simplify the check"
			require.True(t, true == b, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Truef(t, true == b, "msg")                 // want "need to simplify the check"
			require.Truef(t, true == b, "msg with arg %d", 42) // want "need to simplify the check"

			require.False(t, b == false)                         // want "use require.True"
			require.False(t, b == false, "msg")                  // want "use require.True"
			require.False(t, b == false, "msg with arg %d", 42)  // want "use require.True"
			require.Falsef(t, b == false, "msg")                 // want "use require.Truef"
			require.Falsef(t, b == false, "msg with arg %d", 42) // want "use require.Truef"

			require.False(t, false == b)                         // want "use require.True"
			require.False(t, false == b, "msg")                  // want "use require.True"
			require.False(t, false == b, "msg with arg %d", 42)  // want "use require.True"
			require.Falsef(t, false == b, "msg")                 // want "use require.Truef"
			require.Falsef(t, false == b, "msg with arg %d", 42) // want "use require.Truef"

			require.False(t, b != true)                         // want "use require.True"
			require.False(t, b != true, "msg")                  // want "use require.True"
			require.False(t, b != true, "msg with arg %d", 42)  // want "use require.True"
			require.Falsef(t, b != true, "msg")                 // want "use require.Truef"
			require.Falsef(t, b != true, "msg with arg %d", 42) // want "use require.Truef"

			require.False(t, true != b)                         // want "use require.True"
			require.False(t, true != b, "msg")                  // want "use require.True"
			require.False(t, true != b, "msg with arg %d", 42)  // want "use require.True"
			require.Falsef(t, true != b, "msg")                 // want "use require.Truef"
			require.Falsef(t, true != b, "msg with arg %d", 42) // want "use require.Truef"

			require.True(t, b != false)                         // want "need to simplify the check"
			require.True(t, b != false, "msg")                  // want "need to simplify the check"
			require.True(t, b != false, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Truef(t, b != false, "msg")                 // want "need to simplify the check"
			require.Truef(t, b != false, "msg with arg %d", 42) // want "need to simplify the check"

			require.True(t, false != b)                         // want "need to simplify the check"
			require.True(t, false != b, "msg")                  // want "need to simplify the check"
			require.True(t, false != b, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Truef(t, false != b, "msg")                 // want "need to simplify the check"
			require.Truef(t, false != b, "msg with arg %d", 42) // want "need to simplify the check"

			require.False(t, !b)                         // want "use require.True"
			require.False(t, !b, "msg")                  // want "use require.True"
			require.False(t, !b, "msg with arg %d", 42)  // want "use require.True"
			require.Falsef(t, !b, "msg")                 // want "use require.Truef"
			require.Falsef(t, !b, "msg with arg %d", 42) // want "use require.Truef"
		}

		{
			require.Equal(t, b, false)                         // want "use require.False"
			require.Equal(t, b, false, "msg")                  // want "use require.False"
			require.Equal(t, b, false, "msg with arg %d", 42)  // want "use require.False"
			require.Equalf(t, b, false, "msg")                 // want "use require.Falsef"
			require.Equalf(t, b, false, "msg with arg %d", 42) // want "use require.Falsef"

			require.Equal(t, false, b)                         // want "use require.False"
			require.Equal(t, false, b, "msg")                  // want "use require.False"
			require.Equal(t, false, b, "msg with arg %d", 42)  // want "use require.False"
			require.Equalf(t, false, b, "msg")                 // want "use require.Falsef"
			require.Equalf(t, false, b, "msg with arg %d", 42) // want "use require.Falsef"

			require.NotEqual(t, b, true)                         // want "use require.False"
			require.NotEqual(t, b, true, "msg")                  // want "use require.False"
			require.NotEqual(t, b, true, "msg with arg %d", 42)  // want "use require.False"
			require.NotEqualf(t, b, true, "msg")                 // want "use require.Falsef"
			require.NotEqualf(t, b, true, "msg with arg %d", 42) // want "use require.Falsef"

			require.NotEqual(t, true, b)                         // want "use require.False"
			require.NotEqual(t, true, b, "msg")                  // want "use require.False"
			require.NotEqual(t, true, b, "msg with arg %d", 42)  // want "use require.False"
			require.NotEqualf(t, true, b, "msg")                 // want "use require.Falsef"
			require.NotEqualf(t, true, b, "msg with arg %d", 42) // want "use require.Falsef"

			require.False(t, b == true)                         // want "need to simplify the check"
			require.False(t, b == true, "msg")                  // want "need to simplify the check"
			require.False(t, b == true, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Falsef(t, b == true, "msg")                 // want "need to simplify the check"
			require.Falsef(t, b == true, "msg with arg %d", 42) // want "need to simplify the check"

			require.False(t, true == b)                         // want "need to simplify the check"
			require.False(t, true == b, "msg")                  // want "need to simplify the check"
			require.False(t, true == b, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Falsef(t, true == b, "msg")                 // want "need to simplify the check"
			require.Falsef(t, true == b, "msg with arg %d", 42) // want "need to simplify the check"

			require.True(t, b == false)                         // want "use require.False"
			require.True(t, b == false, "msg")                  // want "use require.False"
			require.True(t, b == false, "msg with arg %d", 42)  // want "use require.False"
			require.Truef(t, b == false, "msg")                 // want "use require.Falsef"
			require.Truef(t, b == false, "msg with arg %d", 42) // want "use require.Falsef"

			require.True(t, false == b)                         // want "use require.False"
			require.True(t, false == b, "msg")                  // want "use require.False"
			require.True(t, false == b, "msg with arg %d", 42)  // want "use require.False"
			require.Truef(t, false == b, "msg")                 // want "use require.Falsef"
			require.Truef(t, false == b, "msg with arg %d", 42) // want "use require.Falsef"

			require.True(t, b != true)                         // want "use require.False"
			require.True(t, b != true, "msg")                  // want "use require.False"
			require.True(t, b != true, "msg with arg %d", 42)  // want "use require.False"
			require.Truef(t, b != true, "msg")                 // want "use require.False"
			require.Truef(t, b != true, "msg with arg %d", 42) // want "use require.False"

			require.True(t, true != b)                         // want "use require.False"
			require.True(t, true != b, "msg")                  // want "use require.False"
			require.True(t, true != b, "msg with arg %d", 42)  // want "use require.False"
			require.Truef(t, true != b, "msg")                 // want "use require.False"
			require.Truef(t, true != b, "msg with arg %d", 42) // want "use require.False"

			require.False(t, b != false)                         // want "need to simplify the check"
			require.False(t, b != false, "msg")                  // want "need to simplify the check"
			require.False(t, b != false, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Falsef(t, b != false, "msg")                 // want "need to simplify the check"
			require.Falsef(t, b != false, "msg with arg %d", 42) // want "need to simplify the check"

			require.False(t, false != b)                         // want "need to simplify the check"
			require.False(t, false != b, "msg")                  // want "need to simplify the check"
			require.False(t, false != b, "msg with arg %d", 42)  // want "need to simplify the check"
			require.Falsef(t, false != b, "msg")                 // want "need to simplify the check"
			require.Falsef(t, false != b, "msg with arg %d", 42) // want "need to simplify the check"

			require.True(t, !b)                         // want "use require.False"
			require.True(t, !b, "msg")                  // want "use require.False"
			require.True(t, !b, "msg with arg %d", 42)  // want "use require.False"
			require.Truef(t, !b, "msg")                 // want "use require.Falsef"
			require.Truef(t, !b, "msg with arg %d", 42) // want "use require.Falsef"
		}

		// Valid requires.

		require.True(t, b)
		require.True(t, b, "msg")
		require.True(t, b, "msg with arg %d", 42)
		require.Truef(t, b, "msg")
		require.Truef(t, b, "msg with arg %d", 42)

		require.False(t, b)
		require.False(t, b, "msg")
		require.False(t, b, "msg with arg %d", 42)
		require.Falsef(t, b, "msg")
		require.Falsef(t, b, "msg with arg %d", 42)
	})
}
