package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Antonboom/testifylint/pkg/analyzer"
)

func TestTestifyLint(t *testing.T) {
	pkgs := []string{
		"basic",
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(), pkgs...)
}

// TODO: тесты флагов
// https://github.com/mweb/floatcompare/blob/main/analyzer_test.go
