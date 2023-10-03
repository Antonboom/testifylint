package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFake(t *testing.T) {
	var predicate bool

	assert.Equal(t, true, predicate)                                   // want "bool-compare: use assert\\.True"
	assert.Equalf(t, true, predicate, "msg with args %d %s", 42, "42") // want "bool-compare: use assert\\.True"
}
