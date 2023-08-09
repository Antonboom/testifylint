package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Config implements testifylint configuration.
type Config struct {
	EnabledCheckers []string             `yaml:"enabled-checkers"`
	ExpectedActual  ExpectedActualConfig `yaml:"expected-actual"`
}

type ExpectedActualConfig struct {
	Pattern Regexp `yaml:"pattern"`
}

// ParseFromFile parses Config from filepath. YAML format is expected.
func ParseFromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	return Parse(f)
}

// Parse parses Config from input. YAML format is expected.
func Parse(in io.Reader) (cfg Config, err error) {
	return cfg, yaml.NewDecoder(in).Decode(&cfg)
}
