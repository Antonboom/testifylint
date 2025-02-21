package equalvalues

import (
	"testing"

	v2 "equal-values-different-pkg/v2"
	"github.com/stretchr/testify/assert"
)

func TestEqualValuesChecker(t *testing.T) {
	assert.EqualValues(t, Request{}, Request{}) // want "equal-values: use assert\\.Equal"
	assert.EqualValues(t, Request{}, Response{})
	assert.EqualValues(t, Request{}, v2.Request{}) // THIS IS OK! Using Equal won't work here.
	assert.EqualValues(t, v2.Request{}, v2.Response{})
	assert.EqualValues(t, v2.Response{}, v2.Response{}) // want "equal-values: use assert\\.Equal"
}

type Request struct {
	ID string
}

type Response struct {
	ID string `json:"id"`
}
