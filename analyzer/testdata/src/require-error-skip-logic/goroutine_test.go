package requireerrorskiplogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	go func() {
		assert.NoError(t, nil)
	}()

	go func() {
		assert.NoError(t, nil)

		go func() {
			assert.NoError(t, nil)
			assert.NoError(t, nil)

			go func() {
				assert.NoError(t, nil)
				assert.NoError(t, nil)
				assert.NoError(t, nil)
			}()

			t.Run("", func(t *testing.T) {
				assert.NoError(t, nil) // want "require-error: for error assertions use require"
				assert.NoError(t, nil)
			})

			go concurrentOp(t)
		}()
	}()

	assert.NoError(t, nil) // want "require-error: for error assertions use require"
	assert.NoError(t, nil) // want "require-error: for error assertions use require"

	go concurrentOp(t)
}

func concurrentOp(t *testing.T) {
	assert.NoError(t, nil) // want "require-error: for error assertions use require"
	assert.NoError(t, nil)
}
