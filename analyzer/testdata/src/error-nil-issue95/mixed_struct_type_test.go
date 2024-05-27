package errornilissue95

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// APIResponse can be an error or not.
type APIResponse struct {
	Status   int
	Data     string
	ErrorMsg string
}

func (a APIResponse) Error() string {
	return a.ErrorMsg
}

func Update(a string) (*APIResponse, error) {
	if a == "a" {
		return &APIResponse{Status: 200, Data: "fake"}, nil
	}

	return nil, &APIResponse{Status: 500, ErrorMsg: "Oops"}
}

func TestName(t *testing.T) {
	resp, err := Update("b")
	require.Error(t, err, new(*APIResponse))
	assert.Nil(t, resp)
}
