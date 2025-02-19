package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLenInterface(t *testing.T) {
	const n = 10
	a := newArr(n)
	assert.Equal(t, n, a.Len())
	assert.Len(t, a, n) // Error: "{[         ]}" could not be applied builtin len()
}

type arr struct {
	v []string
}

func newArr(n int) arr {
	return arr{v: make([]string, n)}
}

func (a arr) Len() int {
	return len(a.v)
}
