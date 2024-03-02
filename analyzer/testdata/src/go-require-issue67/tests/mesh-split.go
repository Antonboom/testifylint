package tests

import (
	"fmt"
	"testing"

	"go-require-issue67/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, MeshTrafficSplit)
}

var MeshTrafficSplit = suite.ConformanceTest{
	ShortName:   "MeshTrafficSplit",
	Description: "A mesh client can send traffic to a Service which is split between two versions",
	Features:    []suite.SupportedFeature{},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		cases := []ExpectedResponse{
			{
				Request: Request{
					Host:   "echo",
					Method: "GET",
					Path:   "/v1",
				},
				Response: Response{
					StatusCode: 200,
				},
				Backend: "echo-v1",
			},
			{
				Request: Request{
					Host:   "echo",
					Method: "GET",
					Path:   "/v2",
				},
				Response: Response{
					StatusCode: 200,
				},
				Backend: "echo-v2",
			},
		}
		for i := range cases {
			// Declare tc here to avoid loop variable
			// reuse issues across parallel tests.
			tc := cases[i]
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				// ...
			})
		}
	},
}

// ExpectedResponse defines the response expected for a given request.
type ExpectedResponse struct {
	// Request defines the request to make.
	Request Request

	// ExpectedRequest defines the request that
	// is expected to arrive at the backend. If
	// not specified, the backend request will be
	// expected to match Request.
	ExpectedRequest *ExpectedRequest

	// BackendSetResponseHeaders is a set of headers
	// the echoserver should set in its response.
	BackendSetResponseHeaders map[string]string

	// Response defines what response the test case
	// should receive.
	Response Response

	Backend   string
	Namespace string

	// User Given TestCase name
	TestCaseName string
}

// GetTestCaseName gets the user-defined test case name or generates one from expected response to a given request.
func (er *ExpectedResponse) GetTestCaseName(i int) string {
	// If TestCase name is provided then use that or else generate one.
	if er.TestCaseName != "" {
		return er.TestCaseName
	}

	headerStr := ""
	reqStr := ""

	if er.Request.Headers != nil {
		headerStr = " with headers"
	}

	reqStr = fmt.Sprintf("%d request to '%s%s'%s", i, er.Request.Host, er.Request.Path, headerStr)

	if er.Backend != "" {
		return fmt.Sprintf("%s should go to %s", reqStr, er.Backend)
	}
	return fmt.Sprintf("%s should receive a %d", reqStr, er.Response.StatusCode)
}

// Request can be used as both the request to make and a means to verify
// that echoserver received the expected request. Note that multiple header
// values can be provided, as a comma-separated value.
type Request struct {
	Host             string
	Method           string
	Path             string
	Headers          map[string]string
	UnfollowRedirect bool
	Protocol         string
}

// ExpectedRequest defines expected properties of a request that reaches a backend.
type ExpectedRequest struct {
	Request

	// AbsentHeaders are names of headers that are expected
	// *not* to be present on the request.
	AbsentHeaders []string
}

// Response defines expected properties of a response from a backend.
type Response struct {
	StatusCode    int
	Headers       map[string]string
	AbsentHeaders []string
}
