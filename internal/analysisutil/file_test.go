package analysisutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

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

func TestPkgBaseName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		importPath string
		want       string
	}{
		{"net/http", "http"},
		{"fmt", "fmt"},
		{"github.com/stretchr/testify/assert", "assert"},
		{"example.com/pkg/v2", "pkg"},
		{"example.com/pkg/v10", "pkg"},
		{"example.com/v2", "example.com"},        // version at top level: falls back to domain
		{"example.com/mypkg/v2alpha", "v2alpha"}, // not a pure integer suffix
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			t.Parallel()

			got := analysisutil.PkgBaseName(tt.importPath)
			if got != tt.want {
				t.Errorf("PkgBaseName(%q) = %q, want %q", tt.importPath, got, tt.want)
			}
		})
	}
}

func TestLocalPkgName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		src      string
		pkgPath  string
		wantName string
		wantOK   bool
	}{
		{
			name: "regular import uses last path element",
			src: `package p
import "net/http"
var _ = http.Get`,
			pkgPath:  "net/http",
			wantName: "http",
			wantOK:   true,
		},
		{
			name: "aliased import returns alias",
			src: `package p
import stdhttp "net/http"
var _ = stdhttp.Get`,
			pkgPath:  "net/http",
			wantName: "stdhttp",
			wantOK:   true,
		},
		{
			name: "dot-import returns empty string with ok=true",
			src: `package p
import . "net/http"
var _ = Get`,
			pkgPath:  "net/http",
			wantName: "",
			wantOK:   true,
		},
		{
			name: "blank import returns empty string with ok=false",
			src: `package p
import _ "net/http"`,
			pkgPath:  "net/http",
			wantName: "",
			wantOK:   false,
		},
		{
			name: "package not imported returns empty string with ok=false",
			src: `package p
import "fmt"`,
			pkgPath:  "net/http",
			wantName: "",
			wantOK:   false,
		},
		{
			name: "versioned module path returns non-version element",
			src: `package p
import "example.com/pkg/v2"
`,
			pkgPath:  "example.com/pkg/v2",
			wantName: "pkg",
			wantOK:   true,
		},
		{
			name: "versioned module path with alias returns alias",
			src: `package p
import mypkg "example.com/pkg/v2"
`,
			pkgPath:  "example.com/pkg/v2",
			wantName: "mypkg",
			wantOK:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.src, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			// Use the file's own start position so LocalPkgName finds the file.
			pos := f.Pos()
			name, ok := analysisutil.LocalPkgName([]*ast.File{f}, pos, tt.pkgPath)
			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
		})
	}
}
