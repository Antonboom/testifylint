package config_test

import (
	"flag"
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
	"github.com/Antonboom/testifylint/internal/config"
)

func TestNewDefault(t *testing.T) {
	cfg := config.NewDefault()

	if cfg.EnableAll {
		t.Fatal()
	}
	if len(cfg.DisabledCheckers) != 0 {
		t.Fatal()
	}
	if cfg.DisableAll {
		t.Fatal()
	}
	if len(cfg.EnabledCheckers) != 0 {
		t.Fatal()
	}
	if cfg.BoolCompare.IgnoreCustomTypes {
		t.Fatal()
	}
	if cfg.ExpectedActual.ExpVarPattern.String() != checkers.DefaultExpectedVarPattern.String() {
		t.Fatal()
	}
	if !cfg.Formatter.CheckFormatString {
		t.Fatal()
	}
	if cfg.Formatter.RequireFFuncs {
		t.Fatal()
	}
	if cfg.GoRequire.IgnoreHTTPHandlers {
		t.Fatal()
	}
	if cfg.RequireError.FnPattern.String() != "" {
		t.Fatal()
	}
	if cfg.SuiteExtraAssertCall.Mode != checkers.SuiteExtraAssertCallModeRemove {
		t.Fatal()
	}
}

func TestConfig_Validate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		// Positive.
		{
			name:    "default config",
			cfg:     config.NewDefault(),
			wantErr: false,
		},
		{
			name: "enable-all and disable simultaneously",
			cfg: config.Config{
				EnableAll:        true,
				DisabledCheckers: config.KnownCheckersValue{checkers.NewErrorNil().Name()},
			},
			wantErr: false,
		},
		{
			name: "disable-all and enable simultaneously",
			cfg: config.Config{
				DisableAll:      true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewErrorNil().Name()},
			},
			wantErr: false,
		},
		{
			name: "enable some checkers",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{checkers.NewErrorNil().Name()},
			},
			wantErr: false,
		},
		{
			name: "disable some checkers",
			cfg: config.Config{
				DisabledCheckers: config.KnownCheckersValue{checkers.NewErrorNil().Name()},
			},
			wantErr: false,
		},
		{
			name: "enable and disable simultaneously different checkers",
			cfg: config.Config{
				DisabledCheckers: config.KnownCheckersValue{checkers.NewRequireError().Name()},
				EnabledCheckers:  config.KnownCheckersValue{checkers.NewErrorNil().Name()},
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "enable-all and disable-all simultaneously",
			cfg: config.Config{
				EnableAll:  true,
				DisableAll: true,
			},
			wantErr: true,
		},
		{
			name: "enable-all and enable simultaneously",
			cfg: config.Config{
				EnableAll:       true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewExpectedActual().Name()},
			},
			wantErr: true,
		},
		{
			name: "disable-all and disable simultaneously",
			cfg: config.Config{
				DisableAll:       true,
				DisabledCheckers: config.KnownCheckersValue{checkers.NewExpectedActual().Name()},
			},
			wantErr: true,
		},
		{
			name: "disable-all and no enable",
			cfg: config.Config{
				DisableAll: true,
			},
			wantErr: true,
		},
		{
			name: "enable and disable simultaneously the same checker",
			cfg: config.Config{
				EnabledCheckers:  config.KnownCheckersValue{checkers.NewNilCompare().Name(), checkers.NewExpectedActual().Name()},
				DisabledCheckers: config.KnownCheckersValue{checkers.NewExpectedActual().Name()},
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}
		})
	}
}

func TestBindToFlags(t *testing.T) {
	cfg := config.NewDefault()
	fs := flag.NewFlagSet("TestBindToFlags", flag.PanicOnError)

	config.BindToFlags(&cfg, fs)

	for flagName, defaultVal := range map[string]string{
		"enable-all":                       "false",
		"disable":                          "",
		"disable-all":                      "false",
		"enable":                           "",
		"bool-compare.ignore-custom-types": "false",
		"expected-actual.pattern":          cfg.ExpectedActual.ExpVarPattern.String(),
		"formatter.check-format-string":    "true",
		"formatter.require-f-funcs":        "false",
		"go-require.ignore-http-handlers":  "false",
		"require-error.fn-pattern":         cfg.RequireError.FnPattern.String(),
		"suite-extra-assert-call.mode":     "remove",
	} {
		t.Run(flagName, func(t *testing.T) {
			if v := fs.Lookup(flagName).DefValue; v != defaultVal {
				t.Fatal(v)
			}
		})
	}
}
