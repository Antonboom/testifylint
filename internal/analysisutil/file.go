package analysisutil

import (
	"go/ast"
	"go/token"
	"path"
	"slices"
	"strconv"
	"strings"
)

// PkgBaseName returns the default package name for the given import path.
// It handles versioned module paths (e.g. "example.com/pkg/v2" -> "pkg").
func PkgBaseName(importPath string) string {
	base := path.Base(importPath)
	after, found := strings.CutPrefix(base, "v")
	if !found {
		return base
	}
	v, err := strconv.Atoi(after)
	if err != nil || v < 2 {
		return base
	}
	return path.Base(path.Dir(importPath))
}

// Imports tells if the file imports at least one of the packages.
// If no packages provided then function returns false.
func Imports(file *ast.File, pkgs ...string) bool {
	for _, i := range file.Imports {
		if i.Path == nil {
			continue
		}

		importPath, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if slices.Contains(pkgs, importPath) { // Small O(n).
			return true
		}
	}
	return false
}

// LocalPkgName returns the local name used for an import with the given path in the
// source file that contains pos.
//
// It returns ("", true) for dot-imports, ("name", true) for regular or aliased
// imports and ("", false) when the import isn't found in the file or uses "_".
func LocalPkgName(files []*ast.File, pos token.Pos, pkgPath string) (string, bool) {
	for _, file := range files {
		if file.Pos() > pos || pos > file.End() {
			continue
		}
		for _, imp := range file.Imports {
			importPath, err := strconv.Unquote(imp.Path.Value)
			if err != nil || importPath != pkgPath {
				continue
			}
			if imp.Name != nil {
				if imp.Name.Name == "." {
					return "", true
				}
				if imp.Name.Name == "_" {
					return "", false
				}
				return imp.Name.Name, true
			}
			return PkgBaseName(pkgPath), true
		}
		break
	}
	return "", false
}
