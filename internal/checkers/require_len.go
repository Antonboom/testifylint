package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

const (
	requireLenForIndexReport = "for length assertions guarding index access use require"
	requireLenGuardReport    = "for indexed access use require.Len guard"
	defaultIndent            = "\t"
)

// RequireLen checks fail-fast guards for indexed access.
type RequireLen struct{}

// NewRequireLen constructs RequireLen checker.
func NewRequireLen() RequireLen { return RequireLen{} }
func (RequireLen) Name() string { return "require-len" }

func (checker RequireLen) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	callsByFunc := make(map[funcID][]*callMeta)

	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}
		if len(stack) < 3 {
			return true
		}

		fID := findSurroundingFunc(pass, stack)
		if fID == nil {
			return true
		}

		_, prevIsIfStmt := stack[len(stack)-2].(*ast.IfStmt)
		_, prevIsAssignStmt := stack[len(stack)-2].(*ast.AssignStmt)
		_, prevPrevIsIfStmt := stack[len(stack)-3].(*ast.IfStmt)
		inIfCond := prevIsIfStmt || (prevPrevIsIfStmt && prevIsAssignStmt)

		_, inBoolExpr := stack[len(stack)-2].(*ast.BinaryExpr)

		callExpr := node.(*ast.CallExpr)
		testifyCall := NewCallMeta(pass, callExpr)

		call := &callMeta{
			call:         callExpr,
			testifyCall:  testifyCall,
			rootIf:       findRootIf(stack),
			parentIf:     findNearestNode[*ast.IfStmt](stack),
			parentBlock:  findNearestNode[*ast.BlockStmt](stack),
			inIfCond:     inIfCond,
			inBoolExpr:   inBoolExpr,
			inNoErrorSeq: false,
		}

		callsByFunc[*fID] = append(callsByFunc[*fID], call)
		return testifyCall == nil
	})

	var diagnostics []analysis.Diagnostic

	callsByBlock := map[*ast.BlockStmt][]*callMeta{}
	fileContentCache := make(map[string][]byte)
	for _, calls := range callsByFunc {
		for _, c := range calls {
			if b := c.parentBlock; b != nil {
				callsByBlock[b] = append(callsByBlock[b], c)
			}
		}
	}

	markCallsInNoErrorSequence(callsByBlock)

	for funcInfo, calls := range callsByFunc {
		for i, c := range calls {
			if m := funcInfo.meta; m.isTestCleanup || m.isGoroutine || m.isHTTPHandler {
				continue
			}
			if c.testifyCall == nil || !c.testifyCall.IsAssert {
				continue
			}

			switch c.testifyCall.Fn.NameFTrimmed {
			case "Len", "Lenf":
				if needToSkipBasedOnContext(c, i, calls, callsByBlock) {
					continue
				}
				if shouldRequireLenForIndexedAccess(pass, c, i, calls) {
					diagnostics = append(diagnostics,
						*newDiagnostic(checker.Name(), c.testifyCall, requireLenForIndexReport))
				}
			default:
				if needToSkipForLenGuardContext(c, calls) {
					continue
				}
				if d := newRequireLenGuardDiagnostic(pass, checker.Name(), c, i, calls, fileContentCache); d != nil {
					diagnostics = append(diagnostics, *d)
				}
			}
		}
	}

	return diagnostics
}

func shouldRequireLenForIndexedAccess(
	pass *analysis.Pass,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
) bool {
	if len(currCall.testifyCall.Args) < 2 {
		return false
	}

	collectionExpr := currCall.testifyCall.Args[0]
	requiredLen, ok := isIntBasicLit(currCall.testifyCall.Args[1])
	if !ok || requiredLen <= 0 {
		return false
	}

	collectionStr := analysisutil.NodeString(pass.Fset, collectionExpr)
	if collectionStr == "" {
		return false
	}

	for i := currCallIndex + 1; i < len(otherCalls); i++ {
		if otherCalls[i].parentBlock != currCall.parentBlock {
			continue
		}
		if containsIndexedAccess(pass, otherCalls[i].call, collectionStr, requiredLen) {
			return true
		}
	}
	return false
}

func containsIndexedAccess(pass *analysis.Pass, node ast.Node, collection string, requiredLen int) bool {
	var found bool

	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}

		ie, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		idx, ok := isIntBasicLit(ie.Index)
		if !ok || idx < 0 || idx >= requiredLen {
			return true
		}

		if analysisutil.NodeString(pass.Fset, ie.X) == collection {
			found = true
			return false
		}
		return true
	})

	return found
}

