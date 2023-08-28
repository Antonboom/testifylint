// Test is inspired by
// https://github.com/kubernetes/ingress-nginx/blob/300a60f7244cb056701177989a1a09135626f930/test/e2e/loadbalance/ewma.go#L52
package ginkgo

import (
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("[Load Balancer] EWMA", func() {

	ginkgo.BeforeEach(func() {
		// Some stuff...
	})

	ginkgo.It("does not fail requests", func() {
		host := "load-balance.com"

		WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, "server_name load-balance.com")
			})

		algorithm, err := GetLbAlgorithm(EchoService, 80)
		assert.Nil(ginkgo.GinkgoT(), err)                 // want "error-nil: use assert.NoError"
		assert.Equal(ginkgo.GinkgoT(), algorithm, "ewma") // want "expected-actual: need to reverse actual and expected values"
	})
})

// EchoService name of the deployment for the echo app
const EchoService = "echo"

// WaitForNginxServer waits until the nginx configuration contains a particular server section.
// `cfg` passed to matcher is normalized by replacing all tabs and spaces with single space.
func WaitForNginxServer(name string, matcher func(cfg string) bool) {}

// GetLbAlgorithm returns algorithm identifier for the given backend
func GetLbAlgorithm(serviceName string, servicePort int) (string, error) {
	return "", nil
}
