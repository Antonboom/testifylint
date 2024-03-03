package boolcomparecustomtypes_test

import (
	"testing"

	"bool-compare-custom-types/types"
	"github.com/stretchr/testify/assert"
)

type MyBool bool

func TestBoolCompareChecker_CustomTypes(t *testing.T) {
	var b MyBool
	{
		assert.Equal(t, false, b)
		assert.EqualValues(t, false, b)
		assert.Exactly(t, false, b)

		assert.Equal(t, true, b)
		assert.EqualValues(t, true, b)
		assert.Exactly(t, true, b)

		assert.NotEqual(t, false, b)
		assert.NotEqualValues(t, false, b)

		assert.NotEqual(t, true, b)
		assert.NotEqualValues(t, true, b)

		assert.True(t, b == true)
		assert.True(t, b != false)
		assert.True(t, b == false)
		assert.True(t, b != true)

		assert.False(t, b == true)
		assert.False(t, b != false)
		assert.False(t, b == false)
		assert.False(t, b != true)
	}

	var extB types.Bool
	{
		assert.Equal(t, false, extB)
		assert.EqualValues(t, false, extB)
		assert.Exactly(t, false, extB)

		assert.Equal(t, true, extB)
		assert.EqualValues(t, true, extB)
		assert.Exactly(t, true, extB)

		assert.NotEqual(t, false, extB)
		assert.NotEqualValues(t, false, extB)

		assert.NotEqual(t, true, extB)
		assert.NotEqualValues(t, true, extB)

		assert.True(t, extB == true)
		assert.True(t, extB != false)
		assert.True(t, extB == false)
		assert.True(t, extB != true)

		assert.False(t, extB == true)
		assert.False(t, extB != false)
		assert.False(t, extB == false)
		assert.False(t, extB != true)
	}

	var extSuperB types.SuperBool
	{
		assert.Equal(t, false, extSuperB)
		assert.EqualValues(t, false, extSuperB)
		assert.Exactly(t, false, extSuperB)

		assert.Equal(t, true, extSuperB)
		assert.EqualValues(t, true, extSuperB)
		assert.Exactly(t, true, extSuperB)

		assert.NotEqual(t, false, extSuperB)
		assert.NotEqualValues(t, false, extSuperB)

		assert.NotEqual(t, true, extSuperB)
		assert.NotEqualValues(t, true, extSuperB)

		assert.True(t, extSuperB == true)
		assert.True(t, extSuperB != false)
		assert.True(t, extSuperB == false)
		assert.True(t, extSuperB != true)

		assert.False(t, extSuperB == true)
		assert.False(t, extSuperB != false)
		assert.False(t, extSuperB == false)
		assert.False(t, extSuperB != true)
	}

	// Crazy cases:
	{
		assert.Equal(t, true, types.Bool(extSuperB))
		assert.Equal(t, true, types.SuperBool(b))
		assert.Equal(t, true, bool(types.SuperBool(b))) // want "bool-compare: use assert\\.True"
		assert.True(t, !bool(types.SuperBool(b)))       // want "bool-compare: use assert\\.False"
		assert.False(t, !bool(types.SuperBool(b)))      // want "bool-compare: use assert\\.True"
	}
}
