package checkers

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

const httpMultipleReport = "use httptest.NewRecorder() instead of multiple HTTP assertions for the same handler call"

// HTTPMultiple detects situations like
//
//	assert.HTTPStatusCode(t, handler, "GET", "/path", nil, 200)
//	assert.HTTPBodyContains(t, handler, "GET", "/path", nil, "hello")
//
// and requires
//
//	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/path", nil)
//	w := httptest.NewRecorder()
//	handler(w, req)
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Contains(t, w.Body.String(), "hello")
type HTTPMultiple struct{}

// NewHTTPMultiple constructs HTTPMultiple checker.
func NewHTTPMultiple() HTTPMultiple { return HTTPMultiple{} }
func (HTTPMultiple) Name() string   { return "http-multiple" }

func (checker HTTPMultiple) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil), (*ast.FuncLit)(nil)}, func(node ast.Node) {
		var body *ast.BlockStmt
		switch n := node.(type) {
		case *ast.FuncDecl:
			body = n.Body
		case *ast.FuncLit:
			body = n.Body
		}
		if body == nil {
			return
		}

		diagnostics = append(diagnostics, checker.checkBody(pass, body)...)
	})

	return diagnostics
}

type httpCallKey struct {
	handler, method, url, values string
}

// callInStmt pairs a call with its parent-statement index in the enclosing block
// and records whether it is a direct ast.ExprStmt at that top level.
type callInStmt struct {
	call             *CallMeta
	stmtIdx          int
	isDirectExprStmt bool
}

func (checker HTTPMultiple) checkBody(pass *analysis.Pass, body *ast.BlockStmt) []analysis.Diagnostic {
	groups := make(map[httpCallKey][]callInStmt)

	// Iterate over each top-level statement to track statement indices
	// and determine whether a call is a direct ExprStmt (required for safe fix generation).
	for i, stmt := range body.List {
		ast.Inspect(stmt, func(node ast.Node) bool {
			if node == nil {
				return false
			}
			// Don't cross function literal boundaries; they form independent scopes.
			if _, ok := node.(*ast.FuncLit); ok {
				return false
			}
			ce, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			call := NewCallMeta(pass, ce)
			if call == nil {
				return true
			}
			if !isHTTPAssertion(call) {
				return true
			}
			// HTTP assertions have at least 4 args after t: handler, method, url, values.
			if len(call.Args) < 4 {
				return true
			}
			key := httpCallKey{
				handler: analysisutil.NodeString(pass.Fset, call.Args[0]),
				method:  analysisutil.NodeString(pass.Fset, call.Args[1]),
				url:     analysisutil.NodeString(pass.Fset, call.Args[2]),
				values:  analysisutil.NodeString(pass.Fset, call.Args[3]),
			}
			// A call is only eligible for a fix when it is the direct expression of
			// a top-level ExprStmt — not nested inside an if/for/select/etc. body.
			isDirectExpr := false
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				isDirectExpr = exprStmt.X == ce
			}
			groups[key] = append(groups[key], callInStmt{
				call:             call,
				stmtIdx:          i,
				isDirectExprStmt: isDirectExpr,
			})
			return true
		})
	}

	var diagnostics []analysis.Diagnostic
	for key, calls := range groups {
		if len(calls) < 2 {
			continue
		}
		sort.Slice(calls, func(i, j int) bool {
			return calls[i].call.Pos() < calls[j].call.Pos()
		})

		// A fix is only generated when:
		//   1. All calls sit in consecutive top-level statements.
		//   2. Every call is the direct expression of its top-level ExprStmt (no nesting).
		fixEligible := true
		for i := 1; i < len(calls); i++ {
			if calls[i].stmtIdx != calls[i-1].stmtIdx+1 {
				fixEligible = false
				break
			}
		}
		if fixEligible {
			for _, cis := range calls {
				if !cis.isDirectExprStmt {
					fixEligible = false
					break
				}
			}
		}

		var fix *analysis.SuggestedFix
		if fixEligible {
			fix = checker.generateFix(pass, body, key, calls)
		}

		for i, cis := range calls {
			if i == 0 && fix != nil {
				d := newDiagnostic(checker.Name(), cis.call, httpMultipleReport, *fix)
				diagnostics = append(diagnostics, *d)
			} else {
				d := newDiagnostic(checker.Name(), cis.call, httpMultipleReport)
				diagnostics = append(diagnostics, *d)
			}
		}
	}
	return diagnostics
}

