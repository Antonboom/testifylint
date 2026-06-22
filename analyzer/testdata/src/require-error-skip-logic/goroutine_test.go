package requireerrorskiplogic

import (
	"sync"
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
			assert.Error(t, nil)
			assert.Error(t, nil)

			go func() {
				assert.Error(t, nil)
				assert.Error(t, nil)
				assert.Error(t, nil)
			}()

			t.Run("", func(t *testing.T) {
				assert.Error(t, nil) // want "require-error: for error assertions use require"
				assert.Error(t, nil)
			})

			go concurrentOp(t)
		}()
	}()

	assert.Error(t, nil) // want "require-error: for error assertions use require"
	assert.Error(t, nil) // want "require-error: for error assertions use require"

	go concurrentOp(t)
}

func TestWaitGroupGo(t *testing.T) {
	var wg sync.WaitGroup

	// wg.Go runs callback in a new goroutine — no require-error warnings expected.
	wg.Go(func() {
		assert.NoError(t, nil)
		assert.Error(t, nil)
	})

	// Indirect callbacks are not supported, consistently with `go callback()`.
	callback := func() {
		assert.Error(t, nil) // want "require-error: for error assertions use require"
		assert.Error(t, nil)
	}
	wg.Go(callback)

	wg.Go(func() {
		assert.NoError(t, nil)

		wg.Go(func() {
			assert.Error(t, nil)
			assert.Error(t, nil)
		})
	})

	(*sync.WaitGroup).Go(&wg, func() {
		assert.NoError(t, nil)
		assert.Error(t, nil)
	})

	var embedded struct{ sync.WaitGroup }
	embedded.Go(func() {
		assert.NoError(t, nil)
		assert.Error(t, nil)
	})

	var custom customGo
	custom.Go(func() {
		assert.Error(t, nil) // want "require-error: for error assertions use require"
		assert.Error(t, nil)
	})

	assert.Error(t, nil) // want "require-error: for error assertions use require"
	assert.Error(t, nil) // want "require-error: for error assertions use require"
	wg.Wait()
}

type customGo struct{}

func (customGo) Go(f func()) { f() }

func concurrentOp(t *testing.T) {
	assert.Error(t, nil) // want "require-error: for error assertions use require"
	assert.Error(t, nil)
}
