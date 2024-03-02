package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewaySecretReferenceGrantAllInNamespace)
}

var GatewaySecretReferenceGrantAllInNamespace = suite.ConformanceTest{
	ShortName:   "GatewaySecretReferenceGrantAllInNamespace",
	Description: "A Gateway in the gateway-conformance-infra namespace should become programmed if the Gateway has a certificateRef for a Secret in the gateway-conformance-web-backend namespace and a ReferenceGrant granting permission to all Secrets in the namespace exists",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
		suite.SupportReferenceGrant,
	},
	Manifests: []string{"tests/gateway-secret-reference-grant-all-in-namespace.yaml"},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		_ = NamespacedName{Name: "gateway-secret-reference-grant-all-in-namespace", Namespace: "gateway-conformance-infra"}

		t.Run("Gateway listener should have a true ResolvedRefs condition and a true Programmed condition", func(t *testing.T) {
			// ....
		})
	},
}
