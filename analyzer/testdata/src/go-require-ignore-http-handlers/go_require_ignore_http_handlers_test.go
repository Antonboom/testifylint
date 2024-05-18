package gorequireignorehttphandlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestServer_Require(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open("some file.json")
		require.NoError(t, err)

		data, err := io.ReadAll(file)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
	}))
	defer ts.Close()

	client := ts.Client()
	client.Timeout = 10 * time.Second

	req, err := http.NewRequest("GET", ts.URL+"/require", nil)
	require.NoError(t, err)

	statusCode := make(chan int)
	go func() {
		resp, err := client.Do(req)
		require.NoError(t, err) // want "go-require: require must only be used in the goroutine running the test function"
		defer func() {
			require.NoError(t, resp.Body.Close()) // want "go-require: require must only be used in the goroutine running the test function"
		}()
		statusCode <- resp.StatusCode
	}()

	require.Equal(t, http.StatusOK, <-statusCode)
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
