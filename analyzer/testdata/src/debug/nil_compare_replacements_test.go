package debug

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

var (
	nilChan          chan struct{}
	nilFunc          func()
	nilInterface     any
	nilMap           map[int]int
	nilPointer       *int
	nilSlice         []int
	nilUnsafePointer unsafe.Pointer

	notNilChan          = make(chan struct{})
	notNilFunc          = tMock{}.Errorf
	notNilInterface     = any(notNilChan)
	notNilMap           = map[int]int{1: 1}
	notNilPointer       = new(int)
	notNilSlice         = []int{1}
	notNilUnsafePointer = unsafe.Pointer(notNilPointer)
)

var nillables = map[string]any{
	"nilChan":          nilChan,
	"nilFunc":          nilFunc,
	"nilInterface":     nilInterface,
	"nilMap":           nilMap,
	"nilPointer":       nilPointer,
	"nilSlice":         nilSlice,
	"nilUnsafePointer": nilUnsafePointer,

	"notNilChan":          notNilChan,
	"notNilFunc":          notNilFunc,
	"notNilInterface":     notNilInterface,
	"notNilMap":           notNilMap,
	"notNilPointer":       notNilPointer,
	"notNilSlice":         notNilSlice,
	"notNilUnsafePointer": notNilUnsafePointer,
}

func TestNilReplacements(t *testing.T) {
	var tm tMock

	for varName, v := range nillables {
		t.Run(varName, func(t *testing.T) {
			for assrnName, assrn := range map[string]assertion{
				"Equal":       assert.Equal,
				"EqualValues": assert.EqualValues,
				"Exactly":     assert.Exactly,
			} {
				t.Run(assrnName+"(v,nil)", func(t *testing.T) {
					original := assrn(tm, v, nil)
					replacement := assert.Nil(tm, v)
					assert.Equal(t, original, replacement, "not an equivalent replacement")
				})

				t.Run(assrnName+"(nil,v)", func(t *testing.T) {
					original := assrn(tm, nil, v)
					replacement := assert.Nil(tm, v)
					assert.Equal(t, original, replacement, "not an equivalent replacement")
				})
			}
		})
	}
}

func TestNotNilReplacements(t *testing.T) {
	var tm tMock

	for varName, v := range nillables {
		t.Run(varName, func(t *testing.T) {
			for assrnName, assrn := range map[string]assertion{
				"NotEqual":       assert.NotEqual,
				"NotEqualValues": assert.NotEqualValues,
			} {
				t.Run(assrnName+"(v,nil)", func(t *testing.T) {
					original := assrn(tm, v, nil)
					replacement := assert.NotNil(tm, v)
					assert.Equal(t, original, replacement, "not an equivalent replacement")
				})

				t.Run(assrnName+"(nil,v)", func(t *testing.T) {
					original := assrn(tm, nil, v)
					replacement := assert.NotNil(tm, v)
					assert.Equal(t, original, replacement, "not an equivalent replacement")
				})
			}
		})
	}
}
