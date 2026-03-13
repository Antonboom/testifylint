package requireerrorskiplogic

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type checkerIssue287 struct {
	foo  []any
	bar  []any
	name string
}

// collect.Errorf is not a testify assertion; it should not be reported.
func (c *checkerIssue287) cond(collect *assert.CollectT) {
	if c.foo == nil && c.bar == nil {
		collect.Errorf("precondition not satisfied")
	}
	require.NotEmpty(collect, c.name)
}
