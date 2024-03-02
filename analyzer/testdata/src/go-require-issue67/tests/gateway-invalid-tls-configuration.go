package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewayInvalidTLSConfiguration)
}

var GatewayInvalidTLSConfiguration = suite.ConformanceTest{
	ShortName:   "GatewayInvalidTLSConfiguration",
	Description: "A Gateway should fail to become ready if the Gateway has an invalid TLS configuration",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
	},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		testCases := []struct {
			name                  string
			gatewayNamespacedName NamespacedName
		}{
			{
				name:                  "Nonexistent secret referenced as CertificateRef in a Gateway listener",
				gatewayNamespacedName: NamespacedName{Name: "gateway-certificate-nonexistent-secret", Namespace: "gateway-conformance-infra"},
			},
			{
				name:                  "Unsupported group resource referenced as CertificateRef in a Gateway listener",
				gatewayNamespacedName: NamespacedName{Name: "gateway-certificate-unsupported-group", Namespace: "gateway-conformance-infra"},
			},
			{
				name:                  "Unsupported kind resource referenced as CertificateRef in a Gateway listener",
				gatewayNamespacedName: NamespacedName{Name: "gateway-certificate-unsupported-kind", Namespace: "gateway-conformance-infra"},
			},
			{
				name:                  "Malformed secret referenced as CertificateRef in a Gateway listener",
				gatewayNamespacedName: NamespacedName{Name: "gateway-certificate-malformed-secret", Namespace: "gateway-conformance-infra"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// ...
			})
		}
	},
}

type NamespacedName struct {
	Namespace string
	Name      string
}
