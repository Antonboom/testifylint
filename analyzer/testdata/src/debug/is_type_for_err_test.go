package debug

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTypeForError(t *testing.T) {
	var err error
	assert.IsType(t, nil, err)
	assert.IsType(t, error(nil), err)

	err = (*http.MaxBytesError)(nil)
	assert.IsType(t, (*http.MaxBytesError)(nil), err)
	assert.NoError(t, err)

	err = &http.MaxBytesError{}
	assert.IsType(t, &http.MaxBytesError{}, err)
	assert.IsType(t, err, &http.MaxBytesError{})
	assert.IsType(t, new(http.MaxBytesError), err)
	assert.ErrorAs(t, err, new(*http.MaxBytesError))

	assert.IsType(t, io.EOF, err)
}
