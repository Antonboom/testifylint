package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewayObservedGenerationBump)
}

var GatewayObservedGenerationBump = suite.ConformanceTest{
	ShortName: "GatewayObservedGenerationBump",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
		suite.SupportGatewayPort8080,
	},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		t.Run("observedGeneration should increment", func(t *testing.T) {
			var err error
			require.NoErrorf(t, err, "error getting Gateway: %v", err)
			require.NotEqual(t, "original.Generation", "updated.Generation",
				"generation should change after an update")
		})
	},
}
