package analyzer_test

import (
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/analyzer"
	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestTestifyLint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir             string
		enabledCheckers []string
	}{
		{dir: "base-test", enabledCheckers: []string{checkers.NewBoolCompare().Name()}},
		{dir: "ginkgo"},
		{dir: "pkg-alias"},
	}

	for _, tt := range cases {
		tt := tt

		t.Run(tt.dir, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.New()
			if len(tt.enabledCheckers) > 0 {
				if err := anlzr.Flags.Set("enable", strings.Join(tt.enabledCheckers, ",")); err != nil {
					t.Fatal(err)
				}
			}
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, tt.dir)
		})
	}
}

func TestTestifyLint_Checkers(t *testing.T) {
	t.Parallel()

	for _, checker := range checkers.All() {
		checker := checker

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
