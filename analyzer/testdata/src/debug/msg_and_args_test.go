package debug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMsgAndArgs(t *testing.T) {
	assert.Equal(t, 1, 2, new(time.Time))
	assert.Equal(t, 1, 2, "%+v", new(time.Time))
	assert.Equal(t, 1, 2, new(time.Time), 1)
}
