package requireerrorskiplogic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSomeServer(t *testing.T) {
	httptest.NewServer(http.HandlerFunc(func(hres http.ResponseWriter, hreq *http.Request) {
		var req MyRequest
		err := json.NewDecoder(hreq.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "42", req.ID)
	}))

	httptest.NewServer(http.HandlerFunc(func(hres http.ResponseWriter, hreq *http.Request) {
		var req MyRequest
		err := json.NewDecoder(hreq.Body).Decode(&req)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "42", req.ID)
	}))

	httptest.NewServer(&handler{t: t})
	httptest.NewServer(http.HandlerFunc((&handler{t: t}).Handle))
	httptest.NewServer(handlerClosure(t))
}

var _ http.Handler = (*handler)(nil)

type handler struct {
	t *testing.T
}

func (h *handler) ServeHTTP(hres http.ResponseWriter, hreq *http.Request) {
	var req MyRequest
	err := json.NewDecoder(hreq.Body).Decode(&req)
	assert.NoError(h.t, err)
	assert.Equal(h.t, "42", req.ID)
}

func (h *handler) Handle(hres http.ResponseWriter, hreq *http.Request) {
	var req MyRequest
	err := json.NewDecoder(hreq.Body).Decode(&req)
	assert.NoError(h.t, err)
	assert.Equal(h.t, "42", req.ID)
}

func handlerClosure(t *testing.T) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(hres http.ResponseWriter, hreq *http.Request) {
		var req MyRequest
		err := json.NewDecoder(hreq.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "42", req.ID)
	})
}

type SomeServerSuite struct {
	suite.Suite
}

func TestSomeServerSuite(t *testing.T) {
	suite.Run(t, &SomeServerSuite{})
}

func (s *SomeServerSuite) TestServer() {
	httptest.NewServer(http.HandlerFunc(s.handler))
	httptest.NewServer(s)
}

func (s *SomeServerSuite) ServeHTTP(hres http.ResponseWriter, hreq *http.Request) {
	var req MyRequest
	err := json.NewDecoder(hreq.Body).Decode(&req)
	s.Require().NoError(err)
	s.Equal("42", req.ID)
}

func (s *SomeServerSuite) handler(hres http.ResponseWriter, hreq *http.Request) {
	var req MyRequest
	err := json.NewDecoder(hreq.Body).Decode(&req)
	s.Require().NoError(err)
	s.Equal("42", req.ID)
}

type MyRequest struct {
	ID string
}
