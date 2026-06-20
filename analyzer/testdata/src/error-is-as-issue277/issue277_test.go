package errorisasissue277

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// a structurally implements `assert.TestingT` and error.
// The analyzer must not mistake it for the testing argument of a package-level assertion.
type a struct{}

func (a) Errorf(format string, args ...any) {}
func (a) Error() string                     { return "" }

func TestIssue277(t *testing.T) {
	as := assert.New(t)
	as.IsType(&a{}, &a{})
	as.IsNotType(&a{}, &a{})
}
