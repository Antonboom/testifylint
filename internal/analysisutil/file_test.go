package analysisutil_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestIsTestFile(t *testing.T) {
	cases := []struct {
		filename   string
		fileSrc    string
		isTestFile bool
	}{
		{
			filename:   "service.go",
			fileSrc:    "package service",
			isTestFile: false,
		},
		{
			filename:   "service_unix.go",
			fileSrc:    "package service",
			isTestFile: false,
		},
		{
			filename:   "service_test.go",
			fileSrc:    "package service",
			isTestFile: true,
		},
		{
			filename:   "service_test.go",
			fileSrc:    "package service_test",
			isTestFile: true,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			fset := token.NewFileSet()

			file, err := parser.ParseFile(fset, tt.filename, tt.fileSrc, parser.PackageClauseOnly)
			if err != nil {
				t.Fatal(err)
			}

			if tt.isTestFile != analysisutil.IsTestFile(fset, file) {
				t.FailNow()
			}
		})
	}
}

func TestIsTestFile_FileIsNotInFileSet(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "service_test.go", "package service_test", parser.PackageClauseOnly)
	if err != nil {
		t.Fatal(err)
	}

	if analysisutil.IsTestFile(token.NewFileSet(), file) {
		t.FailNow()
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

	notImported := []string{
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
