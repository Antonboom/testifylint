package basic

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// "use NoError instead of Nil for error check"

func TestNotNilInsteadOfError(t *testing.T) {
	// f-functions
}

func TestNilInsteadOfNoError(t *testing.T) {
	err := operation()

	assert.Nil(t, err)                         // want "for a better message use assert.NoError instead"
	assert.Nilf(t, err, "msg")                 // want "for a better message use assert.NoErrorf instead"
	assert.Nilf(t, err, "msg with arg %d", 42) // want "for a better message use assert.NoErrorf instead"

	require.Nil(t, err)                         // want "for a better message use require.NoError instead"
	require.Nilf(t, err, "msg")                 // want "for a better message use require.NoErrorf instead"
	require.Nilf(t, err, "msg with arg %d", 42) // want "for a better message use require.NoErrorf instead"

	assert.NoError(t, err)
	assert.NoErrorf(t, err, "msg")
	assert.NoErrorf(t, err, "msg with arg %d", 42)

	require.NoError(t, err)
	require.NoErrorf(t, err, "msg")
	require.NoErrorf(t, err, "msg with arg %d", 42)
}

func operation() error {
	return errors.New("invalid")
}
