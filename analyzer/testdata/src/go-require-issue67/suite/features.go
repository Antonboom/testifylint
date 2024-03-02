package suite

// SupportedFeature allows opting in to additional conformance tests at an
// individual feature granularity.
type SupportedFeature string

const (
	// This option indicates support for Gateway.
	// Opting out of this is allowed only for GAMMA-only implementations
	SupportGateway SupportedFeature = "Gateway"
)

const (
	// This option indicates that the Gateway can also use port 8080
	SupportGatewayPort8080 SupportedFeature = "GatewayPort8080"

	// SupportGatewayStaticAddresses option indicates that the Gateway is capable
	// of allocating pre-determined addresses, rather than dynamically having
	// addresses allocated for it.
	SupportGatewayStaticAddresses SupportedFeature = "GatewayStaticAddresses"
)

const (
	// This option indicates support for ReferenceGrant.
	SupportReferenceGrant SupportedFeature = "ReferenceGrant"
)
