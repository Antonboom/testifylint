package requireerrorskiplogic

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestartCountByLogDir(t *testing.T) {
	for _, tc := range []struct {
		filenames    []string
		restartCount int
	}{
		{
			filenames:    []string{"0.log.rotated-log"},
			restartCount: 1,
		},
	} {
		t.Run("", func(t *testing.T) {
			count, err := calcRestartCountByLogDir(tc.filenames)
			if assert.NoError(t, err) {
				assert.Equal(t, count, tc.restartCount)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			}

			if noErr := assert.NoError(t, err); noErr {
				assert.Error(t, err)
				assert.Equal(t, count, tc.restartCount)
			} else {
				assert.Equal(t, count, tc.restartCount)
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	foo := "hello world"
	if b, err := json.Marshal(foo); assert.NoError(t, err) {
		assert.Equal(t, "hello world", b)
	}
}

func TestNewStatusError(t *testing.T) {
	err := NewStatusError("oh no", StatusNotFound)
	assert.EqualError(t, err, "status error: oh no") // want "require-error: for error assertions use require"

	var statusErr StatusError
	if assert.ErrorAs(t, err, &statusErr) {
		assert.Equal(t, StatusNotFound, statusErr.StatusID())
		assert.Error(t, statusErr)
	} else {
		assert.NoError(t, statusErr)
		assert.Equal(t, StatusNotFound, statusErr.StatusID())
	}

	err = NewStatusError("oh no", StatusNotFound)
	assert.EqualError(t, err, "status error: oh no")
}

func TestGetDocFileWithLongLine(t *testing.T) {
	fpath := filepath.Join("testdata", "autogen_exclude_long_line.go")
	_, err := getDoc(fpath)
	assert.NoError(t, err)
}

func TestIdentifierMarker(t *testing.T) {
	cases := []struct{ in, out string }{
		{"unknown field Address in struct literal", "unknown field `Address` in struct literal"},
	}
	p := NewIdentifierMarker()

	for _, c := range cases {
		out, err := p.Process([]Issue{{Text: c.in}})
		assert.NoError(t, err) // want "require-error: for error assertions use require"
		assert.Equal(t, []Issue{{Text: c.out}}, out)
	}
}

func TestValidateAppArmorProfileFormat(t *testing.T) {
	tests := []struct {
		profile     string
		expectValid bool
	}{
		{"", true},
		{"baz", false},
	}

	for _, test := range tests {
		err := ValidateAppArmorProfileFormat(test.profile)
		if test.expectValid {
			assert.NoError(t, err, "Profile %s should be valid", test.profile) // want "require-error: for error assertions use require"
		} else {
			assert.Error(t, err, "Profile %s should not be valid", test.profile) // want "require-error: for error assertions use require"
		}

		t.Run("", func(t *testing.T) {
			if test.expectValid {
				assert.NoError(t, err, "Profile %s should be valid", test.profile)
			} else {
				assert.Error(t, err, "Profile %s should not be valid", test.profile)
			}
		})
	}
}

func TestGetPortForward(t *testing.T) {
	const (
		podName             = "podFoo"
		podNamespace        = "nsFoo"
		podUID       string = "12345678"
	)

	testcases := []struct {
		description string
		podName     string
		expectError bool
	}{{
		description: "success case",
		podName:     podName,
	}, {
		description: "no such pod",
		podName:     "bar",
		expectError: true,
	}}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			description := "streaming - " + tc.description
			ctx := context.Background()
			redirect, err := GetPortForward(ctx, tc.podName, podNamespace, podUID)
			if tc.expectError {
				assert.Error(t, err, description)
			} else {
				assert.NoError(t, err, description) // want "require-error: for error assertions use require"
				assert.Equal(t, "localhost:12345", redirect.Host, description+": redirect")
				assert.NoError(t, err, description)
			}
		})
	}
}

func TestDriverParameter(t *testing.T) {
	testcases := []struct {
		name     string
		filename string
		err      string
		expected string
	}{
		{
			name:     "no such file",
			filename: "no-such-file.yaml",
			err:      "open no-such-file.yaml: no such file or directory",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actual, err := loadDriverDefinition(testcase.filename)
			if testcase.err == "" {
				assert.NoError(t, err, testcase.name) // want "require-error: for error assertions use require"
			} else {
				if assert.Error(t, err, testcase.name) {
					assert.Equal(t, testcase.err, err.Error())
					assert.Error(t, err, testcase.name)
				}
			}
			if err == nil {
				assert.Equal(t, testcase.expected, actual)
			}
		})
	}
}

func TestRunGenCSR(t *testing.T) {
	kubeConfigDir := "kubernetes"
	certDir := kubeConfigDir + "/pki"

	expectedCertificates := []string{
		"apiserver",
		"apiserver-etcd-client",
	}

	err := runGenCSR(nil, nil)
	require.NoError(t, err, "expected runGenCSR to not fail")

	for _, name := range expectedCertificates {
		_, err = TryLoadKeyFromDisk(certDir, name)
		assert.NoErrorf(t, err, "failed to load key file: %s", name) // want "require-error: for error assertions use require"

		_, err = TryLoadCSRFromDisk(certDir, name)
		assert.NoError(t, err, "failed to load CSR file: %s", name)
	}
}

