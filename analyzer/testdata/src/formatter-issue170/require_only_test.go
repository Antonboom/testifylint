package formatterissue170

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	require.True(t, false, fmt.Sprintf("expected %v, got %v", true, false)) // want "formatter: remove unnecessary fmt\\.Sprintf"
}
