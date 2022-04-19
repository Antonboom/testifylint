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
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), pkgs...)
}
