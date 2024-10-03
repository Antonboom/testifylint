package analysisutil_test

import (
	"testing"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

func TestIsJSONLike(t *testing.T) {
	cases := []struct {
		in       string
		expected bool
	}{
		{
			in:       `[{"name": "values-files", "array": ["values-dev.yaml"]}, {"name": "helm-parameters", "map": {"image.tag": "v1.2.3"}}]`,
			expected: true,
		},
		{
			in:       `{"labels":{"aaa":"111"},"annotations":{"ccc":"333"}}`,
			expected: true,
		},
		{
			in:       "{\"message\":\"No user was found in the LDAP server(s) with that username\"}",
			expected: true,
		},
		{
			in:       `"{\n  \"first\": \"Tobi\",\n  \"last\": \"Ferret\"\n}"`,
			expected: true,
		},
		{
			in:       `"{\"message\":\"No user was found in the LDAP server(s) with that username\"}"`,
			expected: true,
		},
		{
			in:       `{"uuid": "b65b1a22-db6d-4f5a-9b3d-7302368a82e6"}`,
			expected: true,
		},
		{
			in:       `apiVersion: 3`,
			expected: false,
		},
		{
			in:       `[{}]`,
			expected: false,
		},
		{
			in:       `{{ .TemplateVar }}`,
			expected: false,
		},
		{
			in:       `{{-.TemplateVar}}`,
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			isJSON := analysisutil.IsJSONLike(tt.in)
			if isJSON != tt.expected {
				t.FailNow()
			}
		})
	}
}
