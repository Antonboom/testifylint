package requireerrorskiplogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNegatedIfReturn tests that `if !assert.Error { return }` patterns with
// error assertions are detected and fixed by require-error.
func TestNegatedIfReturn(t *testing.T) {
	var err error

	// Valid: simple pattern – should be reported and fixed.
	if !assert.NoError(t, err) { // want "require-error: for error assertions use require"
		return
	}

	// Valid: f-variant – should be reported and fixed.
	if !assert.NoErrorf(t, err, "msg") { // want "require-error: for error assertions use require"
		return
	}

	// Valid: other error-assertion variants.
	if !assert.Error(t, err) { // want "require-error: for error assertions use require"
		return
	}

	// Skip: complex body (not just return).
	if !assert.NoError(t, err) {
		_ = "something"
		return
	}

	// Skip: positive assertion in if condition (not negated).
	if assert.NoError(t, err) {
		_ = err
	}

	// Skip: has else clause.
	if !assert.NoError(t, err) {
		return
	} else {
		_ = err
	}

	// Skip: !assert.Contains is not an error assertion – not in require-error scope
	// (handled by negated-assert checker instead).
	if !assert.Contains(t, "str", "s") {
		return
	}
}

// TestNegatedIfContinue tests that `if !assert.Error { continue }` patterns
// inside loops are detected and fixed.
func TestNegatedIfContinue(t *testing.T) {
	var err error

	for i := 0; i < 10; i++ {
		if !assert.NoError(t, err) { // want "require-error: for error assertions use require"
			continue
		}
	}
}

// TestNegatedIfReturnLast tests the case where the negated-if pattern is the
// last assertion in the function (correctly reported, not skipped by "last assertion" logic).
func TestNegatedIfReturnLast(t *testing.T) {
	var err error
	if !assert.NoError(t, err) { // want "require-error: for error assertions use require"
		return
	}
}

// TestNegatedIfCompound tests the compound || pattern where ALL conditions are
// negated error assertions.
func TestNegatedIfCompound(t *testing.T) {
	var err, err2 error

	if !assert.NoError(t, err) || !assert.NoError(t, err2) { // want "require-error: for error assertions use require"
		return
	}

	// Skip: not all conditions are negated error assertions (Contains is not an error assertion).
	if !assert.NoError(t, err) || !assert.Contains(t, "str", "s") {
		return
	}

	// Skip: && is not supported (different semantics).
	if !assert.NoError(t, err) && !assert.Error(t, err) {
		return
	}
}

// TestNegatedIfElseChain tests that else-if patterns are not broken by fixes.
func TestNegatedIfElseChain(t *testing.T) {
	var err error

	// Skip: the outer if has an else clause.
	if !assert.NoError(t, err) {
		return
	} else if !assert.Error(t, err) {
		return
	}
}
