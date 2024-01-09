package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
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
			"github.com/stretchr/testify",
			"github.com/stretchr/testify/assert",
			"github.com/stretchr/testify/http",
			"github.com/stretchr/testify/mock",
			"github.com/stretchr/testify/require",
			"github.com/stretchr/testify/suite",
		},
	}
}

func (UselessImportTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("UselessImportTestsGenerator.ErroredTemplate").
		Parse(uselessImportTestTmpl))
}

func (UselessImportTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("UselessImportTestsGenerator.GoldenTemplate").
		Parse(uselessImportTestTmplGolden))
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
	{{ range $.Packages }}
	_ "{{.}}" // want "{{printf $.ReportFmt . }}"
	{{- end}}
	
	_ "github.com/pmezard/go-difflib/difflib"
	. "database/sql"
)

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

const uselessImportTestTmplGolden = header + `
package {{ .CheckerName.AsPkgName }}

import "testing"
import "strings"

import (
)

import (
	"net/url"
	_ "gopkg.in/yaml.v3"
	// Testing.
	
	_ "github.com/pmezard/go-difflib/difflib"
	. "database/sql"
)

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
