package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmptyReplacement(t *testing.T) {
	var tm tMock

	var nilSlice []string
	emptySlice := make([]string, 0)
	notEmptySlice := make([]string, 1)
	threeElemSlice := make([]string, 3)

	cases := []struct {
		name          string
		original      func() bool
		replaced      func() bool
		shouldBeEqual bool
	}{
		// Positive.
		{
			name:          "NotEqual nilSlice 0",
			original:      func() bool { return assert.NotEqual(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqual 0 nilSlice",
			original:      func() bool { return assert.NotEqual(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues nilSlice 0",
			original:      func() bool { return assert.NotEqualValues(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues 0 nilSlice",
			original:      func() bool { return assert.NotEqualValues(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater nilSlice 0",
			original:      func() bool { return assert.Greater(tm, len(nilSlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less 0 nilSlice",
			original:      func() bool { return assert.Less(tm, 0, len(nilSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual nilSlice 1",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(nilSlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual 1 nilSlice",
			original:      func() bool { return assert.LessOrEqual(tm, 1, len(nilSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},

		{
			name:          "NotEqual emptySlice 0",
			original:      func() bool { return assert.NotEqual(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqual 0 emptySlice",
			original:      func() bool { return assert.NotEqual(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues emptySlice 0",
			original:      func() bool { return assert.NotEqualValues(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues 0 emptySlice",
			original:      func() bool { return assert.NotEqualValues(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater emptySlice 0",
			original:      func() bool { return assert.Greater(tm, len(emptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less 0 emptySlice",
			original:      func() bool { return assert.Less(tm, 0, len(emptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual emptySlice 1",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(emptySlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual 1 emptySlice",
			original:      func() bool { return assert.LessOrEqual(tm, 1, len(emptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},

		{
			name:          "NotEqual notEmptySlice 0",
			original:      func() bool { return assert.NotEqual(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqual 0 notEmptySlice",
			original:      func() bool { return assert.NotEqual(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues notEmptySlice 0",
			original:      func() bool { return assert.NotEqualValues(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "NotEqualValues 0 notEmptySlice",
			original:      func() bool { return assert.NotEqualValues(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater notEmptySlice 0",
			original:      func() bool { return assert.Greater(tm, len(notEmptySlice), 0) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less 0 notEmptySlice",
			original:      func() bool { return assert.Less(tm, 0, len(notEmptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual notEmptySlice 1",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(notEmptySlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual 1 notEmptySlice",
			original:      func() bool { return assert.LessOrEqual(tm, 1, len(notEmptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "GreaterOrEqual threeElemSlice 1",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(threeElemSlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, threeElemSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "LessOrEqual 1 threeElemSlice",
			original:      func() bool { return assert.LessOrEqual(tm, 1, len(threeElemSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, threeElemSlice) },
			shouldBeEqual: true,
		},

		// Negative.
		{
			name:          "Greater nilSlice 1",
			original:      func() bool { return assert.Greater(tm, len(nilSlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater emptySlice 1",
			original:      func() bool { return assert.Greater(tm, len(emptySlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Greater notEmptySlice 1",
			original:      func() bool { return assert.Greater(tm, len(notEmptySlice), 1) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: false,
		},

		{
			name:          "Less 1 nilSlice",
			original:      func() bool { return assert.Less(tm, 1, len(nilSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, nilSlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less 1 emptySlice",
			original:      func() bool { return assert.Less(tm, 1, len(emptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, emptySlice) },
			shouldBeEqual: true,
		},
		{
			name:          "Less 1 notEmptySlice",
			original:      func() bool { return assert.Less(tm, 1, len(notEmptySlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, notEmptySlice) },
			shouldBeEqual: false,
		},

		{
			name:          "GreaterOrEqual threeElemSlice 4",
			original:      func() bool { return assert.GreaterOrEqual(tm, len(threeElemSlice), 4) },
			replaced:      func() bool { return assert.NotEmpty(tm, threeElemSlice) },
			shouldBeEqual: false,
		},
		{
			name:          "LessOrEqual threeElemSlice 2",
			original:      func() bool { return assert.LessOrEqual(tm, 4, len(threeElemSlice)) },
			replaced:      func() bool { return assert.NotEmpty(tm, threeElemSlice) },
			shouldBeEqual: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldBeEqual, tt.original() == tt.replaced())
		})
	}
}
