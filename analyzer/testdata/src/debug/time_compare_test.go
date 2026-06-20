package debug

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeEquality(t *testing.T) {
	var tt time.Time
	assert.True(t, tt.IsZero())
	assert.Zero(t, tt)

	tn := time.Now()
	serialized, err := json.Marshal(tn)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(serialized, &tt))

	// Failed.
	/*
		-(time.Time) 2025-09-14 13:38:22.653679 +0300 MSK m=+0.002093459
		+(time.Time) 2025-09-14 13:38:22.653679 +0300 MSK
	*/
	assert.Equal(t, tn, tt)

	// Passed.
	assert.True(t, tn.Equal(tt))
	assert.Equal(t, tn.Unix(), tt.Unix())
	assert.Equal(t, tn.UnixMilli(), tt.UnixMilli())
	assert.Equal(t, tn.UnixMicro(), tt.UnixMicro())
	assert.Equal(t, tn.UnixNano(), tt.UnixNano())
	assert.WithinDuration(t, tn, tt, 0)
	assert.Equal(t, tn.Truncate(0), tt)
	assert.Equal(t, tn.Round(0), tt)
	assert.Equal(t, 0, tn.Compare(tt))

	s1 := metric{Value: 100, t: time.Now()}
	s2 := metric{Value: 100, t: time.Now()}
	assert.EqualExportedValues(t, s1, s2)

	tn = tn.In(time.FixedZone("CEST", 2*60*60))
	// Passed.
	assert.True(t, tn.Equal(tt))
	// Failed.
	assert.Equal(t, tn.Round(0), tt)
}

func TestTimeCompares(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Second)

	assert.Less(t, t1, t2)
	assert.True(t, t1.Compare(t2) < 0)
	assert.True(t, t1.Before(t2))
	assert.False(t, t1.Equal(t2))
	assert.Equal(t, -1, t1.Compare(t2))
	assert.Greater(t, t2, t1)
}

func TestZeroTime(t *testing.T) {
	var ts, zero time.Time
	now := time.Now()

	assert.Equal(t, time.Time{}, ts)
	assert.Equal(t, zero, ts)
	assert.True(t, ts.IsZero())
	assert.Zero(t, ts)

	assert.NotEqual(t, time.Time{}, now)
	assert.NotEqual(t, zero, now)
	assert.False(t, now.IsZero())
	assert.NotZero(t, now)

	assert.True(t, ts.Equal(zero))
	assert.True(t, ts.Equal(time.Time{}))
	assert.Equal(t, 0, ts.Compare(zero))
}

type metric struct {
	Value int
	t     time.Time
}
