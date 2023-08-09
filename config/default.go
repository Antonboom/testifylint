package config

import (
	"io"
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

// Default is default testifylint config.
var Default = Config{
	EnabledCheckers: checkers.EnabledByDefault(),
	ExpectedActual: ExpectedActualConfig{
		Pattern: Regexp{checkers.DefaultExpectedVarPattern},
	},
}

// DumpDefault dumps more complex YML than just marshalled Default,
// output contains comments and extra padding.
func DumpDefault(out io.Writer) error {
	return defaultConfTmpl.Execute(out, struct {
		Default  Config
		Checkers string
	}{
		Default:  Default,
		Checkers: buildCheckersYML(),
	})
}

var defaultConfTmpl = template.Must(template.New(".testifylint.yml").Parse(`enabled-checkers:
{{ .Checkers }}
expected-actual:
  # Pattern for expected variable name.
  pattern: {{ .Default.ExpectedActual.Pattern }}
`))

func buildCheckersYML() string {
	var result strings.Builder
	for _, checkerName := range checkers.All() {
		if checkers.IsEnabledByDefault(checkerName) {
			result.WriteString("  - " + checkerName + "\n")
		} else {
			result.WriteString("  # - " + checkerName + "\n")
		}
	}
	return result.String()
}
