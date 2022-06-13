package config_test

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/Antonboom/testifylint/pkg/config"
)

var configExamplePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	configExamplePath = filepath.Join(filepath.Dir(currentFile), "..", "..", "config.example.yml")
}

func TestDefaultConfigFilled(t *testing.T) {
	cfg := config.Default

	if len(cfg.Checkers.Enable) == 0 {
		t.FailNow()
	}
	if len(cfg.Checkers.Disable) == 0 {
		t.FailNow()
	}
	if cfg.ExpectedActual.Pattern == "" {
		t.FailNow()
	}
}

func TestDefaultConfigIsValid(t *testing.T) {
	err := config.Validate(config.Default)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestConfigExampleIsActual(t *testing.T) {
	cfg, err := config.ParseFromFile(configExamplePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(cfg, config.Default) {
		t.Fatal("run `task config:dump`")
	}
}
