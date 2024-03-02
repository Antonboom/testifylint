package tests

import (
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, TLSRouteInvalidReferenceGrant)
}

var TLSRouteInvalidReferenceGrant = suite.ConformanceTest{
	ShortName: "TLSRouteInvalidReferenceGrant",
	Description: "A single TLSRoute in the gateway-conformance-infra namespace, with a backendRef in another namespace" +
		" without valid ReferenceGrant, should have the ResolvedRefs condition set to False",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
	},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		t.Run("TLSRoute with BackendRef in another namespace and no ReferenceGrant covering the Service has a "+
			"ResolvedRefs Condition with status False and Reason RefNotPermitted", func(t *testing.T) {
		})
	},
}
