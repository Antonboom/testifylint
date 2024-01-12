package debug

import (
	"testing"

	"github.com/ghetzel/testify/assert"
)

func TestNilReplacement(t *testing.T) {
	var tm tMock

	var (
		nilChan chan struct{}
		//nilFunc func()
		//nilInterface any
		//nilMap map[int]int
		//nilPointer *int
		//nilSlice []int
		//nilUnsafePointer unsafe.Pointer
	)

	var (
		notNilChan = make(chan struct{})
		//notNilFunc = TestNilReplacement
		//notNilInterface = any(notNilChan)
		//notNilMap = map[int]int{1: 1}
		//notNilPointer = new(int)
		//notNilSlice = []int{1}
		//notNilUnsafePointer = unsafe.Pointer(notNilPointer)
	)

	cases := []struct {
		name          string
		original      func() bool
		replaced      func() bool
		shouldBeEqual bool
	}{
		// Positive.
		{
			name:          "Equal nilChan 0",
			original:      func() bool { return assert.Equal(tm, nilChan, nil) },
			replaced:      func() bool { return assert.Nil(tm, nilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal 0 nilChan",
			original:      func() bool { return assert.Equal(tm, nil, nilChan) },
			replaced:      func() bool { return assert.Nil(tm, nilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal notNilChan 0",
			original:      func() bool { return assert.Equal(tm, notNilChan, nil) },
			replaced:      func() bool { return assert.Nil(tm, notNilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "Equal 0 notNilChan",
			original:      func() bool { return assert.Equal(tm, nil, notNilChan) },
			replaced:      func() bool { return assert.Nil(tm, notNilChan) },
			shouldBeEqual: true,
		},

		{
			name:          "EqualValues nilChan 0",
			original:      func() bool { return assert.EqualValues(tm, nilChan, nil) },
			replaced:      func() bool { return assert.Nil(tm, nilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues 0 nilChan",
			original:      func() bool { return assert.EqualValues(tm, nil, nilChan) },
			replaced:      func() bool { return assert.Nil(tm, nilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues notNilChan 0",
			original:      func() bool { return assert.EqualValues(tm, notNilChan, nil) },
			replaced:      func() bool { return assert.Nil(tm, notNilChan) },
			shouldBeEqual: true,
		},
		{
			name:          "EqualValues 0 notNilChan",
			original:      func() bool { return assert.EqualValues(tm, nil, notNilChan) },
			replaced:      func() bool { return assert.Nil(tm, notNilChan) },
			shouldBeEqual: true,
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
