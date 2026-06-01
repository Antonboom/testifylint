package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"sort"
	"strings"

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
		tArg                      string
		requireQualifier          string
		requireImportIsAccessible bool
	)
	if currCall.testifyCall.IsPkg && (len(currCall.testifyCall.ArgsRaw) >= 2) {
		rawTArg := currCall.testifyCall.ArgsRaw[0]
		if implementsTestingT(pass, rawTArg) {
			exprStmt, ok = findExprStmtForCall(currCall)
			tArg = analysisutil.NodeString(pass.Fset, rawTArg)
			requireQualifier, requireImportIsAccessible = analysisutil.LocalPkgName(
				pass.Files, currCall.call.Pos(), testify.RequirePkgPath)
		}
	}

	needImportEdit := !requireImportIsAccessible || requireQualifier == ""
	if ok && requireQualifier == "" {
		requireQualifier = availableRequireQualifier(pass.Files, currCall.call.Pos())
	}

	missingGuards := make([]guardRequirement, 0, 4)
	seenGuards := make(map[string]struct{}, 4)
	for _, target := range indexedAccesses(pass, currCall.call) {
		for _, guard := range missingGuardRequirements(
			pass, currCall, currCallIndex, otherCalls, target, requireQualifier, tArg,
		) {
			if _, seen := seenGuards[guard.textKey]; seen {
				continue
			}
			seenGuards[guard.textKey] = struct{}{}
			missingGuards = append(missingGuards, guard)
		}
	}
	if len(missingGuards) == 0 {
		return nil
	}

	diagnostic := newDiagnostic(checkerName, currCall.testifyCall, requireLenGuardReport)
	if !ok {
		return diagnostic
	}

	additionalEdits := make([]analysis.TextEdit, 0, 2)
	if needImportEdit {
		importEdit, hasImportEdit := addRequireImportTextEdit(pass, currCall.call.Pos(), requireQualifier)
		if !hasImportEdit {
			return diagnostic
		}
		additionalEdits = append(additionalEdits, importEdit)
	}

	indent := lineIndentAtPos(pass, currCall.call.Pos(), fileContentCache)
	var insertBuilder strings.Builder
	for i, guard := range missingGuards {
		if i > 0 {
			insertBuilder.WriteString(indent)
		}
		insertBuilder.WriteString(guard.lineText)
		insertBuilder.WriteByte('\n')
	}
	insertBuilder.WriteString(indent)

	return newDiagnostic(checkerName, currCall.testifyCall, requireLenGuardReport, analysis.SuggestedFix{
		Message: "Insert `require` guards",
		TextEdits: append(additionalEdits,
			analysis.TextEdit{
				Pos:     exprStmt.Pos(),
				End:     exprStmt.Pos(),
				NewText: []byte(insertBuilder.String()),
			},
		),
	})
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
	mapKeys    map[string]ast.Expr
}

type guardRequirement struct {
	textKey  string
	lineText string
}

func indexedAccesses(pass *analysis.Pass, node ast.Node) []indexedAccess {
	var result []indexedAccess
	indexByCollection := make(map[string]int)

	ast.Inspect(node, func(n ast.Node) bool {
		ie, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		collection := analysisutil.NodeString(pass.Fset, ie.X)
		if collection == "" {
			return true
		}

		tt := pass.TypesInfo.TypeOf(ie.X)
		if tt == nil {
			return true
		}
		_, isMap := tt.Underlying().(*types.Map)

		idx := -1
		if v, ok := isIntBasicLit(ie.Index); ok && v >= 0 {
			idx = v
		}
		if !isMap && idx < 0 {
			return true
		}

		i, exists := indexByCollection[collection]
		if !exists {
			indexByCollection[collection] = len(result)
			result = append(result, indexedAccess{collection: collection, maxIndex: idx})
			i = len(result) - 1
		}

		if idx > result[i].maxIndex {
			result[i].maxIndex = idx
		}

		if isMap {
			if result[i].mapKeys == nil {
				result[i].mapKeys = make(map[string]ast.Expr)
			}
			keyStr := analysisutil.NodeString(pass.Fset, ie.Index)
			if keyStr != "" {
				result[i].mapKeys[keyStr] = ie.Index
			}
		}

		return true
	})

	return result
}

func missingGuardRequirements(
	pass *analysis.Pass,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	target indexedAccess,
	requireQualifier string,
	tArg string,
) []guardRequirement {
	var reqs []guardRequirement

	if len(target.mapKeys) > 0 {
		mapKeys := make([]string, 0, len(target.mapKeys))
		for key := range target.mapKeys {
			mapKeys = append(mapKeys, key)
		}
		sort.Strings(mapKeys)

		for _, mapKey := range mapKeys {
			if hasContainsGuard(pass, currCall, currCallIndex, otherCalls, target.collection, mapKey) {
				continue
			}
			reqs = append(reqs, guardRequirement{
				textKey:  "contains:" + target.collection + ":" + mapKey,
				lineText: fmt.Sprintf("%s.Contains(%s, %s, %s)", requireQualifier, tArg, target.collection, mapKey),
			})
		}
		return reqs
	}

	requiredLen := target.maxIndex + 1
	if requiredLen <= 0 || hasLenGuard(pass, currCall, currCallIndex, otherCalls, target.collection, requiredLen) {
		return reqs
	}

	if requiredLen == 1 {
		reqs = append(reqs, guardRequirement{
			textKey:  "notempty:" + target.collection,
			lineText: fmt.Sprintf("%s.NotEmpty(%s, %s)", requireQualifier, tArg, target.collection),
		})
		return reqs
	}

	reqs = append(reqs, guardRequirement{
		textKey:  fmt.Sprintf("gte:%s:%d", target.collection, requiredLen),
		lineText: fmt.Sprintf("%s.GreaterOrEqual(%s, len(%s), %d)", requireQualifier, tArg, target.collection, requiredLen),
	})
	return reqs
}