func TestDecode_Errors(t *testing.T) {
	t.Run("invalid base64", func(t *testing.T) {
		err := Decode(`{"page_size":50,"last":1670502502}`, new(Cursor))
		assert.Error(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		err := Decode("eyJwYWdlX3NpemUiOjUwLCJsYXN0IjoxNjcwNTAyNTAy", new(Cursor))
		assert.Error(t, err)
	})
}

func TestRequest_Validate1(t *testing.T) {
	type validater interface{ Validate() error }

	cases := []struct {
		name    string
		request validater
		wantErr bool
	}{
		// Positive.
		{
			name:    "valid request",
			request: struct{ validater }{},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
		assert.NoError(t, tt.request.Validate())
	}
}

func TestRequest_Validate2(t *testing.T) {
	type validater interface{ Validate() error }

	cases := []struct {
		name    string
		request validater
		wantErr bool
	}{
		// Positive.
		{
			name:    "valid request",
			request: struct{ validater }{},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		err := tt.request.Validate()
		if tt.wantErr {
			assert.Error(t, err) // want "require-error: for error assertions use require"
		} else {
			assert.NoError(t, err) // want "require-error: for error assertions use require"
		}
	}

	assert.NoError(t, cases[0].request.Validate())
}

func TestErrors(t *testing.T) {
	cases := []struct {
		name string

		// inputs
		err        error
		statusCode int
		message    string

		// responses
		legacyResponse *NormalResponse
		newResponse    *NormalResponse
		fallbackUseNew bool
		compareErr     bool
	}{
		{
			name: "base case",

			legacyResponse: &NormalResponse{},
			newResponse: &NormalResponse{
				status: http.StatusInternalServerError,
			},
		},
	}

	compareResponses := func(expected *NormalResponse, actual *NormalResponse, compareErr bool) func(t *testing.T) {
		return func(t *testing.T) {
			if expected == nil {
				assert.Nil(t, actual)
				return
			}

			require.NotNil(t, actual)
			assert.Equal(t, expected.status, actual.status)
			if expected.body != nil {
				assert.Equal(t, expected.body.Bytes(), actual.body.Bytes())
			}
			if expected.header != nil {
				assert.EqualValues(t, expected.header, actual.header)
			}
			assert.Equal(t, expected.errMessage, actual.errMessage)
			if compareErr {
				assert.ErrorIs(t, expected.err, actual.err)
			}
		}
	}

	for _, tc := range cases {
		tc := tc
		t.Run(
			tc.name+" Error",
			compareResponses(tc.legacyResponse, Error(
				tc.statusCode,
				tc.message,
				tc.err,
			), tc.compareErr),
		)
	}
}

func PrepareDB(ctx context.Context, t *testing.T, dbName string) (client any, cleanUp func(ctx context.Context)) {
	t.Helper()
	require.NotEmpty(t, dbName)

	_, err := operationWithResult()
	require.NoError(t, err)
	assert.NoError(t, operation()) // want "require-error: for error assertions use require"

	client, err = operationWithResult()
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	{
		err = operation()
	}
	assert.NoError(t, err)

	return client, func(ctx2 context.Context) {
		assert.NoError(t, operation()) // want "require-error: for error assertions use require"
		assert.NoError(t, operation()) // want "require-error: for error assertions use require"
		assert.NoError(t, operation())
	}
}

func TestChatID_Validate(t *testing.T) {
	assert.NoError(t, NewChatID().Validate()) // want "require-error: for error assertions use require"
	assert.Error(t, ChatID{}.Validate())      // want "require-error: for error assertions use require"
	assert.Error(t, ChatIDNil.Validate())
}

func TestCSRDuration(t *testing.T) {
	t.Parallel()

	s := StartTestServerOrDie(t)
	t.Cleanup(s.TearDownFn)

	_, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	t.Cleanup(cancel)

	wantMetricStrings := []string{
		`apiserver_certificates_registry_csr_honored_duration_total{signerName="kubernetes.io/kube-apiserver-client"} 6`,
		`apiserver_certificates_registry_csr_requested_duration_total{signerName="kubernetes.io/kube-apiserver-client"} 7`,
	}
	t.Cleanup(func() {
		_, err := operationWithResult()
		assert.NoError(t, err)

		var body string

		var gotMetricStrings []string
		for _, line := range strings.Split(string(body), "\n") {
			if strings.HasPrefix(line, "apiserver_certificates_registry_") {
				gotMetricStrings = append(gotMetricStrings, line)
			}
		}

		diff := Diff(wantMetricStrings, gotMetricStrings)
		assert.Empty(t, diff)
	})

	_, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	caCert, err := operationWithResult()
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	caPublicKeyFile := path.Join(t.TempDir(), "test-ca-public-key")
	err = os.WriteFile(caPublicKeyFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCert.([]byte)}), os.FileMode(0600))
	require.NoError(t, err)

	caPrivateKeyBytes, err := operationWithResult()
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	caPrivateKeyFile := path.Join(t.TempDir(), "test-ca-private-key")
	err = os.WriteFile(caPrivateKeyFile, caPrivateKeyBytes.([]byte), os.FileMode(0600))
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	tests := []struct {
		name, csrName string
		duration      *time.Duration
		wantDuration  time.Duration
		wantError     string
	}{
		{
			name:         "no duration set",
			duration:     nil,
			wantDuration: 24 * time.Hour,
			wantError:    "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			assert.NoError(t, err) // want "require-error: for error assertions use require"

			csrData, err := MakeCSR(privateKey, &pkix.Name{CommonName: "panda"}, nil, nil)
			require.NoError(t, err)

			csrName, errReq := operationWithResult()
			diff := Diff(tt.wantError, errReq)
			assert.Empty(t, diff)

			if len(tt.wantError) > 0 {
				_, _ = csrData, csrName
				return
			}

			csrObj, err := operationWithResult()
			assert.NoError(t, err) // want "require-error: for error assertions use require"

			_, err = operationWithResult()
			require.NoError(t, err)

			certData, err := operationWithResult()
			assert.NoError(t, err) // want "require-error: for error assertions use require"

			certs, err := operationWithResult()
			assert.NoError(t, err) // want "require-error: for error assertions use require"

			switch l := len(certs.([]any)); l {
			case 1:
				// good
			default:
				t.Errorf("expected 1 cert, got %d", l)
				for i, certificate := range certs.([]any) {
					t.Log(i, certificate)
				}
				t.FailNow()
			}

			_, _ = csrObj, certData
			got := time.Second
			assert.Equal(t, tt.wantDuration, got)
		})
	}
}

