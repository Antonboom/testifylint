package analysisutil

import (
	"encoding/json"
	"strconv"
	"strings"
)

// IsJSONObjectOrArray returns true if the string contains a valid JSON object or array.
func IsJSONObjectOrArray(s string) bool {
	if unquoted, err := strconv.Unquote(s); nil == err {
		s = unquoted
	}
	s = strings.TrimSpace(s)
	if s == "" || (s[0] != '{' && s[0] != '[') {
		return false
	}
	return json.Valid([]byte(s))
}
