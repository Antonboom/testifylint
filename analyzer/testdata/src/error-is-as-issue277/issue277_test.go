package errorisasissue277

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// a implements assert.TestingT (has Errorf method), so trimTArg will remove it
// from Args when used as first argument to IsType/IsNotType via assert.New(t).
type a struct{}

func (a) Errorf(format string, args ...any) {}

func TestIssue277(t *testing.T) {
	as := assert.New(t)

	// These should not panic (assert.New(t) calls don't pass t as first arg).
	// With assert.New(), the first arg to IsType is the "expected" value,
	// but since type a implements TestingT, trimTArg removes it, leaving only 1 arg.
	as.IsType(&a{}, &a{})
	as.IsNotType(&a{}, &a{})
}
