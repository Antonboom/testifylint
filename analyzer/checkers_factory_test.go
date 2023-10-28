package analyzer //nolint:testpackage

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
	"github.com/Antonboom/testifylint/internal/config"
)

func Test_newCheckers(t *testing.T) {
	pattern := regexp.MustCompile(`^expected[A-Z0-9].*`)

	enabledByDefaultRegularCheckers := []checkers.RegularChecker{
		checkers.NewFloatCompare(),
		checkers.NewBoolCompare(),
		checkers.NewEmpty(),
		checkers.NewLen(),
		checkers.NewCompares(),
		checkers.NewErrorNil(),
		checkers.NewNilCompare(),
		checkers.NewErrorIsAs(),
		checkers.NewRequireError(),
		checkers.NewExpectedActual(),
		checkers.NewSuiteExtraAssertCall(),
		checkers.NewSuiteDontUsePkg(),
	}
	allRegularCheckers := []checkers.RegularChecker{
		checkers.NewFloatCompare(),
		checkers.NewBoolCompare(),
		checkers.NewEmpty(),
		checkers.NewLen(),
		checkers.NewCompares(),
		checkers.NewErrorNil(),
		checkers.NewNilCompare(),
		checkers.NewErrorIsAs(),
		checkers.NewRequireError(),
		checkers.NewExpectedActual(),
		checkers.NewSuiteExtraAssertCall(),
		checkers.NewSuiteDontUsePkg(),
	}

	cases := []struct {
		name        string
		cfg         config.Config
		expRegular  []checkers.RegularChecker
		expAdvanced []checkers.AdvancedChecker
	}{
		{
			name:        "no config",
			cfg:         config.Config{},
			expRegular:  enabledByDefaultRegularCheckers,
			expAdvanced: []checkers.AdvancedChecker{},
		},
		{
			name: "no enabled checkers",
			cfg: config.Config{
				EnabledCheckers: []string{},
			},
			expRegular:  enabledByDefaultRegularCheckers,
			expAdvanced: []checkers.AdvancedChecker{},
		},
		{
			name: "no enabled checkers but enable-all true",
			cfg: config.Config{
				EnabledCheckers: []string{},
				EnableAll:       true,
			},
			expRegular: allRegularCheckers,
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewSuiteTHelper(),
			},
		},
		{
			name: "enabled checkers defined",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{
					checkers.NewSuiteTHelper().Name(),
					checkers.NewRequireError().Name(),
					checkers.NewSuiteExtraAssertCall().Name(),
					checkers.NewLen().Name(),
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewLen(),
				checkers.NewRequireError(),
				checkers.NewSuiteExtraAssertCall(),
			},
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewSuiteTHelper(),
			},
		},
		{
			name: "enabled checkers defined but enable-all true",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{
					checkers.NewSuiteTHelper().Name(),
					checkers.NewRequireError().Name(),
					checkers.NewSuiteExtraAssertCall().Name(),
					checkers.NewLen().Name(),
				},
				EnableAll: true,
			},
			expRegular: allRegularCheckers,
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewSuiteTHelper(),
			},
		},
		{
			name: "expected-actual pattern defined",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{checkers.NewExpectedActual().Name()},
				ExpectedActual: config.ExpectedActualConfig{
					ExpVarPattern: config.RegexpValue{Regexp: pattern},
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewExpectedActual().SetExpVarPattern(pattern),
			},
			expAdvanced: []checkers.AdvancedChecker{},
		},
		{
			name: "suite-extra-assert-call mode defined",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{checkers.NewSuiteExtraAssertCall().Name()},
				SuiteExtraAssertCall: config.SuiteExtraAssertCallConfig{
					Mode: checkers.SuiteExtraAssertCallModeRequire,
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewSuiteExtraAssertCall().SetMode(checkers.SuiteExtraAssertCallModeRequire),
			},
			expAdvanced: []checkers.AdvancedChecker{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rc, ac, err := newCheckers(tt.cfg)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.expRegular, rc) {
				t.Fatalf("unexpected regular checkers: %#v", rc)
			}
			if !reflect.DeepEqual(tt.expAdvanced, ac) {
				t.Fatalf("unexpected expAdvanced checkers: %#v", ac)
			}
		})
	}
}

func Test_newCheckers_unknownChecker(t *testing.T) {
	_, _, err := newCheckers(config.Config{EnabledCheckers: config.KnownCheckersValue{"unknown"}})
	if nil == err {
		t.Fatal("no error but expected")
	}
}
