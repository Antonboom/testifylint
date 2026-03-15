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
			// Valid JSON array of objects is now detected.
			in:       `[{}]`,
			expected: true,
		},
		{
			in:       `apiVersion: 3`,
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

func TestIsYAMLLike(t *testing.T) {
	cases := []struct {
		in       string
		expected bool
	}{
		{
			// Multi-line YAML with multiple key-value pairs.
			in: `
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
images:
  - name: foo
    newName: bar
`,
			expected: true,
		},
		{
			// Simple multi-line YAML.
			in:       "name: John\nage: 30\ncity: New York\n",
			expected: true,
		},
		{
			// Single-line YAML is not detected (could be any plain string).
			in:       "kind: Kustomization",
			expected: false,
		},
		{
			// Valid JSON objects are handled by JSONEq, not YAMLEq.
			in:       `{"foo": "bar"}`,
			expected: false,
		},
		{
			// Valid JSON arrays are handled by JSONEq, not YAMLEq.
			in:       `[{"name": "foo"}]`,
			expected: false,
		},
		{
			// Multi-line string with fewer than 2 key-value lines is not detected.
			in:       "kind: Kustomization\n",
			expected: false,
		},
		{
			// Quoted multi-line YAML (as it appears in Go source).
			in:       "`\nkind: Kustomization\napiVersion: v1\nimages:\n  - name: foo\n    newName: bar\n`",
			expected: true,
		},
		{
			// Regular text without YAML structure.
			in:       "hello world\nfoo bar\n",
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			isYAML := analysisutil.IsYAMLLike(tt.in)
			if isYAML != tt.expected {
				t.Fatalf("IsYAMLLike(%q) = %v, want %v", tt.in, isYAML, tt.expected)
			}
		})
	}
}
