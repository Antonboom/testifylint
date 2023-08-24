package analyzer

import (
	"fmt"

	"github.com/Antonboom/testifylint/config"
	"github.com/Antonboom/testifylint/internal/checkers"
)

// newCheckers accepts linter config and returns slices of enabled checkers.
func newCheckers(cfg config.Config) ([]checkers.RegularChecker, []checkers.AdvancedChecker, error) {
	enabledCheckers := cfg.EnabledCheckers
	if len(enabledCheckers) == 0 {
		enabledCheckers = checkers.EnabledByDefault()
	}

	checkers.SortByPriority(enabledCheckers)

	regularCheckers := make([]checkers.RegularChecker, 0, len(enabledCheckers))
	advancedCheckers := make([]checkers.AdvancedChecker, 0, len(enabledCheckers)/2)

	for _, name := range enabledCheckers {
		ch, ok := checkers.Get(name)
		if !ok {
			return nil, nil, fmt.Errorf("unknown checker %q", name)
		}

		switch c := ch.(type) {
		case *checkers.ExpectedActual:
			c.SetExpVarPattern(cfg.ExpectedActual.ExpVarPattern.Regexp)
		}

		switch casted := ch.(type) {
		case checkers.RegularChecker:
			regularCheckers = append(regularCheckers, casted)
		case checkers.AdvancedChecker:
			advancedCheckers = append(advancedCheckers, casted)
		}
	}

	return regularCheckers, advancedCheckers, nil
}
