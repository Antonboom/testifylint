package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// maxImportNameRetries is the number of numeric suffixes tried when finding a
// non-conflicting local name for an import (e.g., http1 … http9).
// Nine attempts is sufficient because import name collisions are extremely rare
// in practice, and any file with 10+ different imports sharing the same base
// name would be pathological.
const maxImportNameRetries = 9

// AddImportFix returns the local qualifier name for pkgPath in the file containing pos
// and an optional TextEdit to add the import if it is not already present.
//
// Return values:
//   - ("",   nil,  true)  — pkgPath is dot-imported (symbols accessible without qualifier)
//   - (name, nil,  true)  — pkgPath is already imported under name
//   - (name, edit, true)  — pkgPath is absent and a non-conflicting import can be inserted
//   - ("",   nil,  false) — pkgPath is blank-imported or all candidate names are taken
func AddImportFix(files []*ast.File, pos token.Pos, pkgPath string) (name string, edit *analysis.TextEdit, ok bool) {
	// Use LocalPkgName to check whether pkgPath is already imported in this file.
	localName, imported := analysisutil.LocalPkgName(files, pos, pkgPath)
	if imported {
		return localName, nil, true
	}

	// LocalPkgName returns ("", false) for both a blank import and a missing import.
	// Find the file containing pos so we can distinguish the two cases and, if the
	// package is absent, build the TextEdit to insert it.
	var file *ast.File
	for _, f := range files {
		if f.Pos() <= pos && pos <= f.End() {
			file = f
			break
		}
	}
	if file == nil {
		return "", nil, false
	}

	// If the package appears with a blank name it is not accessible — bail out.
	for _, imp := range file.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil || p != pkgPath {
			continue
		}
		// LocalPkgName already returned true for dot/regular/aliased imports;
		// only a blank import can reach this point.
		if imp.Name != nil && imp.Name.Name == "_" {
			return "", nil, false // blank import: package not accessible
		}
	}

	// Package is not imported — compute a non-conflicting local name and insert.
	preferred := analysisutil.PkgBaseName(pkgPath)
	localName = FreshImportLocalName(file, preferred, pkgPath)
	if localName == "" {
		return "", nil, false
	}

	specText := strconv.Quote(pkgPath)
	if localName != preferred {
		specText = localName + " " + specText
	}
	textEdit := importInsertEdit(file, pkgPath, specText)
	return localName, &textEdit, true
}

// isStdlibPath reports whether pkgPath belongs to the Go standard library.
// A standard library package has no dot in its first path segment.
func isStdlibPath(pkgPath string) bool {
	// Split into at most 2 parts at the first slash
	firstSegment, _, _ := strings.Cut(pkgPath, "/")
	return !strings.Contains(firstSegment, ".")
}

// FreshImportLocalName returns a non-conflicting local name for an import of pkgPath.
// It starts with preferred and tries preferred1, preferred2, …, up to maxImportNameRetries.
// Returns "" if no suitable name is found.
//
// The check covers both existing import local names and package-level (file-scope)
// identifiers, because imported package names and top-level declarations share the
// same file namespace.
func FreshImportLocalName(file *ast.File, preferred, pkgPath string) string {
	used := usedImportLocalNames(file, pkgPath)
	topLevel := FileTopLevelNames(file)

	isAvailable := func(name string) bool {
		_, importTaken := used[name]
		_, topLevelTaken := topLevel[name]
		return !importTaken && !topLevelTaken
	}

	if isAvailable(preferred) {
		return preferred
	}
	for i := 1; i <= maxImportNameRetries; i++ {
		candidate := fmt.Sprintf("%s%d", preferred, i)
		if isAvailable(candidate) {
			return candidate
		}
	}
	return ""
}

// usedImportLocalNames returns the set of local names already used by imports in
// the file, excluding any import of excludePath (the package we want to add).
func usedImportLocalNames(file *ast.File, excludePath string) map[string]struct{} {
	names := make(map[string]struct{})
	for _, imp := range file.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil || p == excludePath {
			continue
		}
		if imp.Name == nil {
			names[analysisutil.PkgBaseName(p)] = struct{}{}
		} else if imp.Name.Name != "_" && imp.Name.Name != "." {
			names[imp.Name.Name] = struct{}{}
		}
	}
	return names
}

