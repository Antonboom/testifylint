package config

import (
	"regexp"

	"gopkg.in/yaml.v3"
)

// Regexp is a special wrapper for regexp.Regexp YAML unmarshalling support.
// Original regexp is available through Regexp.Regexp.
type Regexp struct {
	*regexp.Regexp
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (re *Regexp) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	compiled, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	re.Regexp = compiled
	return nil
}
