package analyzer

import (
	"regexp"
	"sort"

	"github.com/Antonboom/testifylint/internal/checker"
	"github.com/Antonboom/testifylint/pkg/config"
)

// newCheckers accepts validated config and returns slice of enabled checkers.
func newCheckers(cfg config.Config) []checker.Checker {
	enabledByDefault := checker.EnabledByDefaultCheckers()
	result := make([]checker.Checker, 0, len(enabledByDefault))

	enabledCheckers := make(map[string]struct{}, len(enabledByDefault))
	for _, ch := range enabledByDefault {
		enabledCheckers[ch] = struct{}{}
	}

	if cfg.Checkers.DisableAll {
		enabledCheckers = make(map[string]struct{})
	}

	for _, ch := range cfg.Checkers.Enable {
		enabledCheckers[ch] = struct{}{}
	}
	for _, ch := range cfg.Checkers.Disable {
		delete(enabledCheckers, ch)
	}

	for name := range enabledCheckers {
		ch, ok := checker.Get(name)
		if !ok {
			panic("unregistered checker: " + name)
		}

		switch c := ch.(type) {
		case *checker.ExpectedActual:
			if p := cfg.ExpectedActual.Pattern; p != "" {
				c.SetExpPattern(regexp.MustCompile(cfg.ExpectedActual.Pattern))
			}
		}
		result = append(result, ch)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority() < result[j].Priority()
	})
	return result
}
