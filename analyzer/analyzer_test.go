package analyzer_test

import (
	"github.com/Antonboom/testifylint/internal/checkers"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/analyzer"
	"github.com/Antonboom/testifylint/config"
)

func TestTestifyLint(t *testing.T) {
	for _, checker := range checkers.All() {
		checker := checker // https://go.dev/wiki/LoopvarExperiment

		t.Run(checker, func(t *testing.T) {
			t.Parallel()

			cfg := config.Config{EnabledCheckers: []string{checker}}
			pkg := filepath.Join("checkers", checker)
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(cfg), pkg)
		})
	}
}
