package gorequireissue66issue73_test

import (
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

// NOTE(a.telyshev): Neither `assert` nor `require` is the best way to test an HTTP handler:
// it leads to redundant stack traces (up to runtime assembler), as well as undefined behaviour (in `require` case).
// Use HTTP mechanisms (status code, headers, response data) and place assertions in the main test function.

func TestServer_Assert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open("some file.json")
		if !assert.NoError(t, err) {
			return
		}

		data, err := io.ReadAll(file)
		if !assert.NoError(t, err) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client := ts.Client()

	req, err := http.NewRequest("GET", ts.URL+"/assert", nil)
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServer_Require(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open("some file.json")
		require.NoError(t, err) // want "go-require: do not use require in http handlers"

		data, err := io.ReadAll(file)
		require.NoError(t, err) // want "go-require: do not use require in http handlers"

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error()) // want "go-require: do not use assert\\.FailNow in http handlers"
		}
	}))
	defer ts.Close()

	client := ts.Client()
	client.Timeout = 10 * time.Second

	req, err := http.NewRequest("GET", ts.URL+"/require", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

type ServerSuite struct {
	suite.Suite
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, &ServerSuite{})
}

func (s *ServerSuite) TestServer() {
	httptest.NewServer(http.HandlerFunc(s.handler))
}

func (s *ServerSuite) handler(w http.ResponseWriter, _ *http.Request) {
	s.T().Helper()

	file, err := os.Open("some file.json")
	s.Require().NoError(err) // want "go-require: do not use require in http handlers"

	data, err := io.ReadAll(file)
	if !s.NoError(err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if !s.NoError(err) {
		s.FailNow(err.Error()) // want "go-require: do not use s\\.FailNow in http handlers"
	}
}
