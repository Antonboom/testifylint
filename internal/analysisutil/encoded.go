package analysisutil

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

var yamlKeyValueLineRe = regexp.MustCompile(`^\s*\w[\w.-]*\s*:(\s|$)`)

// IsJSONLike returns true if the string contains a valid JSON object or array.
func IsJSONLike(s string) bool {
	if unquoted, err := strconv.Unquote(s); err == nil {
		s = unquoted
	}
	s = strings.TrimSpace(s)
	if s == "" || (s[0] != '{' && s[0] != '[') {
		return false
	}
	return json.Valid([]byte(s))
}

// IsYAMLLike returns true if the string appears to be a multi-line YAML document.
// It returns false for strings that are valid JSON (which should use JSONEq instead).
func IsYAMLLike(s string) bool {
	if unquoted, err := strconv.Unquote(s); err == nil {
		s = unquoted
	}

	// Valid JSON is a subset of YAML; those comparisons should use JSONEq.
	s = strings.TrimSpace(s)
	if json.Valid([]byte(s)) {
		return false
	}

	// Only detect multi-line YAML documents.
	if !strings.Contains(s, "\n") {
		return false
	}

	kvCount := 0
	for _, line := range strings.Split(s, "\n") {
		if yamlKeyValueLineRe.MatchString(line) {
			kvCount++
		}
	}
	return kvCount >= 2
}
