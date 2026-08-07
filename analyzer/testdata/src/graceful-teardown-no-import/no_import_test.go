package gracefulteardownnoimport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGracefulTeardownNoImport verifies that the autofix adds the "assert" import
// when it is absent from the file.
func TestGracefulTeardownNoImport(t *testing.T) {
	var err error
	t.Cleanup(func() {
		require.NoError(t, err) // want "graceful-teardown: do not use require in cleanup code"
	})
}
