package config

import (
	"regexp"

	"gopkg.in/yaml.v3"
)

// Regexp is a special wrapper for YAML marshalling and unmarshalling support.
// Original regexp is available through Regexp.Regexp.
type Regexp struct {
	*regexp.Regexp
}

// NewRegexp compiles Regexp from input string.
// Returns error if string is invalid regular expression.
func NewRegexp(s string) (Regexp, error) {
	r, err := regexp.Compile(s)
	if err != nil {
		return Regexp{}, err
	}
	return Regexp{r}, nil
}

// MustRegexp compiles Regexp from input string.
// Panics if string is invalid regular expression.
func MustRegexp(s string) Regexp {
	r, err := NewRegexp(s)
	if err != nil {
		panic(err)
	}
	return r
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (re *Regexp) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	compiled, err := NewRegexp(s)
	if err != nil {
		return err
	}

	*re = compiled
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (re Regexp) MarshalYAML() (any, error) {
	if re.Regexp == nil {
		return nil, nil
	}
	return re.Regexp.String(), nil
}
