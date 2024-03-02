package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewaySecretMissingReferenceGrant)
}

var GatewaySecretMissingReferenceGrant = suite.ConformanceTest{
	ShortName:   "GatewaySecretMissingReferenceGrant",
	Description: "A Gateway in the gateway-conformance-infra namespace should fail to become programmed if the Gateway has a certificateRef for a Secret in the gateway-conformance-web-backend namespace and a ReferenceGrant granting permission to the Secret does not exist",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
		suite.SupportReferenceGrant,
	},
	Manifests: []string{"tests/gateway-secret-missing-reference-grant.yaml"},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		_ = NamespacedName{Name: "gateway-secret-missing-reference-grant", Namespace: "gateway-conformance-infra"}

		t.Run("Gateway listener should have a false ResolvedRefs condition with reason RefNotPermitted", func(t *testing.T) {
			// ...
		})
	},
}
