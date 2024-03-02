package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewaySecretInvalidReferenceGrant)
}

var GatewaySecretInvalidReferenceGrant = suite.ConformanceTest{
	ShortName:   "GatewaySecretInvalidReferenceGrant",
	Description: "A Gateway in the gateway-conformance-infra namespace should fail to become ready if the Gateway has a certificateRef for a Secret in the gateway-conformance-web-backend namespace and a ReferenceGrant exists but does not grant permission to that specific Secret",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
		suite.SupportReferenceGrant,
	},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		_ = NamespacedName{Name: "gateway-secret-invalid-reference-grant", Namespace: "gateway-conformance-infra"}

		t.Run("Gateway listener should have a false ResolvedRefs condition with reason RefNotPermitted", func(t *testing.T) {
			// ...
		})
	},
}
