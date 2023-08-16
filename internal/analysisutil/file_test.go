package analysisutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestIsTestFile(t *testing.T) {
	fset := token.NewFileSet()

	svcFile, err := parser.ParseFile(fset, "service.go", `package service`, parser.PackageClauseOnly)
	if err != nil {
		t.Fatal(err)
	}

	svcTestFile, err := parser.ParseFile(fset, "service_test.go", `package servicetest`, parser.PackageClauseOnly)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{svcFile, svcTestFile},
	}

	if v := analysisutil.IsTestFile(pass, svcFile); v {
		t.Errorf("expected IsTestFile == false")
	}
	if v := analysisutil.IsTestFile(pass, svcTestFile); !v {
		t.Errorf("expected IsTestFile == true")
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
