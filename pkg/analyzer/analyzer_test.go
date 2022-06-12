package analyzer_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/pkg/analyzer"
	"github.com/Antonboom/testifylint/pkg/config"
)

func TestTestifyLint(t *testing.T) {
	cfg := config.Config{
		Checkers: config.CheckersConfig{
			Disable: []string{
				"compares",
				"require-error",
			},
		},
	}
	pkgs := []string{
		filepath.Join("checkers", "most-of"),
		"negative",
		// "pkg-alias",
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), pkgs...)
}

func TestTestifyLint_SeparateCheckers(t *testing.T) {
	checkers := []string{
		"compares",
		"require-error",
		"suite-no-extra-assert-call",
		"suite-thelper",
	}

	for _, checker := range checkers {
		t.Run(checker, func(t *testing.T) {
			cfg := config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: true,
					Enable:     []string{checker},
				},
			}
			path := filepath.Join("checkers", checker)
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), path)
		})
	}
}
