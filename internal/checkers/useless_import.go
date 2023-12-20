package checkers

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	nullContainer                 = "_"
	uselessImportCallReportFormat = "avoid the import of %s as _ since it doesn't do anything"
)

var importPaths = map[string]bool{
	"\"github.com/stretchr/testify/assert\"":  true,
	"\"github.com/stretchr/testify/require\"": true,
	"\"github.com/stretchr/testify/suite\"":   true,
	"\"github.com/stretchr/testify/mock\"":    true,
	"\"github.com/stretchr/testify/http\"":    true,
}

// UselessImport detects useless imports of testify as _.
// These imports are useless since testify doesn't do any magic with init() function.
// It detects situation like this:
//
//	import (
//		"testing"
//
//		_ "github.com/stretchr/testify/assert"
//		_ "github.com/stretchr/testify/require"
//		_ "github.com/stretchr/testify/suite"
//		_ "github.com/stretchr/testify/mock"
//		_ "github.com/stretchr/testify/http"
//	)
type UselessImport struct{}

// NewUselessImport constructs UselessImport checker.
func NewUselessImport() UselessImport { return UselessImport{} }
func (UselessImport) Name() string    { return "useless-import" }

func (checker UselessImport) Check(pass *analysis.Pass, _ *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	for _, file := range pass.Files {
		for _, imp := range file.Imports {
			if imp.Name != nil && imp.Name.Name == nullContainer && importPaths[imp.Path.Value] {
				msg := fmt.Sprintf(uselessImportCallReportFormat, imp.Path.Value)
				fix := &analysis.SuggestedFix{
					Message: fmt.Sprintf("Remove import of %s as _", imp.Path.Value),
					TextEdits: []analysis.TextEdit{
						{
							Pos:     imp.Pos(),
							End:     imp.End(),
							NewText: []byte(""),
						},
					},
				}
				d := newDiagnostic(checker.Name(), imp, msg, fix)
				diagnostics = append(diagnostics, *d)
			}
		}
	}
	return diagnostics
}
