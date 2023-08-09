package analysisutil_test

import (
	"go/parser"
	"go/token"
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
			path:     "vendor/github.com/stretchr/testify/require",
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			isPkg := analysisutil.IsPkg(tt.pkg, tt.name, tt.path)
			if isPkg != tt.expected {
				t.Fatalf("unexpected result for case: %v", tt.pkg.String())
			}
		})
	}
}

func TestImports(t *testing.T) {
	fset := token.NewFileSet()

	src := `package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	assert.Equal(t, 4, 2*2)
}`

	f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("import", func(t *testing.T) {
		for _, imp := range []string{
			"testing",
			"github.com/stretchr/testify/assert",
		} {
			t.Run(imp, func(t *testing.T) {
				if !analysisutil.Imports(f, imp) {
					t.FailNow()
				}
			})
		}
	})

	t.Run("do not import", func(t *testing.T) {
		for _, imp := range []string{
			"",
			"net/http",
			"net/http/httptest",
			"github.com/stretchr/testify/suite",
			"github.com/stretchr/testify/require",
			"vendor/github.com/stretchr/testify/require",
		} {
			t.Run(imp, func(t *testing.T) {
				if analysisutil.Imports(f, imp) {
					t.FailNow()
				}
			})
		}
	})
}
