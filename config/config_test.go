package config_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/config"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		cfg     string
		expCfg  config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: `
enabled-checkers:
  - error-is
  - expected-actual
  - another-yet-checker
expected-actual:
  exp-var-pattern: ^want$
`,
			expCfg: config.Config{
				EnabledCheckers: []string{"error-is", "expected-actual", "another-yet-checker"},
				ExpectedActual: config.ExpectedActualConfig{
					ExpVarPattern: config.Regexp{regexp.MustCompile(`^want$`)},
				},
			},
		},
		{
			name: "empty config",
			cfg: `
enabled-checkers:
expected-actual:
`,
			expCfg: config.Config{},
		},
		{
			name: "invalid yaml",
			cfg: `
enabled-checkers:
  - error-is
  -- expected-actual
`,
			wantErr: true,
		},
		{
			name: "invalid expected-actual config",
			cfg: `
expected-actual:
  exp-var-pattern: ^want($
`,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			parsedCfg, err := config.Parse(strings.NewReader(tt.cfg))

			if tt.wantErr {
				if nil == err {
					t.Fatalf("no error but expected, value: %+v", parsedCfg)
				}
			} else {
				if !reflect.DeepEqual(tt.expCfg, parsedCfg) {
					t.Fatal("config parsed incorrectly")
				}
			}
		})
	}
}

func TestParseFromFile_InvalidPath(t *testing.T) {
	_, err := config.ParseFromFile("unknown")
	if err == nil {
		t.Fatal("no error but expected")
	}
}
