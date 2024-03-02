package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, HTTPRouteReferenceGrant)
}

var HTTPRouteReferenceGrant = suite.ConformanceTest{
	ShortName:   "HTTPRouteReferenceGrant",
	Description: "A single HTTPRoute in the gateway-conformance-infra namespace, with a backendRef in the gateway-conformance-web-backend namespace, should attach to Gateway in the gateway-conformance-infra namespace",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
	},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		t.Run("Simple HTTP request should reach web-backend", func(t *testing.T) {
			var err error
			require.NoError(t, err)
		})

		_, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var err error
		require.NoError(t, err)

		t.Run("Simple HTTP request should return 500 after deleting the relevant reference grant", func(t *testing.T) {
			var err error
			require.NoError(t, err)
		})
	},
}
