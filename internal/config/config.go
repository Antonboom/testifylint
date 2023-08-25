package config

import (
	"flag"

	"github.com/Antonboom/testifylint/internal/checkers"
)

// NewDefault builds default testifylint config.
func NewDefault() Config {
	return Config{
		EnabledCheckers: checkers.EnabledByDefault(),
		ExpectedActual: ExpectedActualConfig{
			ExpVarPattern: RegexpValue{checkers.DefaultExpectedVarPattern},
		},
	}
}

// Config implements testifylint configuration.
type Config struct {
	EnableAll       bool
	EnabledCheckers KnownCheckersValue
	ExpectedActual  ExpectedActualConfig
}

// ExpectedActualConfig implements configuration of checkers.ExpectedActual.
type ExpectedActualConfig struct {
	ExpVarPattern RegexpValue
}

// BindToFlags binds Config fields to according flags.
func BindToFlags(cfg *Config, fs *flag.FlagSet) {
	fs.BoolVar(&cfg.EnableAll, "enable-all", false, "enable all checkers")
	fs.Var(&cfg.EnabledCheckers, "enable", "comma separated list of enabled checkers")
	fs.Var(&cfg.ExpectedActual.ExpVarPattern, "expected-actual.pattern", "regexp for expected variable name")
}
