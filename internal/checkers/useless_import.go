package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

var packagesNotIntendedForBlankImport = map[string]struct{}{
	"github.com/stretchr/testify":         {},
	"github.com/stretchr/testify/assert":  {},
	"github.com/stretchr/testify/http":    {},
	"github.com/stretchr/testify/mock":    {},
	"github.com/stretchr/testify/require": {},
	"github.com/stretchr/testify/suite":   {},
}

// UselessImport detects useless blank imports of testify packages.
// These imports are useless since testify doesn't do any magic with init() function.
//
// The checker detects situations like
//
//	import (
//		"testing"
//
//		_ "github.com/stretchr/testify"
//		_ "github.com/stretchr/testify/assert"
//		_ "github.com/stretchr/testify/http"
//		_ "github.com/stretchr/testify/mock"
//		_ "github.com/stretchr/testify/require"
//		_ "github.com/stretchr/testify/suite"
//	)
//
// and requires
//
//	import (
//		"testing"
//	)
type UselessImport struct{}

// NewUselessImport constructs UselessImport checker.
func NewUselessImport() UselessImport { return UselessImport{} }
func (UselessImport) Name() string    { return "useless-import" }

func (checker UselessImport) Check(pass *analysis.Pass, _ *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	for _, file := range pass.Files {
		if len(file.Imports) == 0 {
			continue
		}

		for _, decl := range file.Decls {
			impDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if impDecl.Tok != token.IMPORT {
				continue
			}

			for _, spec := range impDecl.Specs {
				imp := spec.(*ast.ImportSpec)
				if imp.Name == nil || imp.Name.Name != "_" {
					continue
				}

				pkg, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					continue
				}
				if _, ok := packagesNotIntendedForBlankImport[pkg]; !ok {
					continue
				}

				msg := fmt.Sprintf("avoid blank import of %s as it does nothing", pkg)

				impStart, impEnd := getImportRange(impDecl, imp)
				fix := &analysis.SuggestedFix{
					Message: "Remove blank import of " + pkg,
					TextEdits: []analysis.TextEdit{
						{
							Pos:     impStart - 1,
							End:     impEnd,
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

func getImportRange(impDecl *ast.GenDecl, impSpec *ast.ImportSpec) (token.Pos, token.Pos) {
	start, end := impSpec.Pos(), impSpec.End()

	if len(impDecl.Specs) == 1 {
		start = impDecl.Pos()
	}
	if impSpec.Comment != nil {
		end = impSpec.Comment.End()
	}

	return start, end
}
