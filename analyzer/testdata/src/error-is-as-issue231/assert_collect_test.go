package errorisasissue231

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventuallyAsserts(t *testing.T) {
	require.EventuallyWithT(t,
		func(c *assert.CollectT) {
			_, err := strconv.Atoi("a")
			if err != nil {
				c.Errorf("failed: %v", err)

				c.Errorf(fmt.Sprintf("failed: %v", err)) // want "formatter: remove unnecessary fmt\\.Sprintf"
				c.Errorf("failed: %v")                   // want "formatter: c\\.Errorf format %v reads arg #1, but call has 0 args"
				return
			}
		},
		time.Second,
		time.Millisecond,
	)
}
