package redundantassert

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedundantAssert(t *testing.T) {
	t.Parallel()

	var (
		errA   = errors.New("err A")
		errB   = errors.New("err B")
		target *os.PathError
	)

	assObj := assert.New(t)

	// Invalid.
	{
		assert.Error(t, errA) // want "redundant-assert: remove assert\\.Error: assert\\.ErrorContains already asserts that the error is not nil"
		assert.ErrorContains(t, errA, "err")

		assert.EqualError(t, errA, "err A")
		assert.Error(t, errA) // want "redundant-assert: remove assert\\.Error: assert\\.EqualError already asserts that the error is not nil"

		assObj.Error(errA) // want "redundant-assert: remove assert\\.Error: assObj\\.ErrorAs already asserts that the error is not nil"
		assObj.ErrorAs(errA, &target)

		assert.Errorf(t, errA, "msg") // want "redundant-assert: remove assert\\.Errorf: assObj\\.ErrorAs already asserts that the error is not nil"
		require.ErrorIs(t, errA, io.EOF)
	}

	// Valid.
	{
		{
			assert.Error(t, errA)
			assert.ErrorContains(t, errB, "err")
		}

		{
			assert.Error(t, errA)
		}

		if assert.Error(t, errA) {
			assert.ErrorContains(t, errA, "err")
		}

		{
			{
				assert.Error(t, errA)
			}
			assert.ErrorContains(t, errA, "err")
		}

		require.Error(t, errA)
		require.ErrorContains(t, errA, "err")
	}
}
