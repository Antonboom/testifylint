package checkers

import (
	"go/ast"
	"regexp"
	"slices"

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
	if tv, ok := pass.TypesInfo.Types[e]; ok && tv.Value != nil {
		return analysisutil.IsJSONLike(tv.Value.String())
	}

	if id, ok := e.(*ast.Ident); ok && jsonIdentRe.MatchString(id.Name) {
		if hasWordAfterPattern(id.Name, jsonNegativeWordRe) {
			return false
		}
		return hasBytesType(pass, e) || hasStringType(pass, e)
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
	return slices.ContainsFunc(splitIntoWords(s), re.MatchString)
}

func splitIntoWords(s string) []string {
	return wordsRe.FindAllString(s, -1)
}
