package debug

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorAsGenericTarget(t *testing.T) {
	var err error = new(os.PathError)

	errorAs[*os.PathError]()(t, err)
	errorAs[*os.LinkError]()(t, err)
}

func errorAs[T error]() require.ErrorAssertionFunc {
	return func(t require.TestingT, err error, msgAndArgs ...interface{}) {
		var target T
		require.ErrorAs(t, err, &target, msgAndArgs...)
	}
}