// FileTopLevelNames returns the set of names declared at file (package) scope:
// function names, variable names, constant names, and type names.
// These share the file-level namespace with imported package qualifiers.
func FileTopLevelNames(file *ast.File) map[string]struct{} {
	names := make(map[string]struct{})
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name != nil {
				names[d.Name.Name] = struct{}{}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec: // var and const
					for _, name := range s.Names {
						names[name.Name] = struct{}{}
					}
				case *ast.TypeSpec:
					names[s.Name.Name] = struct{}{}
				}
			}
		}
	}
	return names
}

// importInsertEdit returns a TextEdit that inserts specText into the file's import
// block at the correct position for pkgPath (stdlib imports are placed after the last
// existing stdlib spec to stay in the stdlib group; third-party imports are appended
// before the closing ')'). If no suitable parenthesized block exists, a new standalone
// import statement is created before the first non-import declaration.
func importInsertEdit(file *ast.File, pkgPath, specText string) analysis.TextEdit {
	stdlib := isStdlibPath(pkgPath)

	// Prefer inserting into an existing parenthesized import block (not "C").
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT || !genDecl.Lparen.IsValid() {
			continue
		}
		// Skip cgo import blocks to avoid disrupting the association with doc comments.
		hasCgo := false
		for _, spec := range genDecl.Specs {
			imp := spec.(*ast.ImportSpec)
			if p, _ := strconv.Unquote(imp.Path.Value); p == "C" {
				hasCgo = true
				break
			}
		}
		if hasCgo {
			continue
		}

		if stdlib {
			// Insert after the last stdlib spec so that stdlib imports stay grouped.
			var lastStdlib ast.Spec
			for _, spec := range genDecl.Specs {
				imp := spec.(*ast.ImportSpec)
				p, err := strconv.Unquote(imp.Path.Value)
				if err == nil && isStdlibPath(p) {
					lastStdlib = spec
				}
			}
			if lastStdlib != nil {
				// Insert on a new line immediately after the last stdlib spec.
				return analysis.TextEdit{
					Pos:     lastStdlib.End(),
					End:     lastStdlib.End(),
					NewText: []byte("\n\t" + specText),
				}
			}
			// No existing stdlib spec — insert at the top of the import block.
			if len(genDecl.Specs) > 0 {
				return analysis.TextEdit{
					Pos:     genDecl.Specs[0].Pos(),
					End:     genDecl.Specs[0].Pos(),
					NewText: []byte(specText + "\n\t"),
				}
			}
		}

		// Non-stdlib (third-party or local): insert at the alphabetical position among
		// non-stdlib specs to respect GCI import-group ordering. Scanning specs in file
		// order and inserting before the first non-stdlib entry whose path is
		// lexicographically greater preserves any blank-line group separators that GCI
		// may have added, because blank lines live between spec positions and are
		// unaffected by inserting text at a spec boundary.
		for _, spec := range genDecl.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			p, err := strconv.Unquote(imp.Path.Value)
			if err != nil || isStdlibPath(p) {
				continue
			}
			if p > pkgPath {
				// Insert before this spec; the file already has the tab indent before it,
				// so we append "\n\t" after specText to keep the displaced spec indented.
				return analysis.TextEdit{
					Pos:     spec.Pos(),
					End:     spec.Pos(),
					NewText: []byte(specText + "\n\t"),
				}
			}
		}

		// All non-stdlib specs are alphabetically ≤ pkgPath, or there are none:
		// append before the closing ')'.
		return analysis.TextEdit{
			Pos:     genDecl.Rparen,
			End:     genDecl.Rparen,
			NewText: []byte("\t" + specText + "\n"),
		}
	}

	// No suitable parenthesized block — insert a standalone import statement
	// before the first non-import declaration.
	var insertPos token.Pos
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			continue
		}
		insertPos = decl.Pos()
		break
	}
	if !insertPos.IsValid() {
		insertPos = file.End()
	}
	return analysis.TextEdit{
		Pos:     insertPos,
		End:     insertPos,
		NewText: []byte("import " + specText + "\n\n"),
	}
}
