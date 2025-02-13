package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v2 "testdata/debug/v2"
)

type S []string

func TestEqualValues(t *testing.T) {
	var f1 func()
	var f2 func() bool
	var f3 func()

	assert.Equal(t, []string{"1"}, S{"1"})       // Not equal.
	assert.EqualValues(t, []string{"1"}, S{"1"}) // Equal.

	assert.Equal(t, f1, f2)                     // Invalid operation: cannot take func type as argument.
	assert.EqualValues(t, f1, f2)               // Not equal.
	assert.EqualValues(t, f1, f3)               // Equal.
	assert.EqualValues(t, func() {}, func() {}) // Not equal.

	assert.Equal(t, []v2.HPAScalingPolicy{}, []HPAScalingPolicy{})       // Not equal.
	assert.EqualValues(t, []v2.HPAScalingPolicy{}, []HPAScalingPolicy{}) // Not equal.

	assert.Equal(t, Taint{}, v2.Taint{})       // Not equal.
	assert.EqualValues(t, Taint{}, v2.Taint{}) // Equal.

	assert.Equal(t, []Taint{}, []v2.Taint{})       // Not equal.
	assert.EqualValues(t, []Taint{}, []v2.Taint{}) // Not equal.

	/* Panic, see https://github.com/stretchr/testify/issues/1699
	a := []int{1, 2}
	b := (*[3]int)(nil)
	assert.EqualValues(t, a, b)
	*/
}

// HPAScalingPolicy is a single policy which must hold true for a specified past interval.
type HPAScalingPolicy struct {
	// value contains the amount of change which is permitted by the policy.
	// It must be greater than zero
	Value int32 `json:"value" protobuf:"varint,2,opt,name=value"`

	// periodSeconds specifies the window of time for which the policy should hold true.
	// PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
	PeriodSeconds int32 `json:"periodSeconds" protobuf:"varint,3,opt,name=periodSeconds"`
}

// The node this Taint is attached to has the "effect" on
// any pod that does not tolerate the Taint.
type Taint struct {
	// Required. The taint key to be applied to a node.
	Key string
	// The taint value corresponding to the taint key.
	// +optional
	Value string
}
