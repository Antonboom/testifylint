package notstdfuncs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	myerrors "not-std-funcs/errors"
)

func TestCustomStuff(t *testing.T) {
	var err error
	errSentinel := errors.New("unexpected")
	var scores []float64

	assert.Equal(t, nil, scores) // want "nil-compare: use assert\\.Nil"

	assert.True(t, myerrors.Is(err, 2))
	assert.False(t, Is(err, errSentinel))
	assert.True(t, As(err, &err))
	assert.Equal(t, 3, len(scores))

	require.True(t, myerrors.Is(err, 3))
	require.False(t, Is(err, errSentinel))
	require.True(t, As(err, &err))
	require.Equal(t, 4, len(scores))
}

func Is(err, target error) bool {
	return true
}

func As(err error, target any) bool {
	return false
}

func len[T any](arr T) int {
	return 3
}
