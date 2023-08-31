package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFake(t *testing.T) {
	var predicate bool

	assert.Equal(t, true, predicate)
	assert.Equalf(t, true, predicate, "msg with args %d %s", 42, "42")
}
