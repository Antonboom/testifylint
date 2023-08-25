package config_test

import (
	"flag"
	"slices"
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
	"github.com/Antonboom/testifylint/internal/config"
)

func TestNewDefault(t *testing.T) {
	cfg := config.NewDefault()
	if !slices.Equal(cfg.EnabledCheckers, checkers.EnabledByDefault()) {
		t.Fatal()
	}
	if cfg.ExpectedActual.ExpVarPattern.String() != checkers.DefaultExpectedVarPattern.String() {
		t.Fatal()
	}
}

func TestBindToFlags(t *testing.T) {
	cfg := config.NewDefault()
	fs := flag.NewFlagSet("bind-to-flags-test", flag.PanicOnError)

	config.BindToFlags(&cfg, fs)

	if fs.Lookup("enable-all").DefValue != "false" {
		t.Fatal()
	}
	if fs.Lookup("enable").DefValue != strings.Join(cfg.EnabledCheckers, ",") {
		t.Fatal()
	}
	if fs.Lookup("expected-actual.pattern").DefValue != cfg.ExpectedActual.ExpVarPattern.String() {
		t.Fatal()
	}
}
