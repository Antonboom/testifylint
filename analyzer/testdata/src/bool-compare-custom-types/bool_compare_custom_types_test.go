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
		assert.Equal(t, false, b)       // want "bool-compare: use assert\\.False"
		assert.EqualValues(t, false, b) // want "bool-compare: use assert\\.False"
		assert.Exactly(t, false, b)

		assert.Equal(t, true, b)       // want "bool-compare: use assert\\.True"
		assert.EqualValues(t, true, b) // want "bool-compare: use assert\\.True"
		assert.Exactly(t, true, b)

		assert.NotEqual(t, false, b)       // want "bool-compare: use assert\\.True"
		assert.NotEqualValues(t, false, b) // want "bool-compare: use assert\\.True"

		assert.NotEqual(t, true, b)       // want "bool-compare: use assert\\.False"
		assert.NotEqualValues(t, true, b) // want "bool-compare: use assert\\.False"

		assert.True(t, b == true)  // want "bool-compare: need to simplify the assertion"
		assert.True(t, b != false) // want "bool-compare: need to simplify the assertion"
		assert.True(t, b == false) // want "bool-compare: use assert\\.False"
		assert.True(t, b != true)  // want "bool-compare: use assert\\.False"

		assert.False(t, b == true)  // want "bool-compare: need to simplify the assertion"
		assert.False(t, b != false) // want "bool-compare: need to simplify the assertion"
		assert.False(t, b == false) // want "bool-compare: use assert\\.True"
		assert.False(t, b != true)  // want "bool-compare: use assert\\.True"
	}

	var extB types.Bool
	{
		assert.Equal(t, false, extB)       // want "bool-compare: use assert\\.False"
		assert.EqualValues(t, false, extB) // want "bool-compare: use assert\\.False"
		assert.Exactly(t, false, extB)

		assert.Equal(t, true, extB)       // want "bool-compare: use assert\\.True"
		assert.EqualValues(t, true, extB) // want "bool-compare: use assert\\.True"
		assert.Exactly(t, true, extB)

		assert.NotEqual(t, false, extB)       // want "bool-compare: use assert\\.True"
		assert.NotEqualValues(t, false, extB) // want "bool-compare: use assert\\.True"

		assert.NotEqual(t, true, extB)       // want "bool-compare: use assert\\.False"
		assert.NotEqualValues(t, true, extB) // want "bool-compare: use assert\\.False"

		assert.True(t, extB == true)  // want "bool-compare: need to simplify the assertion"
		assert.True(t, extB != false) // want "bool-compare: need to simplify the assertion"
		assert.True(t, extB == false) // want "bool-compare: use assert\\.False"
		assert.True(t, extB != true)  // want "bool-compare: use assert\\.False"

		assert.False(t, extB == true)  // want "bool-compare: need to simplify the assertion"
		assert.False(t, extB != false) // want "bool-compare: need to simplify the assertion"
		assert.False(t, extB == false) // want "bool-compare: use assert\\.True"
		assert.False(t, extB != true)  // want "bool-compare: use assert\\.True"
	}

	var extSuperB types.SuperBool
	{
		assert.Equal(t, false, extSuperB)       // want "bool-compare: use assert\\.False"
		assert.EqualValues(t, false, extSuperB) // want "bool-compare: use assert\\.False"
		assert.Exactly(t, false, extSuperB)

		assert.Equal(t, true, extSuperB)       // want "bool-compare: use assert\\.True"
		assert.EqualValues(t, true, extSuperB) // want "bool-compare: use assert\\.True"
		assert.Exactly(t, true, extSuperB)

		assert.NotEqual(t, false, extSuperB)       // want "bool-compare: use assert\\.True"
		assert.NotEqualValues(t, false, extSuperB) // want "bool-compare: use assert\\.True"

		assert.NotEqual(t, true, extSuperB)       // want "bool-compare: use assert\\.False"
		assert.NotEqualValues(t, true, extSuperB) // want "bool-compare: use assert\\.False"

		assert.True(t, extSuperB == true)  // want "bool-compare: need to simplify the assertion"
		assert.True(t, extSuperB != false) // want "bool-compare: need to simplify the assertion"
		assert.True(t, extSuperB == false) // want "bool-compare: use assert\\.False"
		assert.True(t, extSuperB != true)  // want "bool-compare: use assert\\.False"

		assert.False(t, extSuperB == true)  // want "bool-compare: need to simplify the assertion"
		assert.False(t, extSuperB != false) // want "bool-compare: need to simplify the assertion"
		assert.False(t, extSuperB == false) // want "bool-compare: use assert\\.True"
		assert.False(t, extSuperB != true)  // want "bool-compare: use assert\\.True"
	}

	// Crazy cases:
	{
		assert.Equal(t, true, types.Bool(extSuperB))    // want "bool-compare: use assert\\.True"
		assert.Equal(t, true, types.SuperBool(b))       // want "bool-compare: use assert\\.True"
		assert.Equal(t, true, bool(types.SuperBool(b))) // want "bool-compare: use assert\\.True"
		assert.True(t, !bool(types.SuperBool(b)))       // want "bool-compare: use assert\\.False"
		assert.False(t, !bool(types.SuperBool(b)))      // want "bool-compare: use assert\\.True"
	}
}

func TestBoolCompareChecker_CustomTypes_Format(t *testing.T) {
	var predicate MyBool
	assert.Equal(t, true, predicate)                                   // want "bool-compare: use assert\\.True"
	assert.Equal(t, true, predicate, "msg")                            // want "bool-compare: use assert\\.True"
	assert.Equal(t, true, predicate, "msg with arg %d", 42)            // want "bool-compare: use assert\\.True"
	assert.Equal(t, true, predicate, "msg with args %d %s", 42, "42")  // want "bool-compare: use assert\\.True"
	assert.Equalf(t, true, predicate, "msg")                           // want "bool-compare: use assert\\.Truef"
	assert.Equalf(t, true, predicate, "msg with arg %d", 42)           // want "bool-compare: use assert\\.Truef"
	assert.Equalf(t, true, predicate, "msg with args %d %s", 42, "42") // want "bool-compare: use assert\\.Truef"
}