func (checker HTTPMultiple) generateFix(
	pass *analysis.Pass,
	body *ast.BlockStmt,
	key httpCallKey,
	calls []callInStmt,
) *analysis.SuggestedFix {
	// All calls must be package calls (not object method calls) so we can safely
	// extract the TestingT variable name from ArgsRaw[0].
	// All calls must also use the same selector package (all assert or all require)
	// to avoid silently changing assertion semantics.
	pkg := calls[0].call.SelectorXStr
	for _, cis := range calls {
		if !cis.call.IsPkg() || cis.call.SelectorXStr != pkg {
			return nil
		}
	}

	firstStmt := body.List[calls[0].stmtIdx]

	// Resolve the local qualifier for net/http/httptest and get an optional import
	// TextEdit if the package is not yet imported. Uses addImportFix so that the fix
	// is offered regardless of whether httptest is already imported.
	httptestQualifier, httptestImportEdit, ok := addImportFix(pass.Files, firstStmt.Pos(), "net/http/httptest")
	if !ok {
		// Blank import or all candidate names exhausted — skip fix.
		return nil
	}

	// Resolve the local qualifier for net/http to remap raw string/int literals
	// (e.g. "GET" → http.MethodGet, 200 → http.StatusOK) in the generated fix.
	httpQual, httpImportEdit, httpAvail := httpNetPkgName(pass, firstStmt.Pos())

	// Remap the method argument if it is a plain string literal (e.g. "GET" → http.MethodGet).
	methodArg := key.method
	if httpAvail {
		if remapped, ok := remapMethodArg(calls[0].call.Args[1], httpQual); ok {
			methodArg = remapped
		}
	}

	// Extract the TestingT variable name from the first call to use t.Context()
	// in NewRequestWithContext. This avoids importing "context" entirely.
	tVar := analysisutil.NodeString(pass.Fset, calls[0].call.ArgsRaw[0])

	// Derive indentation from the column of the first statement.
	// Column is 1-indexed and counts bytes, so col-1 = number of leading tab/space bytes.
	// Go source files always use tabs (enforced by gofmt), so this is a safe assumption.
	col := pass.Fset.Position(firstStmt.Pos()).Column
	indent := strings.Repeat("\t", col-1)
	innerIndent := indent + "\t"

	// Collect replacement assertion lines for every call in the group.
	// The TestingT variable name and package (assert/require) are derived per-call
	// inside httpAssertionReplacement to use each call's original expressions.
	var assertLines []string
	for _, cis := range calls {
		lines := httpAssertionReplacement(pass, cis.call, httpQual, httpAvail)
		if lines == nil {
			return nil // Cannot auto-fix this particular assertion type.
		}
		assertLines = append(assertLines, lines...)
	}

	// Build the replacement block: a scoped { } avoids variable re-declaration
	// when multiple groups are fixed in the same function at once.
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(innerIndent)
	// Use t.Context() so the request context is automatically cancelled when the test ends,
	// without requiring an explicit "context" import.
	_, _ = fmt.Fprintf(&sb, "req := %s.NewRequestWithContext(%s.Context(), %s, %s, nil)\n",
		httptestQualifier, tVar, methodArg, key.url)
	if key.values != "nil" {
		// testify's HTTP helpers set req.URL.RawQuery = values.Encode() — mirror that here.
		sb.WriteString(innerIndent)
		_, _ = fmt.Fprintf(&sb, "req.URL.RawQuery = %s.Encode()\n", key.values)
	}
	sb.WriteString(innerIndent)
	_, _ = fmt.Fprintf(&sb, "rr := %s.NewRecorder()\n", httptestQualifier)
	sb.WriteString(innerIndent)
	sb.WriteString(key.handler + "(rr, req)\n")
	for _, line := range assertLines {
		sb.WriteString(innerIndent)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	sb.WriteString(indent)
	sb.WriteString("}")

	// Extend the end of the edit range to the start of the next line so that any
	// trailing inline comment (e.g. // want "...") on the last statement is also removed.
	lastStmt := body.List[calls[len(calls)-1].stmtIdx]
	f := pass.Fset.File(lastStmt.End())
	lastLine := f.Line(lastStmt.End())
	var endPos token.Pos
	if lastLine < f.LineCount() {
		endPos = f.LineStart(lastLine + 1)
		sb.WriteString("\n") // replace the consumed line-ending
	} else {
		endPos = token.Pos(f.Base() + f.Size())
	}

	textEdits := []analysis.TextEdit{{
		Pos:     firstStmt.Pos(),
		End:     endPos,
		NewText: []byte(sb.String()),
	}}
	if httptestImportEdit != nil {
		textEdits = append(textEdits, *httptestImportEdit)
	}
	if httpImportEdit != nil {
		textEdits = append(textEdits, *httpImportEdit)
	}

	return &analysis.SuggestedFix{
		Message:   "Use httptest.NewRecorder()",
		TextEdits: textEdits,
	}
}

// httpAssertionReplacement maps one HTTP testify assertion to its httptest equivalent line(s).
// call.Args layout (after t is stripped): [handler, method, url, values, <assertion-specific args...>]
// The TestingT variable name and package (assert/require) are extracted from call itself.
// httpQual is the local qualifier for "net/http" (empty for dot-import), httpAvail indicates
// whether net/http constants can be used in the generated code.
// Returns nil when the assertion cannot be automatically fixed.
func httpAssertionReplacement(pass *analysis.Pass, call *CallMeta, httpQual string, httpAvail bool) []string {
	// Use the actual TestingT expression from the original call instead of hard-coding "t".
	t := analysisutil.NodeString(pass.Fset, call.ArgsRaw[0])
	pkg := call.SelectorXStr
	extra := call.Args[4:] // args after handler/method/url/values
	fSuffix := ""
	if call.Fn.IsFmt {
		fSuffix = "f"
	}

	argsStr := func(args []ast.Expr) string {
		if len(args) == 0 {
			return ""
		}
		parts := make([]string, len(args))
		for i, a := range args {
			parts[i] = analysisutil.NodeString(pass.Fset, a)
		}
		return ", " + strings.Join(parts, ", ")
	}

	// statusConst returns the qualified net/http constant reference (e.g. "http.StatusBadRequest")
	// when httpAvail is true, or the bare constName for dot-imports.
	// Falls back to the constName without qualifier when net/http is unavailable.
	statusConst := func(constName string) string {
		if !httpAvail {
			return constName
		}
		return httpNetQualifiedConst(httpQual, constName)
	}

	switch call.Fn.NameFTrimmed {
	case "HTTPStatusCode":
		if len(extra) < 1 {
			return nil
		}
		code := analysisutil.NodeString(pass.Fset, extra[0])
		// Remap a raw integer literal (e.g. 200) to the corresponding http.StatusXxx constant.
		if httpAvail {
			if constName, ok := remapStatusCodeArg(extra[0]); ok {
				code = statusConst(constName)
			}
		}
		msg := argsStr(extra[1:])
		return []string{fmt.Sprintf("%s.Equal%s(%s, %s, rr.Code%s)", pkg, fSuffix, t, code, msg)}

	case "HTTPBodyContains":
		if len(extra) < 1 {
			return nil
		}
		str := analysisutil.NodeString(pass.Fset, extra[0])
		msg := argsStr(extra[1:])
		return []string{fmt.Sprintf("%s.Contains%s(%s, rr.Body.String(), %s%s)", pkg, fSuffix, t, str, msg)}

	case "HTTPBodyNotContains":
		if len(extra) < 1 {
			return nil
		}
		str := analysisutil.NodeString(pass.Fset, extra[0])
		msg := argsStr(extra[1:])
		return []string{fmt.Sprintf("%s.NotContains%s(%s, rr.Body.String(), %s%s)", pkg, fSuffix, t, str, msg)}

	case "HTTPError":
		msg := argsStr(extra)
		return []string{fmt.Sprintf("%s.GreaterOrEqual%s(%s, rr.Code, %s%s)", pkg, fSuffix, t, statusConst("StatusBadRequest"), msg)}

	case "HTTPSuccess":
		msg := argsStr(extra)
		return []string{
			fmt.Sprintf("%s.GreaterOrEqual%s(%s, rr.Code, %s%s)", pkg, fSuffix, t, statusConst("StatusOK"), msg),
			fmt.Sprintf("%s.Less%s(%s, rr.Code, %s%s)", pkg, fSuffix, t, statusConst("StatusMultipleChoices"), msg),
		}

	case "HTTPRedirect":
		msg := argsStr(extra)
		return []string{
			fmt.Sprintf("%s.GreaterOrEqual%s(%s, rr.Code, %s%s)", pkg, fSuffix, t, statusConst("StatusMultipleChoices"), msg),
			fmt.Sprintf("%s.Less%s(%s, rr.Code, %s%s)", pkg, fSuffix, t, statusConst("StatusBadRequest"), msg),
		}
	}

	return nil // HTTPBody or unknown — skip fix.
}

func isHTTPAssertion(call *CallMeta) bool {
	switch call.Fn.NameFTrimmed {
	case "HTTPBody", "HTTPBodyContains", "HTTPBodyNotContains",
		"HTTPError", "HTTPRedirect", "HTTPStatusCode", "HTTPSuccess":
		return true
	}
	return false
}

// remapMethodArg returns the qualified net/http constant for a plain HTTP method string literal
// (e.g. "GET" → http.MethodGet). Returns the qualified constant and true when a remap is found.
func remapMethodArg(arg ast.Expr, httpQual string) (string, bool) {
	bl, ok := arg.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", false
	}
	unquoted, err := strconv.Unquote(bl.Value)
	if err != nil {
		return "", false
	}
	constName, found := httpMethod[strings.ToUpper(unquoted)]
	if !found {
		return "", false
	}
	return httpNetQualifiedConst(httpQual, constName), true
}

// remapStatusCodeArg returns the net/http constant name for a plain integer status code literal
// (e.g. 200 → "StatusOK"). Returns the constant name and true when a remap is found.
func remapStatusCodeArg(arg ast.Expr) (string, bool) {
	bl, ok := arg.(*ast.BasicLit)
	if !ok || bl.Kind != token.INT {
		return "", false
	}
	v := constant.MakeFromLiteral(bl.Value, token.INT, 0)
	if v.Kind() != constant.Int {
		return "", false
	}
	intVal, exact := constant.Int64Val(v)
	if !exact || int64(int(intVal)) != intVal {
		return "", false
	}
	constName, found := httpStatusCode[int(intVal)]
	if !found {
		return "", false
	}
	return constName, true
}
