package formatter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatterChecker(t *testing.T) {
	assert.Equal(t, 1, 2, new(time.Time))        // want "formatter: do not use non-string value as first element of msgAndArgs"
	assert.Equal(t, 1, 2, "%+v", new(time.Time)) // want "formatter: use assert\\.Equalf"
	assert.Equalf(t, 1, 2, "%+v", new(time.Time))
}