func hasLenGuard(
	pass *analysis.Pass,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	collection string,
	requiredLen int,
) bool {
	for i := guardSearchStartIndex(pass, currCall, currCallIndex, otherCalls); i < currCallIndex; i++ {
		c := otherCalls[i]
		if c.parentBlock != currCall.parentBlock || c.testifyCall == nil || c.testifyCall.IsAssert {
			continue
		}

		switch c.testifyCall.Fn.NameFTrimmed {
		case "Len":
			if len(c.testifyCall.Args) < 2 {
				continue
			}
			lenCollection := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[0])
			assertedLen, isBasicLit := isIntBasicLit(c.testifyCall.Args[1])
			if lenCollection == collection && isBasicLit && assertedLen >= requiredLen {
				return true
			}
		case "NotEmpty":
			if requiredLen != 1 || len(c.testifyCall.Args) < 1 {
				continue
			}
			notEmptyCollection := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[0])
			if notEmptyCollection == collection {
				return true
			}
		case "GreaterOrEqual":
			if len(c.testifyCall.Args) < 2 {
				continue
			}
			lenArg, ok := isBuiltinLenCall(pass, c.testifyCall.Args[0])
			if !ok || analysisutil.NodeString(pass.Fset, lenArg) != collection {
				continue
			}
			assertedMinLen, isBasicLit := isIntBasicLit(c.testifyCall.Args[1])
			if isBasicLit && assertedMinLen >= requiredLen {
				return true
			}
		}
	}
	return false
}

func hasContainsGuard(
	pass *analysis.Pass,
	currCall *callMeta,
	currCallIndex int,
	otherCalls []*callMeta,
	collection string,
	key string,
) bool {
	for i := guardSearchStartIndex(pass, currCall, currCallIndex, otherCalls); i < currCallIndex; i++ {
		c := otherCalls[i]
		if c.parentBlock != currCall.parentBlock || c.testifyCall == nil || c.testifyCall.IsAssert {
			continue
		}
		if c.testifyCall.Fn.NameFTrimmed != "Contains" || len(c.testifyCall.Args) < 2 {
			continue
		}
		containsCollection := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[0])
		if containsCollection != collection {
			continue
		}
		containsKey := analysisutil.NodeString(pass.Fset, c.testifyCall.Args[1])
		if containsKey == key {
			return true
		}
	}
	return false
}

func guardSearchStartIndex(pass *analysis.Pass, currCall *callMeta, currCallIndex int, otherCalls []*callMeta) int {
	start := 0
	for i := range currCallIndex {
		c := otherCalls[i]
		if c.parentBlock != currCall.parentBlock || c.testifyCall == nil {
			continue
		}
		if isRequireLenErrorCheck(c.testifyCall.Fn.NameFTrimmed) || callContainsErrorArg(pass, c.testifyCall) {
			start = i + 1
		}
	}
	return start
}

func isRequireLenErrorCheck(nameFTrimmed string) bool {
	switch nameFTrimmed {
	case "Error", "ErrorIs", "ErrorAs", "EqualError", "ErrorContains", "NoError", "NotErrorIs":
		return true
	}
	return false
}

// callContainsErrorArg returns true if any non-message/args argument in the
// testify call has an error type. This covers cases like assert.Equal(t, err, nil)
// or assert.Nil(t, err) where the error is outside of msg and args.
func callContainsErrorArg(pass *analysis.Pass, cm *CallMeta) bool {
	limit := nonMessageArgLimit(pass, cm)
	for _, arg := range cm.Args[:limit] {
		t := pass.TypesInfo.TypeOf(arg)
		if t == nil {
			continue
		}
		if types.Implements(t, errorIface) || types.Implements(types.NewPointer(t), errorIface) {
			return true
		}
	}
	return false
}

// nonMessageArgLimit returns the number of non-message/args arguments in a
// testify call. Message/format arguments are at the end (variadic msgAndArgs
// or explicit msgfmt+args for formatted variants).
func nonMessageArgLimit(pass *analysis.Pass, cm *CallMeta) int {
	sig := cm.Fn.Signature
	if sig == nil {
		return len(cm.Args)
	}

	paramsLen := sig.Params().Len()
	if len(cm.ArgsRaw) > 0 && implementsTestingT(pass, cm.ArgsRaw[0]) {
		paramsLen--
	}
	if paramsLen < 0 {
		return 0
	}

	limit := len(cm.Args)
	if sig.Variadic() {
		// Variadic parameter is msgAndArgs; base assertion args are fixed params.
		limit = paramsLen - 1
		// Formatted assertions have explicit message arg before msgAndArgs.
		if cm.Fn.IsFmt {
			limit--
		}
	}
	if limit < 0 {
		return 0
	}
	if limit > len(cm.Args) {
		return len(cm.Args)
	}
	return limit
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
