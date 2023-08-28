package pkgalias

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"
)

func TestPkgAlias(t *testing.T) {
	var predicate bool

	a.Equal(t, true, predicate) // want "bool-compare: use a\\.True"
	r.Equal(t, true, predicate) // want "bool-compare: use r\\.True"
}
