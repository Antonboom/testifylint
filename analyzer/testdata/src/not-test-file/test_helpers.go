package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertTrue(t *testing.T, predicate bool) {
	t.Helper()

	assert.Equal(t, true, predicate)                                   // want "bool-compare: use assert\\.True"
	assert.Equalf(t, true, predicate, "msg with args %d %s", 42, "42") // want "bool-compare: use assert\\.True"
}
