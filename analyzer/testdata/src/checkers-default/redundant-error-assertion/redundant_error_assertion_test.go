package redundanterrorassertion

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedundantErrorAssertion(t *testing.T) {
	t.Parallel()

	var (
		errA   = errors.New("err A")
		errB   = errors.New("err B")
		target *os.PathError
	)

	assObj, reqObj := assert.New(t), require.New(t)

	// Invalid.
	{
		assert.Error(t, errA) // want "redundant-error-assertion: remove assert\\.Error: assert\\.ErrorContains already asserts that the error is not nil"
		assert.ErrorContains(t, errA, "err")
	}

	{
		require.Error(t, errA) // want "redundant-error-assertion: remove require\\.Error: assert\\.ErrorContains already asserts that the error is not nil"
		assert.ErrorContains(t, errA, "err")
	}

	{
		require.Error(t, errA) // want "redundant-error-assertion: remove require\\.Error: require\\.ErrorContains already asserts that the error is not nil"
		require.ErrorContains(t, errA, "err")
	}

	{
		assert.Error(t, errA) // want "redundant-error-assertion: remove assert\\.Error: require\\.ErrorContains already asserts that the error is not nil"
		require.ErrorContains(t, errA, "err")
	}

	{
		assert.EqualError(t, errA, "err A")
		assert.Error(t, errA) // want "redundant-error-assertion: remove assert\\.Error: assert\\.EqualError already asserts that the error is not nil"
	}

	{
		require.EqualError(t, errA, "err A")
		require.Error(t, errA) // want "redundant-error-assertion: remove require\\.Error: require\\.EqualError already asserts that the error is not nil"
	}

	{
		assObj.Error(errA) // want "redundant-error-assertion: remove assObj\\.Error: assObj\\.ErrorAs already asserts that the error is not nil"
		assObj.ErrorAs(errA, &target)
	}

	{
		reqObj.Error(errA) // want "redundant-error-assertion: remove reqObj\\.Error: reqObj\\.ErrorAs already asserts that the error is not nil"
		reqObj.ErrorAs(errA, &target)
	}

	{
		assert.Errorf(t, errA, "msg") // want "redundant-error-assertion: remove assert\\.Errorf: reqObj\\.ErrorAs already asserts that the error is not nil"
		reqObj.ErrorAs(errA, &target)
	}

	{
		require.Errorf(t, errA, "msg") // want "redundant-error-assertion: remove require\\.Errorf: reqObj\\.ErrorAs already asserts that the error is not nil"
		reqObj.ErrorAs(errA, &target)
	}

	{
		require.Errorf(t, errA, "msg") // want "redundant-error-assertion: remove require\\.Errorf: require\\.ErrorIs already asserts that the error is not nil"
		require.ErrorIs(t, errA, io.EOF)
	}

	// Valid.
	{
		{
			assert.Error(t, errA)
			assert.ErrorContains(t, errB, "err")
		}

		{
			require.Error(t, errA)
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

		{
			errA = errB
			assert.Error(t, errA)

			errA = io.EOF
			assert.ErrorContains(t, errA, "EOF")
		}
	}
}
