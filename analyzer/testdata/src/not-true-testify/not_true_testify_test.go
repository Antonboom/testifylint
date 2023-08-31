// Test inspired by
// https://github.com/ghetzel/hydra/blob/ea80690dfd4e28fa87f2b310bb524eee880824d5/gen_test.go#L1
package nottruetestify

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestGenerateBasic(t *testing.T) {
	assert := require.New(t)

	assert.Equal(NewComponent(`NonexistingThing`).String(), `NonexistingThing{}`)

	win := NewComponent(`ApplicationWindow`)
	win.Set(`visible`, true)
	win.Set(`color`, `#FF00CC`)

	assert.Equal(win.String(), "ApplicationWindow {\n  visible: true;\n  color: \"#FF00CC\";\n}")
}

type Component struct {
	Type       string                 `yaml:"type,omitempty"       json:"type,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty" json:"properties,omitempty"`
}

func NewComponent(ctype string) *Component {
	return &Component{
		Type: ctype,
	}
}

func (self *Component) Set(key string, value interface{}) {
	if self.Properties == nil {
		self.Properties = make(map[string]interface{})
	}

	self.Properties[key] = value
}

func (self *Component) String() string {
	return self.Type
}
