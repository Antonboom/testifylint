package checkers_test

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestFreshImportLocalName_TopLevelCollision(t *testing.T) {
	// A file with a top-level identifier named "http" must not receive
	// an import whose local qualifier is also "http", because that would
	// re-declare the name in the file block and cause a compile error.
	src := `package foo

import "testing"

var http = struct{}{}

func TestFoo(t *testing.T) { _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	name := checkers.FreshImportLocalName(f, "http", "net/http")
	if name != "http1" {
		t.Errorf("FreshImportLocalName = %q, want %q", name, "http1")
	}
}

func TestFreshImportLocalName_FuncCollision(t *testing.T) {
	// A file with a top-level function named "http" should also trigger
	// the fallback to "http1".
	src := `package foo

import "testing"

func http() {}

func TestFoo(t *testing.T) { _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	name := checkers.FreshImportLocalName(f, "http", "net/http")
	if name != "http1" {
		t.Errorf("FreshImportLocalName = %q, want %q", name, "http1")
	}
}

func TestFreshImportLocalName_NoCollision(t *testing.T) {
	// A file without any "http" top-level identifier should receive
	// the preferred import name "http".
	src := `package foo

import "testing"

func TestFoo(t *testing.T) { _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	name := checkers.FreshImportLocalName(f, "http", "net/http")
	if name != "http" {
		t.Errorf("FreshImportLocalName = %q, want %q", name, "http")
	}
}

func TestFileTopLevelNames(t *testing.T) {
	src := `package foo

import "testing"

const myConst = 1

var myVar = 2

type MyType struct{}

func myFunc() {}

func TestFoo(t *testing.T) { _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	names := checkers.FileTopLevelNames(f)

	for _, want := range []string{"myConst", "myVar", "MyType", "myFunc", "TestFoo"} {
		if _, ok := names[want]; !ok {
			t.Errorf("FileTopLevelNames missing %q", want)
		}
	}
	// Import names should NOT appear in top-level names.
	if _, ok := names["testing"]; ok {
		t.Errorf("FileTopLevelNames unexpectedly contains import name %q", "testing")
	}
}

// applyAndFormat calls AddImportFix, applies the resulting TextEdit to src,
// formats the result with go/format, and returns the local import name and
// the final formatted source. Both the result and the expected string are passed
// through go/format so that whitespace differences (tabs vs spaces) do not affect
// the comparison.
func applyAndFormat(t *testing.T, src, pkgPath string) (name, result string) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo_test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	n, edit, ok := checkers.AddImportFix([]*ast.File{f}, f.Pos()+1, pkgPath)
	if !ok {
		t.Fatalf("AddImportFix returned ok=false")
	}

	srcBytes := []byte(src)
	if edit != nil {
		offset := fset.Position(edit.Pos).Offset
		endOffset := fset.Position(edit.End).Offset
		srcBytes = append(srcBytes[:offset:offset], append(edit.NewText, srcBytes[endOffset:]...)...)
	}

	formatted, err := format.Source(srcBytes)
	if err != nil {
		t.Fatalf("format.Source error: %v\nresult before formatting:\n%s", err, srcBytes)
	}
	return n, string(formatted)
}

// TestAddImportFix_AlreadyImported verifies that when pkgPath is already present
// AddImportFix returns the correct local name and a nil edit.
func TestAddImportFix_AlreadyImported(t *testing.T) {
	src := `package foo

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) { _ = assert.Equal; _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo_test.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	name, edit, ok := checkers.AddImportFix([]*ast.File{f}, f.Pos()+1, "github.com/stretchr/testify/assert")
	if !ok {
		t.Fatal("AddImportFix returned ok=false for already-imported package")
	}
	if name != "assert" {
		t.Errorf("AddImportFix name = %q, want %q", name, "assert")
	}
	if edit != nil {
		t.Error("AddImportFix returned non-nil edit for already-imported package")
	}
}

// TestAddImportFix_GCIOrdering verifies that when a package is alphabetically
// before an existing non-stdlib import, AddImportFix inserts it before that import
// so that go/format keeps the imports sorted within their GCI group.
func TestAddImportFix_GCIOrdering(t *testing.T) {
	// File already imports "require"; adding "assert" (a < r) must place assert
	// before require after go/format sorts within the group.
	src := `package foo

import (
"testing"

"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) { _ = require.NoError; _ = t }
`
	// want uses the same structure; go/format normalises whitespace on both sides.
	want := `package foo

import (
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) { _ = require.NoError; _ = t }
`
	wantFmt, err := format.Source([]byte(want))
	if err != nil {
		t.Fatalf("format.Source on want: %v", err)
	}

	name, got := applyAndFormat(t, src, "github.com/stretchr/testify/assert")
	if name != "assert" {
		t.Errorf("AddImportFix name = %q, want %q", name, "assert")
	}
	if got != string(wantFmt) {
		t.Errorf("unexpected result:\nwant:\n%s\ngot:\n%s", wantFmt, got)
	}
}

// TestAddImportFix_GCIOrderingPreservesGroup verifies that when a file has GCI
// blank-line-separated import groups, inserting an external package before a
// lexicographically-greater non-stdlib spec keeps it in the correct group (i.e.
// it does not end up after local imports).
func TestAddImportFix_GCIOrderingPreservesGroup(t *testing.T) {
	// File has an external group ("require") and a local group ("localmod/pkg").
	// Adding "assert" (a < r) must place it before "require", staying in the
	// external group, not in the local group.
	src := `package foo

import (
"testing"

"github.com/stretchr/testify/require"

"github.com/localmod/pkg"
)

func TestFoo(t *testing.T) { _ = require.NoError; _ = pkg.Foo; _ = t }
`
	want := `package foo

import (
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

"github.com/localmod/pkg"
)

func TestFoo(t *testing.T) { _ = require.NoError; _ = pkg.Foo; _ = t }
`
	wantFmt, err := format.Source([]byte(want))
	if err != nil {
		t.Fatalf("format.Source on want: %v", err)
	}

	name, got := applyAndFormat(t, src, "github.com/stretchr/testify/assert")
	if name != "assert" {
		t.Errorf("AddImportFix name = %q, want %q", name, "assert")
	}
	if got != string(wantFmt) {
		t.Errorf("unexpected result:\nwant:\n%s\ngot:\n%s", wantFmt, got)
	}
}

// TestAddImportFix_BlankImported verifies that AddImportFix returns ok=false when
// the target package is blank-imported (import _ "pkg"), because the package is
// present but its symbols are inaccessible via a qualifier.
func TestAddImportFix_BlankImported(t *testing.T) {
	src := `package foo

import (
"testing"

_ "github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) { _ = t }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo_test.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, _, ok := checkers.AddImportFix([]*ast.File{f}, f.Pos()+1, "github.com/stretchr/testify/assert")
	if ok {
		t.Error("AddImportFix returned ok=true for blank-imported package; want ok=false")
	}
}
