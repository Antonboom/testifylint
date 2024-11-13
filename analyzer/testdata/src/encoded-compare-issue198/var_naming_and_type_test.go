package encodedcompareissue198

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonResultFromFullResult(t *testing.T) {
	jsonData := new(bytes.Buffer)
	var cr *commonResult

	crFromJSON := new(commonResult)
	err := json.Unmarshal(jsonData.Bytes(), crFromJSON)
	require.NoError(t, err)

	assert.Equal(t, cr, crFromJSON)
}

func TestCommonResultFromFullAndCompactJSON(t *testing.T) {
	compactJSONData := new(bytes.Buffer)
	fullJSONData := new(bytes.Buffer)

	crFromCompactJSON := new(commonResult)
	crFromFullJSON := new(commonResult)

	err := json.NewDecoder(compactJSONData).Decode(crFromCompactJSON)
	require.NoError(t, err)

	err = json.NewDecoder(fullJSONData).Decode(crFromCompactJSON)
	require.NoError(t, err)

	assert.Equal(t, crFromCompactJSON, crFromFullJSON)
}

type commonResult struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}
