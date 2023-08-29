package main

import (
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestCheckerName(t *testing.T) {
	cases := []struct {
		name        string
		transformer func(CheckerName) string
		expected    string
	}{
		{
			name:        checkers.NewSuiteExtraAssertCall().Name(),
			transformer: CheckerName.AsPkgName,
			expected:    "suiteextraassertcall",
		},
		{
			name:        checkers.NewSuiteDontUsePkg().Name(),
			transformer: CheckerName.AsTestName,
			expected:    "TestSuiteDontUsePkgChecker",
		},
		{
			name:        checkers.NewFloatCompare().Name(),
			transformer: CheckerName.AsSuiteName,
			expected:    "FloatCompareCheckerSuite",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if n := tt.transformer(CheckerName(tt.name)); n != tt.expected {
				t.Fatal(n)
			}
		})
	}
}
