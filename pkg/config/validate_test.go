package config_test

import (
	"testing"

	"github.com/Antonboom/testifylint/pkg/config"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: false,
					Enable:     []string{"len", "empty", "expected-actual"},
					Disable:    []string{"error-is"},
				},
				ExpectedActual: config.ExpectedActualConfig{
					Pattern: `^want$`,
				},
			},
			wantErr: false,
		},
		{
			name: "disable all and disable",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: true,
					Enable:     []string{"error"},
					Disable:    []string{"error-is"},
				},
			},
			wantErr: true,
		},
		{
			name: "disable all and no enabled checkers",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: true,
					Enable:     []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown enabled checker",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: false,
					Enable:     []string{"bugaga"},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown disabled checker",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: false,
					Disable:    []string{"bugaga"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid expected-actual pattern",
			cfg: config.Config{
				Checkers: config.CheckersConfig{
					DisableAll: false,
					Disable:    []string{"error-is"},
				},
				ExpectedActual: config.ExpectedActualConfig{
					Pattern: `(bugaga`,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			t.Log(err)

			if tt.wantErr && err == nil {
				t.Fatal("no error but expected")
			}
		})
	}
}
