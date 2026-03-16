package checkers_test

import (
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
