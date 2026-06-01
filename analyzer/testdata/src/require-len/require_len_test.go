package requirelen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLenGuard(t *testing.T) {
	arr := []int{0, 1}

	assert.Len(t, arr, 2)      // want "require-len: for length assertions guarding index access use require"
	assert.Positive(t, arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardf(t *testing.T) {
	arr := []int{0, 1}

	assert.Lenf(t, arr, 2, "msg") // want "require-len: for length assertions guarding index access use require"
	assert.Positive(t, arr[1])    // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardRequire(t *testing.T) {
	arr := []int{0, 1}

	require.Len(t, arr, 2)
	assert.Positive(t, arr[1])
}

func TestLenGuardIndexOutOfCheckedRange(t *testing.T) {
	arr := []int{0, 1}

	assert.Len(t, arr, 2)
	assert.Positive(t, arr[2]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardInsertedFromIndexAccess(t *testing.T) {
	arr := []int{0, 1}

	assert.Positive(t, arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardInsertedUsesGreatestIndex(t *testing.T) {
	arr := []int{0, 1, 2}

	assert.Equal(t, arr[0]+arr[2], 2) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardInsertedSingleIndexUsesNotEmpty(t *testing.T) {
	arr := []int{0}

	assert.Positive(t, arr[0]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardDifferentBlock(t *testing.T) {
	arr := []int{0, 1}

	if true {
		assert.Len(t, arr, 2)
	}
	if true {
		consumeInt(arr[1])
	}
}

func TestLenGuardWithoutAutofixForNestedCall(t *testing.T) {
	arr := []int{0, 1}

	consumeBool(assert.Positive(t, arr[1])) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardWithAssertionObject(t *testing.T) {
	arr := []int{0, 1}

	a := assert.New(t)
	a.Positive(arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardWithAssertLenOnly(t *testing.T) {
	arr := []int{0, 1, 2}

	assert.Len(t, arr, 3)      // want "require-len: for length assertions guarding index access use require"
	assert.Positive(t, arr[2]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardWithInsufficientRequireLen(t *testing.T) {
	arr := []int{0, 1, 2}

	require.Len(t, arr, 1)
	assert.Positive(t, arr[2]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardExpectedActualMapBothSides(t *testing.T) {
	var want, res map[int]string

	assert.Equal(t, want[0], res[0]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardExpectedActualMapMissingOneSide(t *testing.T) {
	var want, res map[int]string

	require.Contains(t, want, 1)
	assert.Equal(t, want[1], res[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardExpectedActualSliceIndexOne(t *testing.T) {
	want := []string{"a", "b"}
	res := []string{"a", "b", "c"}

	assert.Equal(t, want[1], res[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardMustBeAfterErrorCheck(t *testing.T) {
	arr := []int{0, 1}
	err := consumeErr()

	require.GreaterOrEqual(t, len(arr), 2)
	require.NoError(t, err)
	assert.Positive(t, arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestContainsGuardMustBeAfterErrorCheck(t *testing.T) {
	var m map[int]string
	err := consumeErr()

	require.Contains(t, m, 1)
	require.NoError(t, err)
	assert.Equal(t, m[1], "v") // want "require-len: for indexed access use require\\.Len guard"
}

func TestLenGuardMustBeAfterEqualErrCheck(t *testing.T) {
	arr := []int{0, 1}
	err := consumeErr()

	require.GreaterOrEqual(t, len(arr), 2)
	assert.Equal(t, nil, err)
	assert.Positive(t, arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}

func TestContainsGuardMustBeAfterNilErrCheck(t *testing.T) {
	var m map[int]string
	err := consumeErr()

	require.Contains(t, m, 1)
	assert.Nil(t, err)
	assert.Equal(t, m[1], "v") // want "require-len: for indexed access use require\\.Len guard"
}

func consumeInt(_ int) {}

func consumeBool(_ bool) {}

func consumeErr() error { return nil }
