package conformance_test

import (
	"testing"

	"go-require-issue67/suite"
	"go-require-issue67/tests"
)

func TestConformance(t *testing.T) {
	cSuite := new(suite.ConformanceTestSuite)
	cSuite.Run(t, tests.ConformanceTests)
}
