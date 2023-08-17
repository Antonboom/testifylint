package analysisutil_test

import (
	"go/types"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestIsPkg(t *testing.T) {
	cases := []struct {
		pkg        *types.Package
		name, path string
		expected   bool
	}{
		// Positive.
		{
			pkg:      types.NewPackage("testing", "testing"),
			name:     "testing",
			path:     "testing",
			expected: true,
		},
		{
			pkg:      types.NewPackage("net/http", "http"),
			name:     "http",
			path:     "net/http",
			expected: true,
		},
		{
			pkg:      types.NewPackage("github.com/stretchr/testify/assert", "assert"),
			name:     "assert",
			path:     "github.com/stretchr/testify/assert",
			expected: true,
		},
		{
			pkg:      types.NewPackage("vendor/github.com/stretchr/testify/require", "require"),
			name:     "require",
			path:     "github.com/stretchr/testify/require",
			expected: true,
		},

		// Negative.
		{
			pkg:      types.NewPackage("net/http", "http"),
			name:     "http",
			path:     "http",
			expected: false,
		},
		{
			pkg:      types.NewPackage("net/http", "http"),
			name:     "httptest",
			path:     "net/http",
			expected: false,
		},
		{
			pkg:      types.NewPackage("net/http", "http"),
			name:     "httptest",
			path:     "net/http/httptest",
			expected: false,
		},
		{
			pkg:      types.NewPackage("vendor/github.com/stretchr/testify/require", "require"),
			name:     "require",
			path:     "vendor/github.com/stretchr/testify/require", // Invalid import path.
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			isPkg := analysisutil.IsPkg(tt.pkg, tt.name, tt.path)
			if isPkg != tt.expected {
				t.FailNow()
			}
		})
	}
}
