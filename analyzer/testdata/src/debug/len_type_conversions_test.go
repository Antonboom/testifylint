package debug

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLenTypeConversions(t *testing.T) {
	const resp = "Multi-языковая string 你好世界 🙂"
	const respLen = len(resp) // 47

	assert.Equal(t, respLen, len(resp))
	assert.Equal(t, respLen, len([]byte(resp)))
	assert.Equal(t, respLen, len(json.RawMessage(resp)))
	assert.Len(t, resp, respLen)
	assert.Len(t, []byte(resp), respLen)
	assert.Len(t, json.RawMessage(resp), respLen)
}

func TestBytesReadability(t *testing.T) {
	resp := []byte(`{"status": 200}`)
	assert.Len(t, resp, 3)         // "[123 34 115 116 97 116 117 115 34 58 32 50 48 48 125]" should have 3 item(s), but has 15
	assert.Len(t, string(resp), 3) // "{"status": 200}" should have 3 item(s), but has 15

	var noResp []byte
	assert.Empty(t, resp)              // Should be empty, but was [123 34 115 116 97 116 117 115 34 58 32 50 48 48 125]
	assert.Empty(t, string(resp))      // Should be empty, but was {"status": 200}
	assert.NotEmpty(t, noResp)         // Should NOT be empty, but was []
	assert.NotEmpty(t, string(noResp)) // Should NOT be empty, but was
}
