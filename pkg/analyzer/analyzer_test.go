package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/pkg/analyzer"
	"github.com/Antonboom/testifylint/pkg/config"
)

func TestTestifyLint(t *testing.T) {
	cfg := config.Config{
		Checkers: config.CheckersConfig{
			Disable: []string{"require-error"},
		},
	}
	pkgs := []string{
		"basic",
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), pkgs...)
}

func TestTestifyLint_SpecificCheckers(t *testing.T) {
	checkers := []string{
		"require-error",
		"suite-no-extra-assert-call",
	}

	for _, checker := range checkers {
		t.Run(checker, func(t *testing.T) {
			cfg := config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: true,
					Enable:     []string{checker},
				},
			}
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), checker)
		})
	}
}
