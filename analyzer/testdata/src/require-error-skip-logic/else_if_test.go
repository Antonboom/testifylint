package requireerrorskiplogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplexCondition(t *testing.T) {
	testCases := []struct {
		testName       string
		someValue      any
		someOtherValue any
		expectedError  error
		expectedValue  any
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := operationWithResult()
			if tc.someValue == nil && tc.someOtherValue == nil {
				assert.Nil(t, result)
				assert.NoError(t, err)
			} else if tc.someOtherValue != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedValue, result)
			}
		})
	}
}

func TestCrazyCondition(t *testing.T) {
	testCases := []struct {
		testName       string
		someValue      any
		someOtherValue any
		expectedError  error
		expectedValue  any
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := operationWithResult()
			if tc.someValue == nil && tc.someOtherValue == nil {
				assert.Nil(t, result)
				assert.NoError(t, err)
			} else if tc.someOtherValue != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else if tc.someOtherValue == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedValue, result)
			} else if tc.someValue != nil {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssertInElseIf(t *testing.T) {
	for _, tc := range []struct {
		filenames    []string
		restartCount int
	}{
		{},
	} {
		t.Run("", func(t *testing.T) {
			count, err := calcRestartCountByLogDir(tc.filenames)
			if assert.NoError(t, err) {
				assert.Equal(t, count, tc.restartCount)
				assert.NoError(t, err)
			} else if true {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			} else {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			}

			if true {
				assert.Equal(t, count, tc.restartCount)
				assert.NoError(t, err)
			} else if assert.NoError(t, err) {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			} else {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			}

			assert.NoError(t, operation())
		})
	}
}
