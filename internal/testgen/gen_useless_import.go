package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
	"github.com/Antonboom/testifylint/internal/testify"
)

type UselessImportTestsGenerator struct{}

func (UselessImportTestsGenerator) Checker() checkers.Checker {
	return checkers.NewUselessImport()
}

func (g UselessImportTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + ": avoid blank import of %s as it does nothing"
	)

	return struct {
		CheckerName CheckerName
		ReportFmt   string
		Packages    []string
	}{
		CheckerName: CheckerName(checker),
		ReportFmt:   report,
		Packages: []string{
			testify.ModulePath,
			testify.AssertPkgPath,
			testify.HTTPPkgPath,
			testify.MockPkgPath,
			testify.RequirePkgPath,
			testify.SuitePkgPath,
		},
	}
}

func (UselessImportTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("UselessImportTestsGenerator.ErroredTemplate").
		Parse(uselessImportTestTmpl))
}

func (UselessImportTestsGenerator) GoldenTemplate() Executor {
	// NOTE(a.telyshev): Auto-fixing introduces complexity (a lot of import combinations)
	// into such a simple and rarely used check.
	return nil
}

const uselessImportTestTmpl = header + `
package {{ .CheckerName.AsPkgName }}

import "testing"
{{- range $.Packages }}
import _ "{{.}}" // want "{{printf $.ReportFmt . }}"
{{- end }}
import "strings"

import (
	{{- range $.Packages }}
	_ "{{.}}" // want "{{printf $.ReportFmt . }}"
	{{- end }}
)

import (
	"net/url"
	_ "gopkg.in/yaml.v3"

	// Testing.
	{{- range $.Packages }}
	_ "{{.}}" // want "{{printf $.ReportFmt . }}"
	{{- end}}

	_ "github.com/pmezard/go-difflib/difflib"
	. "database/sql"
)

{{ with $pkg := (index $.Packages 0) -}}
import (
	_ "github.com/stretchr/testify" // want "{{printf $.ReportFmt $pkg }}"
)

import (
	// Test dependencies so that it doesn't get cleaned by glide vc
	_ "github.com/stretchr/testify" // want "{{printf $.ReportFmt $pkg }}"
)
{{- end }}

func TestDummy(t *testing.T) {
	dummy := 1 + 3
	if dummy != 4 {
		t.Errorf("expected %d, but got %d", 4, dummy)
	}

	_ = strings.Builder{}
	_ = url.URL{}
	_ = DB{}
}
`
