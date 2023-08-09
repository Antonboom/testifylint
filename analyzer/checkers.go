package analyzer

import (
	"fmt"
	"github.com/Antonboom/testifylint/config"
	"github.com/Antonboom/testifylint/internal/checkers"
)

// newCheckers accepts linter config and returns slices of enabled checkers.
func newCheckers(cfg config.Config) ([]checkers.CallChecker, []checkers.AdvancedChecker, error) {
	enabledCheckers := cfg.EnabledCheckers
	if len(enabledCheckers) == 0 {
		enabledCheckers = checkers.EnabledByDefault()
	}

	callCheckers := make([]checkers.CallChecker, 0, len(enabledCheckers))
	advancedCheckers := make([]checkers.AdvancedChecker, 0, len(enabledCheckers))

	for _, name := range enabledCheckers {
		ch, ok := checkers.Get(name)
		if !ok {
			return nil, nil, fmt.Errorf("unknown checker %q", name)
		}

		switch c := ch.(type) {
		case *checkers.ExpectedActual:
			c.SetExpPattern(cfg.ExpectedActual.Pattern.Regexp)
		}

		switch casted := ch.(type) {
		case checkers.CallChecker:
			callCheckers = append(callCheckers, casted)
		case checkers.AdvancedChecker:
			advancedCheckers = append(advancedCheckers, casted)
		}
	}

	checkers.SortByPriority(callCheckers)
	checkers.SortByPriority(advancedCheckers)

	return callCheckers, advancedCheckers, nil
}
