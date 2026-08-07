package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// maxImportNameRetries is the number of numeric suffixes tried when finding a
// non-conflicting local name for an import (e.g., http1 … http9).
// Nine attempts is sufficient because import name collisions are extremely rare
// in practice, and any file with 10+ different imports sharing the same base
// name would be pathological.
const maxImportNameRetries = 9

// addImportFix returns the local qualifier name for pkgPath in the file containing pos
// and an optional TextEdit to add the import if it is not already present.
//
// Return values:
//   - ("",   nil,  true)  — pkgPath is dot-imported (symbols accessible without qualifier)
//   - (name, nil,  true)  — pkgPath is already imported under name
//   - (name, edit, true)  — pkgPath is absent and a non-conflicting import can be inserted
//   - ("",   nil,  false) — pkgPath is blank-imported or all candidate names are taken
func addImportFix(files []*ast.File, pos token.Pos, pkgPath string) (name string, edit *analysis.TextEdit, ok bool) {
	// Find the file containing pos.
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

	// Check whether pkgPath is already imported in this file.
	for _, imp := range file.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil || p != pkgPath {
			continue
		}
		if imp.Name != nil {
			switch imp.Name.Name {
			case ".":
				return "", nil, true // dot-import: symbols accessible without qualifier
			case "_":
				return "", nil, false // blank import: package not accessible
			default:
				return imp.Name.Name, nil, true // aliased import
			}
		}
		return importBaseName(pkgPath), nil, true // regular import
	}

	// Package is not imported — compute a non-conflicting local name and insert.
	preferred := importBaseName(pkgPath)
	localName := freshImportLocalName(file, preferred, pkgPath)
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

// importBaseName returns the default local name for the given import path.
// It handles versioned module paths (e.g., "example.com/pkg/v2" → "pkg").
func importBaseName(importPath string) string {
	base := path.Base(importPath)
	if len(base) > 1 && base[0] == 'v' {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			base = path.Base(path.Dir(importPath))
		}
	}
	return base
}

// isStdlibPath reports whether pkgPath belongs to the Go standard library.
// A standard library package has no dot in its first path segment.
func isStdlibPath(pkgPath string) bool {
	slash := strings.IndexByte(pkgPath, '/')
	if slash < 0 {
		slash = len(pkgPath)
	}
	return !strings.Contains(pkgPath[:slash], ".")
}

// freshImportLocalName returns a non-conflicting local name for an import of pkgPath.
// It starts with preferred and tries preferred1, preferred2, …, up to maxImportNameRetries.
// Returns "" if no suitable name is found.
//
// The check covers both existing import local names and package-level (file-scope)
// identifiers, because imported package names and top-level declarations share the
// same file namespace.
func freshImportLocalName(file *ast.File, preferred, pkgPath string) string {
	used := usedImportLocalNames(file, pkgPath)
	topLevel := fileTopLevelNames(file)

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
		if imp.Name != nil {
			if imp.Name.Name != "_" && imp.Name.Name != "." {
				names[imp.Name.Name] = struct{}{}
			}
		} else {
			names[importBaseName(p)] = struct{}{}
		}
	}
	return names
}

// fileTopLevelNames returns the set of names declared at file (package) scope:
// function names, variable names, constant names, and type names.
// These share the file-level namespace with imported package qualifiers.
func fileTopLevelNames(file *ast.File) map[string]struct{} {
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

		// Third-party package (or empty block): append before the closing ')'.
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
