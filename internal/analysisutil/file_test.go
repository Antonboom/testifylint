package analysisutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestPkgBaseName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		importPath string
		want       string
	}{
		{name: "plain", importPath: "github.com/stretchr/testify/require", want: "require"},
		{name: "module v2 suffix", importPath: "example.com/pkg/v2", want: "pkg"},
		{name: "module v1 path segment", importPath: "example.com/v1", want: "v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := analysisutil.PkgBaseName(tt.importPath)
			if got != tt.want {
				t.Fatalf("PkgBaseName(%q) = %q, want %q", tt.importPath, got, tt.want)
			}
		})
	}
}

func TestImports(t *testing.T) {
	fset := token.NewFileSet()

	const src = `package simple

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

	notImported := []string{ //nolint:prealloc // This is simple test.
		"",
		"net/http",
		"net/http/httptest",
		"github.com/stretchr/testify/suite",
		"github.com/stretchr/testify/require",
		"vendor/github.com/stretchr/testify/require",
	}
	if analysisutil.Imports(f, notImported...) {
		t.FailNow()
	}
	if !analysisutil.Imports(f, append(notImported, "testing")...) {
		t.FailNow()
	}
	if !analysisutil.Imports(f, "github.com/stretchr/testify/assert") {
		t.FailNow()
	}
}

func TestLocalPkgName(t *testing.T) {
	t.Parallel()

	const src = `package localpkgname

import (
	"github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"
	. "net/http"
	_ "math"
)

func TestLocalPkgName() {
	assert.Equal(nil, 1, 1)
	r.Len(nil, []int{1}, 1)
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}

	pos := token.Pos(1)
	assertName, ok := analysisutil.LocalPkgName([]*ast.File{file}, pos, "github.com/stretchr/testify/assert")
	if !ok || assertName != "assert" {
		t.Fatalf("LocalPkgName(assert) = (%q, %t), want (assert, true)", assertName, ok)
	}

	requireName, ok := analysisutil.LocalPkgName([]*ast.File{file}, pos, "github.com/stretchr/testify/require")
	if !ok || requireName != "r" {
		t.Fatalf("LocalPkgName(require) = (%q, %t), want (r, true)", requireName, ok)
	}

	dotName, ok := analysisutil.LocalPkgName([]*ast.File{file}, pos, "net/http")
	if !ok || dotName != "" {
		t.Fatalf("LocalPkgName(dot) = (%q, %t), want (\"\", true)", dotName, ok)
	}

	blankName, ok := analysisutil.LocalPkgName([]*ast.File{file}, pos, "math")
	if ok || blankName != "" {
		t.Fatalf("LocalPkgName(blank) = (%q, %t), want (\"\", false)", blankName, ok)
	}
}
