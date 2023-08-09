package config

import (
	"fmt"

	"github.com/Antonboom/testifylint/internal/checkers"
)

func Validate(cfg Config) error {
	return validateEnabledCheckers(cfg.EnabledCheckers)
}

func validateEnabledCheckers(cfgCheckers []string) error {
	for _, checkerName := range cfgCheckers {
		if ok := checkers.IsKnown(checkerName); !ok {
			return fmt.Errorf("unknown checker %q", checkerName)
		}
	}
	return nil
}
