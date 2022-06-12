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
		t.Fatal(err.Error())
	}

	svcTestFile, err := parser.ParseFile(fset, "service_test.go", `package servicetest`, parser.PackageClauseOnly)
	if err != nil {
		t.Fatal(err.Error())
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
