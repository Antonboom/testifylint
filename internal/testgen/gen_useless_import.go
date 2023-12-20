package main

import (
	"strings"
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
		report  = checker + ": avoid the import of %s as _ since it doesn't do anything"
	)

	return struct {
		CheckerName    CheckerName
		ReportTemplate string
		Packages       []string
		GoldenPackages []string
	}{
		CheckerName:    CheckerName(checker),
		ReportTemplate: report,
		Packages: []string{
			"github.com/stretchr/testify/assert",
			"github.com/stretchr/testify/require",
			"github.com/stretchr/testify/suite",
			"github.com/stretchr/testify/mock",
			"github.com/stretchr/testify/http",
		},
		GoldenPackages: []string{},
	}
}

func (UselessImportTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("UselessImportTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(uselessImportTestTmpl))
}

func (UselessImportTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("UselessImportTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(strings.ReplaceAll(uselessImportTestTmpl, "Packages", "GoldenPackages")))
}

const uselessImportTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"
	{{range $.Packages}}
    _ "{{.}}" // want "{{printf $.ReportTemplate . }}"
	{{- end}}
)

func TestDummy(t *testing.T) {
	dummy := 1 + 3
	if dummy != 4 {
		t.Errorf("Dummy should return %d, but got %d", 4, dummy)
	}
}
`
