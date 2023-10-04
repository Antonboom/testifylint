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

	if cfg.EnableAll {
		t.Fatal()
	}
	if !slices.Equal(cfg.EnabledCheckers, checkers.EnabledByDefault()) {
		t.Fatal()
	}
	if cfg.ExpectedActual.ExpVarPattern.String() != checkers.DefaultExpectedVarPattern.String() {
		t.Fatal()
	}
	if cfg.SuiteExtraAssertCall.Mode != checkers.SuiteExtraAssertCallModeRemove {
		t.Fatal()
	}
}

func TestBindToFlags(t *testing.T) {
	cfg := config.NewDefault()
	fs := flag.NewFlagSet("TestBindToFlags", flag.PanicOnError)

	config.BindToFlags(&cfg, fs)

	for flagName, defaultVal := range map[string]string{
		"enable-all":                   "false",
		"enable":                       strings.Join(cfg.EnabledCheckers, ","),
		"expected-actual.pattern":      cfg.ExpectedActual.ExpVarPattern.String(),
		"suite-extra-assert-call.mode": "remove",
	} {
		t.Run(flagName, func(t *testing.T) {
			if v := fs.Lookup(flagName).DefValue; v != defaultVal {
				t.Fatal(v)
			}
		})
	}
}
