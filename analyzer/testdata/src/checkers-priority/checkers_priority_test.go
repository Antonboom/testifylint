package checkerspriority

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckersPriority(t *testing.T) {
	var f float64
	var b bool

	// `empty` > `len` > `expected-actual`
	assert.Equal(t, len([]int{}), 0) // want "empty: use assert\\.Empty"
	assert.Equal(t, len([]int{}), 3) // want "len: use assert\\.Len"

	// `float-compare` > `bool-compare` > `compares` > `expected-actual`
	require.True(t, 42.42 == f) // want "float-compare: use require\\.InEpsilon \\(or InDelta\\)"
	require.Equal(t, b, true)   // want "bool-compare: use require\\.True"
	require.True(t, b == true)  // want "bool-compare: need to simplify the assertion"
	require.False(t, b == b)    // want "compares: use require\\.NotEqual"
}
