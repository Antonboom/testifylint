package basic

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestErrorInsteadOfErrorIs(t *testing.T) {
	err := operation()

	assert.Error(t, err, io.EOF)                                    // want "invalid usage of assert.Error, use assert.ErrorIs instead"
	assert.Error(t, err, new(os.PathError))                         // want "invalid usage of assert.Error, use assert.ErrorIs instead"
	assert.Error(t, err, errors.New("sky is falling"))              // want "invalid usage of assert.Error, use assert.ErrorIs instead"
	assert.Error(t, err, fmt.Errorf("sky is falling %d times", 10)) // want "invalid usage of assert.Error, use assert.ErrorIs instead"

	require.Error(t, err, io.EOF)                                    // want "invalid usage of require.Error, use require.ErrorIs instead"
	require.Error(t, err, new(os.PathError))                         // want "invalid usage of require.Error, use require.ErrorIs instead"
	require.Error(t, err, errors.New("sky is falling"))              // want "invalid usage of require.Error, use require.ErrorIs instead"
	require.Error(t, err, fmt.Errorf("sky is falling %d times", 10)) // want "invalid usage of require.Error, use require.ErrorIs instead"

	// f-functions
}
