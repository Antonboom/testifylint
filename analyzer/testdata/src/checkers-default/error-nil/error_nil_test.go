// Code generated by testifylint/internal/testgen. DO NOT EDIT.

package errornil

import (
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestErrorNilChecker(t *testing.T) {
	var err error

	// Invalid.
	{
		assert.Nil(t, err)                                                   // want "error-nil: use assert\\.NoError"
		assert.Nilf(t, err, "msg with args %d %s", 42, "42")                 // want "error-nil: use assert\\.NoErrorf"
		assert.NotNil(t, err)                                                // want "error-nil: use assert\\.Error"
		assert.NotNilf(t, err, "msg with args %d %s", 42, "42")              // want "error-nil: use assert\\.Errorf"
		assert.Empty(t, err)                                                 // want "error-nil: use assert\\.NoError"
		assert.Emptyf(t, err, "msg with args %d %s", 42, "42")               // want "error-nil: use assert\\.NoErrorf"
		assert.NotEmpty(t, err)                                              // want "error-nil: use assert\\.Error"
		assert.NotEmptyf(t, err, "msg with args %d %s", 42, "42")            // want "error-nil: use assert\\.Errorf"
		assert.Zero(t, err)                                                  // want "error-nil: use assert\\.NoError"
		assert.Zerof(t, err, "msg with args %d %s", 42, "42")                // want "error-nil: use assert\\.NoErrorf"
		assert.NotZero(t, err)                                               // want "error-nil: use assert\\.Error"
		assert.NotZerof(t, err, "msg with args %d %s", 42, "42")             // want "error-nil: use assert\\.Errorf"
		assert.Equal(t, err, nil)                                            // want "error-nil: use assert\\.NoError"
		assert.Equalf(t, err, nil, "msg with args %d %s", 42, "42")          // want "error-nil: use assert\\.NoErrorf"
		assert.Equal(t, nil, err)                                            // want "error-nil: use assert\\.NoError"
		assert.Equalf(t, nil, err, "msg with args %d %s", 42, "42")          // want "error-nil: use assert\\.NoErrorf"
		assert.EqualValues(t, err, nil)                                      // want "error-nil: use assert\\.NoError"
		assert.EqualValuesf(t, err, nil, "msg with args %d %s", 42, "42")    // want "error-nil: use assert\\.NoErrorf"
		assert.EqualValues(t, nil, err)                                      // want "error-nil: use assert\\.NoError"
		assert.EqualValuesf(t, nil, err, "msg with args %d %s", 42, "42")    // want "error-nil: use assert\\.NoErrorf"
		assert.Exactly(t, err, nil)                                          // want "error-nil: use assert\\.NoError"
		assert.Exactlyf(t, err, nil, "msg with args %d %s", 42, "42")        // want "error-nil: use assert\\.NoErrorf"
		assert.Exactly(t, nil, err)                                          // want "error-nil: use assert\\.NoError"
		assert.Exactlyf(t, nil, err, "msg with args %d %s", 42, "42")        // want "error-nil: use assert\\.NoErrorf"
		assert.NotEqual(t, err, nil)                                         // want "error-nil: use assert\\.Error"
		assert.NotEqualf(t, err, nil, "msg with args %d %s", 42, "42")       // want "error-nil: use assert\\.Errorf"
		assert.NotEqual(t, nil, err)                                         // want "error-nil: use assert\\.Error"
		assert.NotEqualf(t, nil, err, "msg with args %d %s", 42, "42")       // want "error-nil: use assert\\.Errorf"
		assert.NotEqualValues(t, err, nil)                                   // want "error-nil: use assert\\.Error"
		assert.NotEqualValuesf(t, err, nil, "msg with args %d %s", 42, "42") // want "error-nil: use assert\\.Errorf"
		assert.NotEqualValues(t, nil, err)                                   // want "error-nil: use assert\\.Error"
		assert.NotEqualValuesf(t, nil, err, "msg with args %d %s", 42, "42") // want "error-nil: use assert\\.Errorf"
		assert.ErrorIs(t, err, nil)                                          // want "error-nil: use assert\\.NoError"
		assert.ErrorIsf(t, err, nil, "msg with args %d %s", 42, "42")        // want "error-nil: use assert\\.NoErrorf"
		assert.NotErrorIs(t, err, nil)                                       // want "error-nil: use assert\\.Error"
		assert.NotErrorIsf(t, err, nil, "msg with args %d %s", 42, "42")     // want "error-nil: use assert\\.Errorf"
	}

	// Valid.
	{
		assert.NoError(t, err)
		assert.NoErrorf(t, err, "msg with args %d %s", 42, "42")
		assert.Error(t, err)
		assert.Errorf(t, err, "msg with args %d %s", 42, "42")
	}

	// Ignored.
	{
		assert.Nil(t, nil)
		assert.Nilf(t, nil, "msg with args %d %s", 42, "42")
		assert.NotNil(t, nil)
		assert.NotNilf(t, nil, "msg with args %d %s", 42, "42")
		assert.Equal(t, err, err)
		assert.Equalf(t, err, err, "msg with args %d %s", 42, "42")
		assert.Equal(t, nil, nil)
		assert.Equalf(t, nil, nil, "msg with args %d %s", 42, "42")
		assert.NotEqual(t, err, err)
		assert.NotEqualf(t, err, err, "msg with args %d %s", 42, "42")
		assert.NotEqual(t, nil, nil)
		assert.NotEqualf(t, nil, nil, "msg with args %d %s", 42, "42")
		assert.Empty(t, err.Error())
		assert.Emptyf(t, err.Error(), "msg with args %d %s", 42, "42")
		assert.NotEmpty(t, err.Error())
		assert.NotEmptyf(t, err.Error(), "msg with args %d %s", 42, "42")
		assert.Zero(t, err.Error())
		assert.Zerof(t, err.Error(), "msg with args %d %s", 42, "42")
		assert.NotZero(t, err.Error())
		assert.NotZerof(t, err.Error(), "msg with args %d %s", 42, "42")
	}
}

func TestErrorNilChecker_ErrorDetection(t *testing.T) {
	errOp := func() error { return io.EOF }
	var a error
	var b withErroredMethod
	_, c := b.Get2()

	assert.Nil(t, a)        // want "error-nil: use assert\\.NoError"
	assert.Nil(t, b.Get1()) // want "error-nil: use assert\\.NoError"
	assert.Nil(t, c)        // want "error-nil: use assert\\.NoError"
	assert.Nil(t, errOp())  // want "error-nil: use assert\\.NoError"
}

func TestErrorNilChecker_ValidNils(t *testing.T) {
	var (
		ptr   *int
		iface any
		ch    chan error
		sl    []error
		fn    func()
		m     map[int]int
		uPtr  unsafe.Pointer
	)

	assert.Nil(t, ptr)
	assert.NotNil(t, ptr)
	assert.Nil(t, iface)
	assert.NotNil(t, iface)
	assert.Nil(t, ch)
	assert.NotNil(t, ch)
	assert.Nil(t, sl)
	assert.NotNil(t, sl)
	assert.Nil(t, fn)
	assert.NotNil(t, fn)
	assert.Nil(t, m)
	assert.NotNil(t, m)
	assert.Nil(t, uPtr)
	assert.NotNil(t, uPtr)
}

type withErroredMethod struct{}

func (withErroredMethod) Get1() error        { return nil }
func (withErroredMethod) Get2() (int, error) { return 0, nil }
