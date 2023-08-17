package main

import (
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestCheckerName_AsPkgName(t *testing.T) {
	name := checkers.NewSuiteNoExtraAssertCall().Name()
	testName := CheckerName(name).AsPkgName()
	if testName != "suitenoextraassertcall" {
		t.Fatal(testName)
	}
}

func TestCheckerName_AsTestName(t *testing.T) {
	name := checkers.NewSuiteNoExtraAssertCall().Name()
	testName := CheckerName(name).AsTestName()
	if testName != "TestSuiteNoExtraAssertCallChecker" {
		t.Fatal(testName)
	}
}