func newRequireLenGuardDiagnostic(
	pass *analysis.Pass,
	checkerName string,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	fileContentCache map[string][]byte,
) *analysis.Diagnostic {
	var (
		exprStmt                  *ast.ExprStmt
		ok                        bool
		requireQualifier          string
		requireImportIsAccessible bool
	)
	if currCall.testifyCall.IsPkg && (len(currCall.testifyCall.ArgsRaw) >= 2) {
		tArg := currCall.testifyCall.ArgsRaw[0]
		if implementsTestingT(pass, tArg) {
			exprStmt, ok = findExprStmtForCall(currCall)
			requireQualifier, requireImportIsAccessible = analysisutil.LocalPkgName(
				pass.Files, currCall.call.Pos(), testify.RequirePkgPath)
		}
	}

	for _, target := range indexedAccesses(pass, currCall.call) {
		requiredLen := target.maxIndex + 1
		if hasLenGuard(pass, currCall, currCallIndex, otherCalls, target.collection, requiredLen) {
			continue
		}

		diagnostic := newDiagnostic(checkerName, currCall.testifyCall, requireLenGuardReport)
		if !ok {
			return diagnostic
		}

		additionalEdits := make([]analysis.TextEdit, 0, 2)
		if !requireImportIsAccessible || requireQualifier == "" {
			requireQualifier = availableRequireQualifier(pass.Files, currCall.call.Pos())

			importEdit, hasImportEdit := addRequireImportTextEdit(pass, currCall.call.Pos(), requireQualifier)
			if !hasImportEdit {
				return diagnostic
			}
			additionalEdits = append(additionalEdits, importEdit)
		}

		tArg := currCall.testifyCall.ArgsRaw[0]
		indent := lineIndentAtPos(pass, currCall.call.Pos(), fileContentCache)
		insertText := fmt.Sprintf("%s.Len(%s, %s, %d)\n%s",
			requireQualifier, analysisutil.NodeString(pass.Fset, tArg), target.collection, requiredLen, indent)
		fixMsg := "Insert `require.Len` guard"
		if requiredLen == 1 {
			insertText = fmt.Sprintf("%s.NotEmpty(%s, %s)\n%s",
				requireQualifier, analysisutil.NodeString(pass.Fset, tArg), target.collection, indent)
			fixMsg = "Insert `require.NotEmpty` guard"
		}
		return newDiagnostic(checkerName, currCall.testifyCall, requireLenGuardReport, analysis.SuggestedFix{
			Message: fixMsg,
			TextEdits: append(additionalEdits,
				analysis.TextEdit{
					Pos:     exprStmt.Pos(),
					End:     exprStmt.Pos(),
					NewText: []byte(insertText),
				},
			),
		})
	}
	return nil
}

func availableRequireQualifier(files []*ast.File, pos token.Pos) string {
	file := fileForPos(files, pos)
	if file == nil {
		return testify.RequirePkgName
	}

	for i := 0; ; i++ {
		qualifier := testify.RequirePkgName
		if i > 0 {
			qualifier = fmt.Sprintf("%s%d", testify.RequirePkgName, i)
		}
		if file.Scope == nil || file.Scope.Lookup(qualifier) == nil {
			return qualifier
		}
	}
}

func addRequireImportTextEdit(pass *analysis.Pass, pos token.Pos, requireQualifier string) (analysis.TextEdit, bool) {
	importSpec := fmt.Sprintf("%q", testify.RequirePkgPath)
	if requireQualifier != testify.RequirePkgName {
		importSpec = fmt.Sprintf("%s %q", requireQualifier, testify.RequirePkgPath)
	}

	for _, file := range pass.Files {
		if file.Pos() > pos || pos > file.End() {
			continue
		}

		if analysisutil.Imports(file, testify.RequirePkgPath) {
			return analysis.TextEdit{}, false
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.IMPORT {
				continue
			}

			if genDecl.Lparen.IsValid() {
				return analysis.TextEdit{
					Pos:     genDecl.Rparen,
					End:     genDecl.Rparen,
					NewText: []byte("\t" + importSpec + "\n"),
				}, true
			}

			if len(genDecl.Specs) == 0 {
				break
			}

			existingImport := analysisutil.NodeString(pass.Fset, genDecl.Specs[0])
			newImportBlock := fmt.Sprintf("import (\n\t%s\n\t%s\n)", existingImport, importSpec)
			return analysis.TextEdit{
				Pos:     genDecl.Pos(),
				End:     genDecl.End(),
				NewText: []byte(newImportBlock),
			}, true
		}

		return analysis.TextEdit{
			Pos:     file.Name.End(),
			End:     file.Name.End(),
			NewText: []byte(fmt.Sprintf("\n\nimport %s\n", importSpec)),
		}, true
	}

	return analysis.TextEdit{}, false
}

