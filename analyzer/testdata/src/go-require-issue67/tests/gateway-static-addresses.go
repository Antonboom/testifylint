package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, GatewayStaticAddresses)
}

// GatewayStaticAddresses tests the implementation's support of deploying
// Gateway resources with static addresses, or in other words addresses
// provided via the specification rather than relying on the underlying
// implementation/network to dynamically assign the Gateway an address.
//
// Running this test against your own implementation is currently a little bit
// messy, as at the time of writing we didn't have great ways to provide the
// test suite with things like known good, or known bad addresses to run the
// test with (as we obviously can't determine that for the implementation).
//
// As such, if you're trying to enable this test for yourself and you're getting
// confused about how to provide addresses, you'll actually do that in the
// conformance test suite BEFORE you even set up and run your tests. Make sure
// you populate the following test suite fields:
//
//   - suite.UsableNetworkAddresses
//   - suite.UnusableNetworkAddresses
//
// With appropriate network addresses for your network environment.
var GatewayStaticAddresses = suite.ConformanceTest{
	ShortName:   "GatewayStaticAddresses",
	Description: "A Gateway in the gateway-conformance-infra namespace should be able to use previously determined addresses.",
	Features: []suite.SupportedFeature{
		suite.SupportGateway,
		suite.SupportGatewayStaticAddresses,
	},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		_, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var err error
		require.NoError(t, err, "error getting Gateway: %v", err)

		var err2 error
		require.NoError(t, err2, "failed to patch Gateway: %v", err)

		t.Logf("verifying that the Gateway %s/%s is now accepted, but is not programmed due to an address that can't be used")
		var err3 error
		require.NoError(t, err3, "error getting Gateway: %v", err)

		t.Logf("patching Gateway %s/%s to remove the unusable address %s")
		var err4 error
		require.NoError(t, err4, "failed to patch Gateway: %v", err)

		t.Logf("verifying that the Gateway %s/%s is accepted and programmed with the usable static address %s assigned")
		require.Equal(t, "usableAddress.Type", "currentGW.Status.Addresses[0].Type",
			"expected address type to match the usable address")
		require.Equal(t, "usableAddress.Value", "currentGW.Status.Addresses[0].Value",
			"expected usable address to be assigned")
	},
}
