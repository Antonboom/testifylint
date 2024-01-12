package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyReplacement(t *testing.T) {
	var tm tMock

	var nilSlice []string
	emptySlice := make([]string, 0)
	notEmptySlice := make([]string, 1)

	cases := []struct {
		name          string
		original      func() bool
		replaced      func() bool
		shouldBeEqual bool
	}{
		// Positive.
		{
			name:          "Equal nilSlice 0",
			original:      func() bool { return assert.Equal(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal 0 nilSlice",
			original:      func() bool { return assert.Equal(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues nilSlice 0",
			original:      func() bool { return assert.EqualValues(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues 0 nilSlice",
			original:      func() bool { return assert.EqualValues(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly nilSlice 0",
			original:      func() bool { return assert.Exactly(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly 0 nilSlice",
			original:      func() bool { return assert.Exactly(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual nilSlice 0",
			original:      func() bool { return assert.LessOrEqual(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual 0 nilSlice",
			original:      func() bool { return assert.LessOrEqual(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less nilSlice 1",
			original:      func() bool { return assert.Less(tm, len(nilSlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater 1 nilSlice",
			original:      func() bool { return assert.Greater(tm, 1, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},

		{
			name:          "Equal emptySlice 0",
			original:      func() bool { return assert.Equal(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal 0 emptySlice",
			original:      func() bool { return assert.Equal(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues emptySlice 0",
			original:      func() bool { return assert.EqualValues(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues 0 emptySlice",
			original:      func() bool { return assert.EqualValues(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly emptySlice 0",
			original:      func() bool { return assert.Exactly(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly 0 emptySlice",
			original:      func() bool { return assert.Exactly(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual emptySlice 0",
			original:      func() bool { return assert.LessOrEqual(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual 0 emptySlice",
			original:      func() bool { return assert.GreaterOrEqual(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less emptySlice 1",
			original:      func() bool { return assert.Less(tm, len(emptySlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater 1 emptySlice",
			original:      func() bool { return assert.Greater(tm, 1, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},

		{
			name:          "Equal notEmptySlice 0",
			original:      func() bool { return assert.Equal(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal 0 notEmptySlice",
			original:      func() bool { return assert.Equal(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues notEmptySlice 0",
			original:      func() bool { return assert.EqualValues(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues 0 notEmptySlice",
			original:      func() bool { return assert.EqualValues(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly notEmptySlice 0",
			original:      func() bool { return assert.Exactly(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Exactly 0 notEmptySlice",
			original:      func() bool { return assert.Exactly(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual notEmptySlice 0",
			original:      func() bool { return assert.LessOrEqual(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual 0 notEmptySlice",
			original:      func() bool { return assert.GreaterOrEqual(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less notEmptySlice 1",
			original:      func() bool { return assert.Less(tm, len(notEmptySlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater 1 notEmptySlice",
			original:      func() bool { return assert.Greater(tm, 1, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},

		// Negative.
		{
			name:          "LessOrEqual nilSlice 1",
			original:      func() bool { return assert.LessOrEqual(tm, len(nilSlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual emptySlice 1",
			original:      func() bool { return assert.LessOrEqual(tm, len(emptySlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual notEmptySlice 1",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(notEmptySlice), 1) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: false,
		},

		{
			name:          "GreaterOrEqual 1 nilSlice",
			original:      func() bool { return assert.GreaterOrEqual(tm, 1, len(nilSlice)) },
			replaced:      func() bool { return assert.Empty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual 1 emptySlice",
			original:      func() bool { return assert.GreaterOrEqual(tm, 1, len(emptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual 1 notEmptySlice",
			original:      func() bool { return assert.GreaterOrEqual(tm, 1, len(notEmptySlice)) },
			replaced:      func() bool { return assert.Empty(tm, notEmptySlice) },
			shouldBeEqual: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldBeEqual {
				assert.Equal(t, tt.original(), tt.replaced())
			} else {
				assert.NotEqual(t, tt.original(), tt.replaced())
			}
		})
	}
}
