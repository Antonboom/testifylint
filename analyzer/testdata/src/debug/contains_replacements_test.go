package debug

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsNotContainsReplacements_String(t *testing.T) {
	var tm tMock

	const substr = "abc123"

	for _, s := range []string{"", "random", substr + substr} {
		t.Run(s, func(t *testing.T) {
			assert.Equal(t,
				assert.True(tm, strings.Contains(s, substr)),
				assert.Contains(tm, s, substr),
			)

			assert.Equal(t,
				assert.False(tm, !strings.Contains(s, substr)),
				assert.Contains(tm, s, substr),
			)

			assert.Equal(t,
				assert.False(tm, strings.Contains(s, substr)),
				assert.NotContains(tm, s, substr),
			)

			assert.Equal(t,
				assert.True(tm, !strings.Contains(s, substr)),
				assert.NotContains(tm, s, substr),
			)
		})
	}
}

func TestContainsNotContainsReplacements_Bytes(t *testing.T) {
	var tm tMock

	subbytes := []byte("abc123")

	for _, s := range [][]byte{nil, []byte("random"), append(subbytes, subbytes...)} {
		t.Run(fmt.Sprintf("[]byte(%s)", s), func(t *testing.T) {
			assert.Equal(t,
				assert.True(tm, bytes.Contains(s, subbytes)),
				assert.Contains(tm, s, subbytes),
			)

			assert.Equal(t,
				assert.False(tm, !bytes.Contains(s, subbytes)),
				assert.Contains(tm, s, subbytes),
			)

			assert.Equal(t,
				assert.False(tm, bytes.Contains(s, subbytes)),
				assert.NotContains(tm, s, subbytes),
			)

			assert.Equal(t,
				assert.True(tm, !bytes.Contains(s, subbytes)),
				assert.NotContains(tm, s, subbytes),
			)
		})
	}
}
