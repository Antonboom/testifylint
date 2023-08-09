package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/config"
)

func TestParse(t *testing.T) {
	const conf = `
enabled-checkers:
  - error-is
  - expected-actual
expected-actual:
  pattern: ^want$
`
	loaded, err := config.Parse(strings.NewReader(conf))
	if err != nil {
		t.Fatal(err)
	}

	expected := config.Config{
		EnabledCheckers: []string{"error-is", "expected-actual"},
		ExpectedActual: config.ExpectedActualConfig{
			Pattern: config.MustRegexp(`^want$`),
		},
	}
	if !reflect.DeepEqual(expected, loaded) {
		t.Fatal("config was parsed incorrectly")
	}
}

func TestParse_EmptyConfig(t *testing.T) {
	const conf = `
enabled-checkers:
expected-actual:
`
	loaded, err := config.Parse(strings.NewReader(conf))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config.Config{}, loaded) {
		t.Fatal("empty config was parsed incorrectly")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	const conf = `
enabled-checkers:
  - error-is
  -- expected-actual
`
	_, err := config.Parse(strings.NewReader(conf))
	if err == nil {
		t.Fatal("no error but expected")
	}
}

func TestParseFromFile_InvalidPath(t *testing.T) {
	_, err := config.ParseFromFile("unknown")
	if err == nil {
		t.Fatal("no error but expected")
	}
}