func fileForPos(files []*ast.File, pos token.Pos) *ast.File {
	for _, file := range files {
		if file.Pos() <= pos && pos <= file.End() {
			return file
		}
	}
	return nil
}

type indexedAccess struct {
	collection string
	maxIndex   int
}

func indexedAccesses(pass *analysis.Pass, node ast.Node) []indexedAccess {
	var result []indexedAccess
	indexByCollection := make(map[string]int)

	ast.Inspect(node, func(n ast.Node) bool {
		ie, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		idx, ok := isIntBasicLit(ie.Index)
		if !ok || idx < 0 {
			return true
		}

		collection := analysisutil.NodeString(pass.Fset, ie.X)
		if collection == "" {
			return true
		}

		i, exists := indexByCollection[collection]
		if !exists {
			indexByCollection[collection] = len(result)
			result = append(result, indexedAccess{collection: collection, maxIndex: idx})
			return true
		}
		if idx > result[i].maxIndex {
			result[i].maxIndex = idx
		}
		return true
	})

	return result
}

func hasLenGuard(
	pass *analysis.Pass,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	collection string,
	requiredLen int,
) bool {
	for i := range currCallIndex {
		c := otherCalls[i]
		if c.parentBlock != currCall.parentBlock || c.testifyCall == nil {
			continue
		}
		switch c.testifyCall.Fn.NameFTrimmed {
		case "Len", "Lenf":
			if c.testifyCall.IsAssert {
				continue
			}
			if len(c.testifyCall.Args) < 2 {
				continue
			}
			lenCollection := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[0])
			if lenCollection != collection {
				continue
			}
			if assertedLen, isBasicLit := isIntBasicLit(c.testifyCall.Args[1]); isBasicLit && (assertedLen < requiredLen) {
				continue
			}
			return true
		case "NotEmpty":
			if c.testifyCall.IsAssert {
				continue
			}
			if requiredLen != 1 || len(c.testifyCall.Args) < 1 {
				continue
			}
			notEmptyCollection := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[0])
			if notEmptyCollection != collection {
				continue
			}
			return true
		}
	}
	return false
}

func needToSkipForLenGuardContext(currCall *callMeta, otherCalls []*callMeta) bool {
	if currCall.inIfCond || currCall.inBoolExpr || currCall.inNoErrorSeq {
		return true
	}

	if currCall.rootIf != nil {
		for _, rootCall := range otherCalls {
			if (rootCall.rootIf == currCall.rootIf) && rootCall.inIfCond {
				return true
			}
		}
	}

	return false
}

func findExprStmtForCall(currCall *callMeta) (*ast.ExprStmt, bool) {
	if currCall.parentBlock == nil {
		return nil, false
	}

	for _, stmt := range currCall.parentBlock.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		if exprStmt.X == currCall.call {
			return exprStmt, true
		}
	}

	return nil, false
}

func lineIndentAtPos(pass *analysis.Pass, pos token.Pos, fileContentCache map[string][]byte) string {
	p := pass.Fset.PositionFor(pos, false)
	if p.Filename == "" || p.Offset < 0 {
		return defaultIndent
	}

	content, ok := fileContentCache[p.Filename]
	if !ok {
		var err error
		content, err = os.ReadFile(p.Filename)
		if err != nil {
			return defaultIndent
		}
		fileContentCache[p.Filename] = content
	}
	if p.Offset > len(content) {
		return defaultIndent
	}

	lineStart := p.Offset
	for lineStart > 0 {
		b := content[lineStart-1]
		if b == '\n' || b == '\r' {
			break
		}
		lineStart--
	}

	lineIndentEnd := lineStart
	for lineIndentEnd < len(content) {
		b := content[lineIndentEnd]
		if b != ' ' && b != '\t' {
			break
		}
		lineIndentEnd++
	}

	if lineIndentEnd == lineStart {
		return defaultIndent
	}

	return string(content[lineStart:lineIndentEnd])
}
