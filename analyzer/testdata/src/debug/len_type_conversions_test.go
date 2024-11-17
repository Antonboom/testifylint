package debug

import (
	"encoding/json"
	"testing"

	"github.com/ghetzel/testify/assert"
)

func TestLenTypeConversions(t *testing.T) {
	const resp = "Mutlti-языковая string 你好世界 🙂"
	const respLen = len(resp) // 48

	assert.Equal(t, respLen, len(resp))
	assert.Equal(t, respLen, len([]byte(resp)))
	assert.Equal(t, respLen, len(json.RawMessage(resp)))
	assert.Len(t, resp, respLen)
	assert.Len(t, []byte(resp), respLen)
	assert.Len(t, json.RawMessage(resp), respLen)
}
