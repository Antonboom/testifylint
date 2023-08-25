package config_test

import (
	"reflect"
	"regexp"
	"slices"
	"testing"

	"github.com/Antonboom/testifylint/internal/config"
)

func TestKnownCheckersValue_Set(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected config.KnownCheckersValue
		wantErr  bool
	}{
		{
			name:     "positive",
			input:    "suite-dont-use-pkg,len,require-error",
			expected: config.KnownCheckersValue{"suite-dont-use-pkg", "len", "require-error"},
		},
		{
			name:    "unknown checker",
			input:   "suite-dont-use-pkg,lenlen,require-error",
			wantErr: true,
		},
		{
			name:    "malformed input 1",
			input:   "suite-dont-use-pkg, len, require-error",
			wantErr: true,
		},
		{
			name:    "malformed input 2",
			input:   "suite-dont-use-pkg,,len,require-error",
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var v config.KnownCheckersValue
			err := v.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}
			if !slices.Equal(tt.expected, v) {
				t.Fatal(v)
			}
		})
	}
}

func TestRegexpValue_String_ZeroValue(t *testing.T) {
	var r config.RegexpValue
	if r.String() != "" {
		t.Fatal()
	}
}

func TestRegexp_Set(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected config.RegexpValue
		wantErr  bool
	}{
		{
			name:     "valid regexp",
			input:    `^expected$`,
			expected: config.RegexpValue{regexp.MustCompile(`^expected$`)},
		},
		{
			name:    "invalid regexp",
			input:   `((.`,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var v config.RegexpValue
			err := v.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.expected, v) {
				t.Fatal(v)
			}
		})
	}
}
