package checkers

import (
	"go/ast"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

var (
	wordsRe = regexp.MustCompile(`[A-Z]+(?:[a-z]*|$)|[a-z]+`) // NOTE(a.telyshev): ChatGPT.

	jsonIdentRe        = regexp.MustCompile(`json|JSON|Json`)
	jsonNegativeWordRe = regexp.MustCompile(`(?i)^(invalid|bad|malformed|broken|corrupt|wrong)$`)
	yamlWordRe         = regexp.MustCompile(`yaml|YAML|Yaml|^(yml|YML|Yml)$`)
)

func isJSONStyleExpr(pass *analysis.Pass, e ast.Expr) bool {
	if id, ok := e.(*ast.Ident); ok && jsonIdentRe.MatchString(id.Name) {
		// Skip identifiers that contain a "negative" word, e.g. "invalidJSON" or "badJSON".
		// Such variables likely hold non-valid JSON (e.g. in error-path tests) and suggesting
		// JSONEq would cause test failures at runtime.
		if hasWordAfterPattern(id.Name, jsonNegativeWordRe) {
			return false
		}

		// If this is a constant with a known value, verify it actually looks like JSON.
		// This avoids false positives for constants like `const contentTypeJSON = "application/json"`.
		if t, ok := pass.TypesInfo.Types[e]; ok && t.Value != nil {
			return analysisutil.IsJSONLike(t.Value.String())
		}
		return hasBytesType(pass, e) || hasStringType(pass, e)
	}

	if t, ok := pass.TypesInfo.Types[e]; ok && t.Value != nil {
		return analysisutil.IsJSONLike(t.Value.String())
	}

	if bl, ok := e.(*ast.BasicLit); ok {
		return bl.Kind == token.STRING && analysisutil.IsJSONLike(bl.Value)
	}

	if args, ok := isFmtSprintfCall(pass, e); ok {
		return isJSONStyleExpr(pass, args[0])
	}

	return false
}

func isYAMLStyleExpr(pass *analysis.Pass, e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && (hasBytesType(pass, e) || hasStringType(pass, e)) && hasWordAfterPattern(id.Name, yamlWordRe)
}

func hasWordAfterPattern(s string, re *regexp.Regexp) bool {
	for _, w := range splitIntoWords(s) {
		if re.MatchString(w) {
			return true
		}
	}
	return false
}

func splitIntoWords(s string) []string {
	return wordsRe.FindAllString(s, -1)
}
