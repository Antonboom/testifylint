package requireerrorskiplogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Run_Positive_DoNothing(t *testing.T) {
	var fakePlugin any

	assert.NoError(t, VerifyZeroAttachCalls(fakePlugin))
	assert.NoErrorf(t, VerifyZeroWaitForAttachCallCount(fakePlugin), "boom!")
	assert.NoError(t, VerifyZeroMountDeviceCallCount(fakePlugin))
	assert.NoError(t, VerifyZeroSetUpCallCount(fakePlugin))
	assert.NoErrorf(t, VerifyZeroTearDownCallCount(fakePlugin), "boom!")
	assert.NoErrorf(t, VerifyZeroDetachCallCount(fakePlugin), "boom!")
}

func PrepareDB(ctx context.Context, t *testing.T, dbName string) (client any, cleanUp func(ctx context.Context)) {
	t.Helper()
	require.NotEmpty(t, dbName)

	_, err := operationWithResult()
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	client, err = operationWithResult()
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	err = operation()
	assert.NoError(t, err)

	return client, func(ctx2 context.Context) {
		assert.NoError(t, operation())
		assert.NoError(t, operation())
		assert.NoError(t, operation())
	}
}

func VerifyZeroAttachCalls(any) error            { return nil } //
func VerifyZeroWaitForAttachCallCount(any) error { return nil } //
func VerifyZeroMountDeviceCallCount(any) error   { return nil } //
func VerifyZeroSetUpCallCount(any) error         { return nil } //
func VerifyZeroDetachCallCount(any) error        { return nil } //
func VerifyZeroTearDownCallCount(any) error      { return nil } //