type RoleRegistration struct {
	Role   RoleDTO
	Grants []string
}

type RoleDTO struct {
	Name string `json:"name"`
}

func TestService_DeclareFixedRoles(t *testing.T) {
	tests := []struct {
		name          string
		registrations []RoleRegistration
		wantErr       bool
		err           error
	}{
		{
			name:    "should work with empty list",
			wantErr: false,
		},
		{
			name: "should add registration",
			registrations: []RoleRegistration{
				{
					Role:   RoleDTO{Name: "fixed:test:test"},
					Grants: []string{"Admin"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := operation()
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.NoError(t, err) // want "require-error: for error assertions use require"

			registrationCnt := 0
			assert.Equal(t, len(tt.registrations), registrationCnt,
				"expected service registration list to contain all test registrations")
		})
	}
}

func createTempPackageJson(t *testing.T, version string) error {
	t.Helper()

	data := struct{}{}
	file, _ := json.MarshalIndent(data, "", " ")

	err := os.WriteFile("package.json", file, 0644)
	assert.NoError(t, err) // want "require-error: for error assertions use require"

	t.Cleanup(func() {
		err := os.RemoveAll("package.json")
		assert.NoError(t, err)
	})
	return nil
}

func operation() error                  { return nil }
func operationWithResult() (any, error) { return nil, nil }

type IdentifierMarker struct{}                                       //
func NewIdentifierMarker() *IdentifierMarker                         { return new(IdentifierMarker) } //
func (im *IdentifierMarker) Process(issues []Issue) ([]Issue, error) { return nil, nil }

type Issue struct {
	FromLinter string
	Text       string
}

func ValidateAppArmorProfileFormat(profile string) error {
	return nil
}

func GetPortForward(ctx context.Context, podName, podNamespace string, podUID string) (*url.URL, error) {
	return nil, nil
}

func getDoc(filePath string) (string, error) {
	return "", nil
}

func loadDriverDefinition(filename string) (string, error) {
	return "", nil
}

func TryLoadKeyFromDisk(pkiPath, name string) (crypto.Signer, error) {
	return nil, nil
}

func TryLoadCSRFromDisk(pkiPath, name string) (*x509.CertificateRequest, error) {
	return nil, nil
}

func runGenCSR(out io.Writer, config *int) error {
	return nil
}

func Decode(in string, to any) error {
	return nil
}

type Cursor struct{}

func Error(status int, message string, err error) *NormalResponse {
	return nil
}

type NormalResponse struct {
	status     int
	body       *bytes.Buffer
	header     http.Header
	errMessage string
	err        error
}

func (r *NormalResponse) Err() error {
	return r.err
}

func calcRestartCountByLogDir(filenames []string) (int, error) { return 0, nil }

type Status int                                 //
const StatusNotFound Status = 404               //
func NewStatusError(msg string, s Status) error { return nil } //
type StatusError struct{ error }                //
func (StatusError) StatusID() Status            { return 0 } //

var ChatIDNil ChatID             //
type ChatID struct{}             //
func NewChatID() ChatID          { return ChatID{} } //
func (c ChatID) Validate() error { return nil }

type TestServer interface {
	TearDownFn()
}

func StartTestServerOrDie(t *testing.T) TestServer {
	return nil
}

func Diff(x, y any) string {
	return ""
}

func MakeCSR(privateKey interface{}, subject *pkix.Name, dnsSANs []string, ipSANs []net.IP) (csr []byte, err error) {
	return nil, nil
}
