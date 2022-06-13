package config_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/Antonboom/testifylint/pkg/config"
)

func TestParseFromFile_InvalidPath(t *testing.T) {
	_, err := config.ParseFromFile("unknown")
	if err == nil {
		t.Fatal("no error but expected")
	}
}

func TestDumpParse(t *testing.T) {
	cfg := config.Config{
		Checkers: config.CheckersConfig{
			DisableAll: true,
			Enable:     []string{"len"},
			Disable:    []string{"empty"},
		},
		ExpectedActual: config.ExpectedActualConfig{
			Pattern: `^want$`,
		},
	}

	b := bytes.NewBuffer(nil)
	if err := config.Dump(cfg, b); err != nil {
		t.Fatal(err.Error())
	}

	exp := `checkers:
  disable-all: true
  enable:
    - len
  disable:
    - empty
expected-actual:
  pattern: ^want$
`
	if rawConf := b.String(); rawConf != exp {
		t.Fatalf("unexpected result: %q", rawConf)
	}

	loaded, err := config.Parse(b)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(loaded, cfg) {
		t.Fatal("parsed back config is different from original")
	}
}
