package requireerrorskiplogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanup(t *testing.T) {
	t.Cleanup(func() {
		assert.Error(t, nil)
		cleanup(t)
	})

	t.Cleanup(func() {
		cleanup(t)
	})

	t.Cleanup(func() {
		assert.Error(t, nil)

		t.Cleanup(func() {
			assert.Error(t, nil)
			cleanup(t)

			t.Cleanup(func() {
				assert.Error(t, nil)
				cleanup(t)
			})
		})
	})
}

func cleanup(t *testing.T) {
	assert.Error(t, nil) // want "require-error: for error assertions use require"
	cleanup(t)
}
