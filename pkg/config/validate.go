package config

import (
	"fmt"
	"regexp"

	"github.com/Antonboom/testifylint/internal/checker"
)

func Validate(cfg Config) error {
	if err := validateCheckersConf(cfg.Checkers); err != nil {
		return err
	}
	if err := validateExpectedActualConf(cfg.ExpectedActual); err != nil {
		return err
	}
	return nil
}

func validateCheckersConf(cfg CheckersConfig) error {
	if cfg.DisableAll {
		if len(cfg.Enable) == 0 {
			return fmt.Errorf("checkers.disable-all is true, but no one checker was enabled")
		}

		if len(cfg.Disable) != 0 {
			return fmt.Errorf("checkers.disable-all and checkers.disable options must not be combined")
		}
	}

	for _, checkerName := range cfg.Enable {
		if ok := checker.IsKnown(checkerName); !ok {
			return fmt.Errorf("checkers.enable: unknown checker %q", checkerName)
		}
	}

	for _, checkerName := range cfg.Disable {
		if ok := checker.IsKnown(checkerName); !ok {
			return fmt.Errorf("checkers.disable: unknown checker %v", checkerName)
		}
	}

	return nil
}

func validateExpectedActualConf(cfg ExpectedActualConfig) error {
	if p := cfg.Pattern; p != "" {
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("expected-actual.pattern: %v", err)
		}
	}
	return nil
}
