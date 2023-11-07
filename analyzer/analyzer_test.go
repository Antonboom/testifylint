package analyzer_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/analyzer"
	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestTestifyLint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir   string
		flags map[string]string
	}{
		{
			dir:   "base-test",
			flags: map[string]string{"disable-all": "true", "enable": checkers.NewBoolCompare().Name()},
		},
		{
			dir:   "checkers-priority",
			flags: map[string]string{"enable-all": "true"},
		},
		{
			dir:   "error-as-target",
			flags: map[string]string{"disable-all": "true", "enable": checkers.NewErrorIsAs().Name()},
		},
		{
			dir: "expected-var-custom-pattern",
			flags: map[string]string{
				"disable-all":             "true",
				"enable":                  checkers.NewExpectedActual().Name(),
				"expected-actual.pattern": "goldenValue",
			},
		},
		{dir: "ginkgo"},
		{
			dir:   "not-std-funcs",
			flags: map[string]string{"enable-all": "true"},
		},
		{dir: "not-test-file"},    // By default, linter checks regular files too.
		{dir: "not-true-testify"}, // Linter ignores stretchr/testify's forks.
		{dir: "pkg-alias"},
		{
			dir: "require-error-fn-pattern",
			flags: map[string]string{
				"disable-all":              "true",
				"enable":                   checkers.NewRequireError().Name(),
				"require-error.fn-pattern": "^(NoErrorf?|NotErrorIsf?)$",
			},
		},
		{
			dir:   "require-error-skip-logic",
			flags: map[string]string{"disable-all": "true", "enable": checkers.NewRequireError().Name()},
		},
		{
			dir: "suite-require-extra-assert-call",
			flags: map[string]string{
				"disable-all":                  "true",
				"enable":                       checkers.NewSuiteExtraAssertCall().Name(),
				"suite-extra-assert-call.mode": "require",
			},
		},
	}

	for _, tt := range cases {
		tt := tt

		t.Run(tt.dir, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.New()
			for k, v := range tt.flags {
				if err := anlzr.Flags.Set(k, v); err != nil {
					t.Fatal(err)
				}
			}
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, tt.dir)
		})
	}
}

func TestTestifyLint_CheckersDefault(t *testing.T) {
	t.Parallel()

	for _, checker := range checkers.All() {
		checker := checker

		t.Run(checker, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.New()
			if err := anlzr.Flags.Set("disable-all", "true"); err != nil {
				t.Fatal(err)
			}
			if err := anlzr.Flags.Set("enable", checker); err != nil {
				t.Fatal(err)
			}

			pkg := filepath.Join("checkers-default", checker)
			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, pkg)
		})
	}
}
