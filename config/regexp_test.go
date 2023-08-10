package config_test

import (
	"github.com/Antonboom/testifylint/config"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRegexp_UnmarshalYAML(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected string
		wantErr  bool
	}{
		{
			name:    "invalid data format",
			in:      `[ 100 ]`,
			wantErr: true,
		},
		{
			name:    "invalid regexp",
			in:      `((.`,
			wantErr: true,
		},
		{
			name:     "valid regexp",
			in:       `^expected$`,
			expected: `^expected$`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var r config.Regexp
			err := yaml.Unmarshal([]byte(tt.in), &r)
			if tt.wantErr {
				if nil == err {
					t.Fatalf("no error but expected, value: %v", r.String())
				}
			} else {
				if tt.expected != r.String() {
					t.Fatal("regexp unmarshalled incorrectly")
				}
			}
		})
	}
}
