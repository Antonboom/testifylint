package analyzer_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/analyzer"
	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestTestifyLint_Base(t *testing.T) {
	t.Parallel()

	anlzr := analyzer.New()
	if err := anlzr.Flags.Set("enable", checkers.NewBoolCompare().Name()); err != nil {
		t.Fatal(err)
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, "base-test")
}

func TestTestifyLint_Checkers(t *testing.T) {
	t.Parallel()

	for _, checker := range checkers.All() {
		checker := checker // https://go.dev/wiki/LoopvarExperiment

		t.Run(checker, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.New()
			if err := anlzr.Flags.Set("enable", checker); err != nil {
				t.Fatal(err)
			}
			pkg := filepath.Join("checkers", checker)
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, pkg)
		})
	}
}
