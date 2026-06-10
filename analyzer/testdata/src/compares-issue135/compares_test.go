package comparesissue135

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression tests for https://github.com/Antonboom/testifylint/issues/135 and
// https://github.com/Antonboom/testifylint/issues/205:
// When an untyped constant appears in a binary expression, Go contextually types
// it to match the other operand. However, when passed as interface{} to a testify
// function the constant reverts to its default type (e.g. int), causing a runtime
// "Elements should be the same type" failure if the other operand has a different
// type (e.g. int32). The fix must cast the constant to the contextual type.
func TestComparesIssue135(t *testing.T) {
	var val int32

	// Untyped constants: math.MaxInt16 is untyped int (default int), but in the
	// binary expression it is contextually typed to int32. The fix must emit an
	// explicit int32(math.MaxInt16) cast so both arguments have the same runtime type.
	assert.True(t, val <= math.MaxInt16)  // want "compares: use assert\\.LessOrEqual"
	assert.True(t, val < math.MaxInt16)   // want "compares: use assert\\.Less"
	assert.True(t, val >= math.MinInt16)  // want "compares: use assert\\.GreaterOrEqual"
	assert.True(t, val > math.MinInt16)   // want "compares: use assert\\.Greater"
	assert.True(t, val == math.MaxInt16)  // want "compares: use assert\\.Equal"
	assert.True(t, val != math.MaxInt16)  // want "compares: use assert\\.NotEqual"
	assert.False(t, val <= math.MaxInt16) // want "compares: use assert\\.Greater"
	assert.False(t, val < math.MaxInt16)  // want "compares: use assert\\.GreaterOrEqual"
	assert.False(t, val >= math.MinInt16) // want "compares: use assert\\.Less"
	assert.False(t, val > math.MinInt16)  // want "compares: use assert\\.LessOrEqual"
	assert.False(t, val == math.MaxInt16) // want "compares: use assert\\.NotEqual"
	assert.False(t, val != math.MaxInt16) // want "compares: use assert\\.Equal"

	// Same-type comparisons: both operands are int32, no cast required.
	var other int32
	assert.True(t, val <= other)  // want "compares: use assert\\.LessOrEqual"
	assert.True(t, val < other)   // want "compares: use assert\\.Less"
	assert.True(t, val >= other)  // want "compares: use assert\\.GreaterOrEqual"
	assert.True(t, val > other)   // want "compares: use assert\\.Greater"
	assert.True(t, val == other)  // want "compares: use assert\\.Equal"
	assert.True(t, val != other)  // want "compares: use assert\\.NotEqual"
	assert.False(t, val <= other) // want "compares: use assert\\.Greater"
	assert.False(t, val < other)  // want "compares: use assert\\.GreaterOrEqual"
	assert.False(t, val >= other) // want "compares: use assert\\.Less"
	assert.False(t, val > other)  // want "compares: use assert\\.LessOrEqual"
	assert.False(t, val == other) // want "compares: use assert\\.NotEqual"
	assert.False(t, val != other) // want "compares: use assert\\.Equal"
}
