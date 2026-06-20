package requireerrorissue287

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type checker struct {
	foo  []any
	bar  []any
	name string
}

func (c *checker) cond(collect *assert.CollectT) {
	assert.NoError(collect, nil) // want "require-error: for error assertions use require"

	if c.foo == nil && c.bar == nil {
		collect.Errorf("precondition not satisfied")
		collect.Errorf("precondition not satisfied: %v") // want "formatter: collect.Errorf format %v reads arg #1, but call has 0 args"
	}
	require.NotEmpty(collect, c.name)
}
