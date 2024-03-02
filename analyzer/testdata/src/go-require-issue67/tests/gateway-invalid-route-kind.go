package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewayInvalidRouteKind)
}

var GatewayInvalidRouteKind = suite.ConformanceTest{
	ShortName:   "GatewayInvalidRouteKind",
	Description: "A Gateway in the gateway-conformance-infra namespace should fail to become ready an invalid Route kind is specified.",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
	},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		t.Run("Gateway listener should have a false ResolvedRefs condition with reason InvalidRouteKinds and no supportedKinds", func(t *testing.T) {
			// ...
		})

		t.Run("Gateway listener should have a false ResolvedRefs condition with reason InvalidRouteKinds and HTTPRoute must be put in the supportedKinds", func(t *testing.T) {
			// ...
		})
	},
}
