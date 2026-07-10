package requirelen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLenGuardWithoutRequireImport(t *testing.T) {
	arr := []int{0, 1}

	assert.Positive(t, arr[1]) // want "require-len: for indexed access use require\\.Len guard"
}
