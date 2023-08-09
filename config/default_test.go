package config_test

import (
	"bytes"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/Antonboom/testifylint/config"
)

var defaultConfigPath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	defaultConfigPath = filepath.Join(filepath.Dir(currentFile), "..", ".testifylint.yml")
}

func TestDumpDefault(t *testing.T) {
	b := bytes.NewBuffer(nil)
	if err := config.DumpDefault(b); err != nil {
		t.Fatal(err)
	}
	rawConfig := b.String()

	loaded, err := config.Parse(b)
	if err != nil {
		t.Fatalf("%v:\n%s", err, rawConfig)
	}

	if !reflect.DeepEqual(loaded, config.Default) {
		t.Fatalf("dumped config is not equal with default config:\n%s", rawConfig)
	}
}

func TestDumpedDefaultConfigIsActual(t *testing.T) {
	cfg, err := config.ParseFromFile(defaultConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	err = config.Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cfg, config.Default) {
		t.Fatal("dumped default config is outdated")
	}
}
