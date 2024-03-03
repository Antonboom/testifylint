package boolcomparecustomtypes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type bool int

func TestBoolCompareChecker_BoolOverride(t *testing.T) {
	var mimic bool
	assert.Equal(t, false, mimic)
	assert.Equal(t, false, mimic)
	assert.EqualValues(t, false, mimic)
	assert.Exactly(t, false, mimic)
	assert.NotEqual(t, false, mimic)
	assert.NotEqualValues(t, false, mimic)
}
