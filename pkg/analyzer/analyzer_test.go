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

func TestTestifyLint_RequireError(t *testing.T) {
	cfg := config.Config{
		Checkers: config.CheckersConfig{
			DisableAll: true,
			Enable:     []string{"require-error"},
		},
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), "require-error")
}
