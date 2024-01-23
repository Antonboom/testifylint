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

	// n := len(elems)
	// n == 0, n <= 0, n < 1
	// 0 == n, 0 >= n, 1 > n
	cases := []struct {
		name        string
		original    func() bool
		replacement func() bool
	}{
		{
			name:        "Equal len(nilSlice) 0",
			original:    func() bool { return assert.Equal(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Equal 0 len(nilSlice)",
			original:    func() bool { return assert.Equal(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Equal len(emptySlice) 0",
			original:    func() bool { return assert.Equal(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Equal 0 len(emptySlice)",
			original:    func() bool { return assert.Equal(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Equal len(notEmptySlice) 0",
			original:    func() bool { return assert.Equal(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
		{
			name:        "Equal 0 len(notEmptySlice)",
			original:    func() bool { return assert.Equal(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},

		{
			name:        "EqualValues len(nilSlice) 0",
			original:    func() bool { return assert.EqualValues(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "EqualValues 0 len(nilSlice)",
			original:    func() bool { return assert.EqualValues(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "EqualValues len(emptySlice) 0",
			original:    func() bool { return assert.EqualValues(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "EqualValues 0 len(emptySlice)",
			original:    func() bool { return assert.EqualValues(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "EqualValues len(notEmptySlice) 0",
			original:    func() bool { return assert.EqualValues(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
		{
			name:        "EqualValues 0 len(notEmptySlice)",
			original:    func() bool { return assert.EqualValues(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},

		{
			name:        "Exactly len(nilSlice) 0",
			original:    func() bool { return assert.Exactly(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Exactly 0 len(nilSlice)",
			original:    func() bool { return assert.Exactly(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Exactly len(emptySlice) 0",
			original:    func() bool { return assert.Exactly(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Exactly 0 len(emptySlice)",
			original:    func() bool { return assert.Exactly(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Exactly len(notEmptySlice) 0",
			original:    func() bool { return assert.Exactly(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
		{
			name:        "Exactly 0 len(notEmptySlice)",
			original:    func() bool { return assert.Exactly(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},

		{
			name:        "LessOrEqual len(nilSlice) 0",
			original:    func() bool { return assert.LessOrEqual(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "GreaterOrEqual 0 len(nilSlice)",
			original:    func() bool { return assert.GreaterOrEqual(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "LessOrEqual len(emptySlice) 0",
			original:    func() bool { return assert.LessOrEqual(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "GreaterOrEqual 0 len(emptySlice)",
			original:    func() bool { return assert.GreaterOrEqual(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "LessOrEqual len(notEmptySlice) 0",
			original:    func() bool { return assert.LessOrEqual(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
		{
			name:        "GreaterOrEqual 0 len(notEmptySlice)",
			original:    func() bool { return assert.GreaterOrEqual(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},

		{
			name:        "Less len(nilSlice) 1",
			original:    func() bool { return assert.Less(tm, len(nilSlice), 1) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Greater 1 len(nilSlice)",
			original:    func() bool { return assert.Greater(tm, 1, len(nilSlice)) },
			replacement: func() bool { return assert.Empty(tm, nilSlice) },
		},
		{
			name:        "Less len(emptySlice) 1",
			original:    func() bool { return assert.Less(tm, len(emptySlice), 1) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Greater 1 len(emptySlice)",
			original:    func() bool { return assert.Greater(tm, 1, len(emptySlice)) },
			replacement: func() bool { return assert.Empty(tm, emptySlice) },
		},
		{
			name:        "Less len(notEmptySlice) 1",
			original:    func() bool { return assert.Less(tm, len(notEmptySlice), 1) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
		{
			name:        "Greater 1 len(notEmptySlice)",
			original:    func() bool { return assert.Greater(tm, 1, len(notEmptySlice)) },
			replacement: func() bool { return assert.Empty(tm, notEmptySlice) },
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.original(), tt.replacement(), "not an equivalent replacement")
		})
	}
}

func TestNotEmptyReplacement(t *testing.T) {
	var tm tMock

	var nilSlice []string
	emptySlice := make([]string, 0)
	notEmptySlice := make([]string, 1)
	threeElemSlice := make([]string, 3)

	// n := len(elems)
	// n != 0, n > 0, n >= 1, n > 1
	// 0 != n, 0 < n, n <= 1, 1 < n
	cases := []struct {
		name        string
		original    func() bool
		replacement func() bool
	}{
		{
			name:        "NotEqual len(nilSlice) 0",
			original:    func() bool { return assert.NotEqual(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "NotEqual 0 len(nilSlice)",
			original:    func() bool { return assert.NotEqual(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "NotEqual len(emptySlice) 0",
			original:    func() bool { return assert.NotEqual(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "NotEqual 0 len(emptySlice)",
			original:    func() bool { return assert.NotEqual(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "NotEqual len(notEmptySlice) 0",
			original:    func() bool { return assert.NotEqual(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "NotEqual 0 len(notEmptySlice)",
			original:    func() bool { return assert.NotEqual(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "NotEqual len(notEmptySlice) 0",
			original:    func() bool { return assert.NotEqual(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "NotEqual 0 len(notEmptySlice)",
			original:    func() bool { return assert.NotEqual(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "NotEqual len(threeElemSlice) 0",
			original:    func() bool { return assert.NotEqual(tm, len(threeElemSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
		{
			name:        "NotEqual 0 len(threeElemSlice)",
			original:    func() bool { return assert.NotEqual(tm, 0, len(threeElemSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},

		{
			name:        "NotEqualValues len(nilSlice) 0",
			original:    func() bool { return assert.NotEqualValues(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "NotEqualValues 0 len(nilSlice)",
			original:    func() bool { return assert.NotEqualValues(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "NotEqualValues len(emptySlice) 0",
			original:    func() bool { return assert.NotEqualValues(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "NotEqualValues 0 len(emptySlice)",
			original:    func() bool { return assert.NotEqualValues(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "NotEqualValues len(threeElemSlice) 0",
			original:    func() bool { return assert.NotEqualValues(tm, len(threeElemSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
		{
			name:        "NotEqualValues 0 len(threeElemSlice)",
			original:    func() bool { return assert.NotEqualValues(tm, 0, len(threeElemSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},

		{
			name:        "Greater len(nilSlice) 0",
			original:    func() bool { return assert.Greater(tm, len(nilSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "Less 0 len(nilSlice)",
			original:    func() bool { return assert.Less(tm, 0, len(nilSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "Greater len(emptySlice) 0",
			original:    func() bool { return assert.Greater(tm, len(emptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "Less 0 len(emptySlice)",
			original:    func() bool { return assert.Less(tm, 0, len(emptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "Greater len(notEmptySlice) 0",
			original:    func() bool { return assert.Greater(tm, len(notEmptySlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "Less 0 len(notEmptySlice)",
			original:    func() bool { return assert.Less(tm, 0, len(notEmptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "Greater len(threeElemSlice) 0",
			original:    func() bool { return assert.Greater(tm, len(threeElemSlice), 0) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
		{
			name:        "Less 0 len(threeElemSlice)",
			original:    func() bool { return assert.Less(tm, 0, len(threeElemSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},

		{
			name:        "GreaterOrEqual len(nilSlice) 1",
			original:    func() bool { return assert.GreaterOrEqual(tm, len(nilSlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "LessOrEqual 1, len(nilSlice)",
			original:    func() bool { return assert.LessOrEqual(tm, 1, len(nilSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "GreaterOrEqual len(emptySlice) 1",
			original:    func() bool { return assert.GreaterOrEqual(tm, len(emptySlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "LessOrEqual 1, len(emptySlice)",
			original:    func() bool { return assert.LessOrEqual(tm, 1, len(emptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "GreaterOrEqual len(notEmptySlice) 1",
			original:    func() bool { return assert.GreaterOrEqual(tm, len(notEmptySlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "LessOrEqual 1, len(notEmptySlice)",
			original:    func() bool { return assert.LessOrEqual(tm, 1, len(notEmptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "GreaterOrEqual len(threeElemSlice) 1",
			original:    func() bool { return assert.GreaterOrEqual(tm, len(threeElemSlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
		{
			name:        "LessOrEqual 1, len(threeElemSlice)",
			original:    func() bool { return assert.LessOrEqual(tm, 1, len(threeElemSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},

		{
			name:        "Greater len(nilSlice) 1",
			original:    func() bool { return assert.Greater(tm, len(nilSlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "Less 1, len(nilSlice)",
			original:    func() bool { return assert.Less(tm, 1, len(nilSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, nilSlice) },
		},
		{
			name:        "Greater len(emptySlice) 1",
			original:    func() bool { return assert.Greater(tm, len(emptySlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "Less 1, len(emptySlice)",
			original:    func() bool { return assert.Less(tm, 1, len(emptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, emptySlice) },
		},
		{
			name:        "Greater len(notEmptySlice) 1",
			original:    func() bool { return assert.Greater(tm, len(notEmptySlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "Less 1, len(notEmptySlice)",
			original:    func() bool { return assert.Less(tm, 1, len(notEmptySlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, notEmptySlice) },
		},
		{
			name:        "Greater len(threeElemSlice) 1",
			original:    func() bool { return assert.Greater(tm, len(threeElemSlice), 1) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
		{
			name:        "Less 1, len(threeElemSlice)",
			original:    func() bool { return assert.Less(tm, 1, len(threeElemSlice)) },
			replacement: func() bool { return assert.NotEmpty(tm, threeElemSlice) },
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.original(), tt.replacement(), "not an equivalent replacement")
		})
	}
}
