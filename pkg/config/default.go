package config

import "github.com/Antonboom/testifylint/internal/checker"

var Default = Config{
	Checkers: CheckersConfig{
		DisableAll: false,
		Enable:     checker.EnabledByDefaultCheckers(),
		Disable:    checker.DisabledByDefaultCheckers(),
	},
	ExpectedActual: ExpectedActualConfig{
		Pattern: checker.DefaultExpectedVarPattern.String(),
	},
}
