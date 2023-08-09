package config_test

import (
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/config"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     string
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: `
enabled-checkers:
  - suite-dont-use-pkg
  - len
  - require-error
expected-actual:
  pattern: ^want$
`,
			wantErr: false,
		},
		{
			name: "no checkers enable",
			cfg: `
enabled-checkers:

expected-actual:
  pattern: ^want$
`,
			wantErr: false,
		},
		{
			name: "no expected-actual section",
			cfg: `
enabled-checkers:
  - suite-dont-use-pkg
`,
			wantErr: false,
		},
		{
			name: "no expected-actual pattern",
			cfg: `
enabled-checkers:
  - suite-dont-use-pkg

expected-actual:
  pattern:
`,
			wantErr: false,
		},
		{
			name: "unknown checker",
			cfg: `
enabled-checkers:
  - suite-dont-use-pkg
  - len
  - require-party-for-everybody
`,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Parse(strings.NewReader(tt.cfg))
			if err != nil {
				t.Fatal(err)
			}

			err = config.Validate(cfg)
			if tt.wantErr && err == nil {
				t.Fatal("no error but expected")
			}
		})
	}
}
