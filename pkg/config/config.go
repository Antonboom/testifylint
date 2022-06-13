package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Checkers       CheckersConfig       `yaml:"checkers"`
	ExpectedActual ExpectedActualConfig `yaml:"expected-actual"`
}

type CheckersConfig struct {
	DisableAll bool     `yaml:"disable-all"`
	Enable     []string `yaml:"enable"`
	Disable    []string `yaml:"disable"`
}

type ExpectedActualConfig struct {
	Pattern string `yaml:"pattern"`
}

func ParseFromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	return Parse(f)
}

func Parse(in io.Reader) (cfg Config, err error) {
	return cfg, yaml.NewDecoder(in).Decode(&cfg)
}

func Dump(cfg Config, out io.Writer) error {
	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)
	return enc.Encode(cfg)
}
