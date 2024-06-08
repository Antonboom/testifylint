package analyzer //nolint:testpackage

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/internal/checkers"
	"github.com/Antonboom/testifylint/internal/config"
)

func Test_newCheckers(t *testing.T) {
	expVarPattern := regexp.MustCompile(`^expected[A-Z0-9].*`)
	fnPattern := regexp.MustCompile(`^NoErrorf?$`)

	enabledByDefaultRegularCheckers := []checkers.RegularChecker{
		checkers.NewFloatCompare(),
		checkers.NewBoolCompare(),
		checkers.NewEmpty(),
		checkers.NewLen(),
		checkers.NewNegativePositive(),
		checkers.NewCompares(),
		checkers.NewErrorNil(),
		checkers.NewNilCompare(),
		checkers.NewErrorIsAs(),
		checkers.NewExpectedActual(),
		checkers.NewSuiteExtraAssertCall(),
		checkers.NewSuiteDontUsePkg(),
		checkers.NewUselessAssert(),
		checkers.NewFormatter(),
	}
	allRegularCheckers := []checkers.RegularChecker{
		checkers.NewFloatCompare(),
		checkers.NewBoolCompare(),
		checkers.NewEmpty(),
		checkers.NewLen(),
		checkers.NewNegativePositive(),
		checkers.NewCompares(),
		checkers.NewErrorNil(),
		checkers.NewNilCompare(),
		checkers.NewErrorIsAs(),
		checkers.NewExpectedActual(),
		checkers.NewSuiteExtraAssertCall(),
		checkers.NewSuiteDontUsePkg(),
		checkers.NewUselessAssert(),
		checkers.NewFormatter(),
	}

	enabledByDefaultAdvancedCheckers := []checkers.AdvancedChecker{
		checkers.NewBlankImport(),
		checkers.NewGoRequire(),
		checkers.NewRequireError(),
		checkers.NewSuiteBrokenParallel(),
		checkers.NewSuiteSubtestRun(),
	}
	allAdvancedCheckers := []checkers.AdvancedChecker{
		checkers.NewBlankImport(),
		checkers.NewGoRequire(),
		checkers.NewRequireError(),
		checkers.NewSuiteBrokenParallel(),
		checkers.NewSuiteSubtestRun(),
		checkers.NewSuiteTHelper(),
	}

	zeroedFormatter := checkers.RegularChecker(checkers.NewFormatter().
		SetCheckFormatString(false).
		SetRequireFFuncs(false))

	cases := []struct {
		name        string
		cfg         config.Config
		expRegular  []checkers.RegularChecker
		expAdvanced []checkers.AdvancedChecker
	}{
		{
			name:        "no config",
			cfg:         config.Config{},
			expRegular:  replace(enabledByDefaultRegularCheckers, zeroedFormatter),
			expAdvanced: enabledByDefaultAdvancedCheckers,
		},
		{
			name:        "default config",
			cfg:         config.NewDefault(),
			expRegular:  enabledByDefaultRegularCheckers,
			expAdvanced: enabledByDefaultAdvancedCheckers,
		},
		{
			name: "enable two checkers only",
			cfg: config.Config{
				DisableAll: true,
				EnabledCheckers: config.KnownCheckersValue{
					checkers.NewRequireError().Name(),
					checkers.NewLen().Name(),
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewLen(),
			},
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewRequireError(),
			},
		},
		{
			name: "disable two checkers only",
			cfg: config.Config{
				EnableAll: true,
				DisabledCheckers: config.KnownCheckersValue{
					checkers.NewRequireError().Name(),
					checkers.NewSuiteTHelper().Name(),
				},
			},
			expRegular: filter(replace(allRegularCheckers, zeroedFormatter), config.KnownCheckersValue{
				checkers.NewRequireError().Name(),
				checkers.NewSuiteTHelper().Name(),
			}),
			expAdvanced: filter(allAdvancedCheckers, config.KnownCheckersValue{
				checkers.NewRequireError().Name(),
				checkers.NewSuiteTHelper().Name(),
			}),
		},
		{
			name: "enable one checker in addition to enabled by default checkers",
			cfg: config.Config{
				EnabledCheckers: config.KnownCheckersValue{
					checkers.NewSuiteTHelper().Name(),
				},
			},
			expRegular:  replace(enabledByDefaultRegularCheckers, zeroedFormatter),
			expAdvanced: allAdvancedCheckers,
		},
		{
			name: "disable three checkers from enabled by default checkers",
			cfg: config.Config{
				DisabledCheckers: config.KnownCheckersValue{
					checkers.NewNilCompare().Name(),
					checkers.NewErrorNil().Name(),
					checkers.NewRequireError().Name(),
				},
			},
			expRegular: filter(replace(enabledByDefaultRegularCheckers, zeroedFormatter), config.KnownCheckersValue{
				checkers.NewNilCompare().Name(),
				checkers.NewErrorNil().Name(),
				checkers.NewRequireError().Name(),
			}),
			expAdvanced: filter(enabledByDefaultAdvancedCheckers, config.KnownCheckersValue{
				checkers.NewNilCompare().Name(),
				checkers.NewErrorNil().Name(),
				checkers.NewRequireError().Name(),
			}),
		},
		{
			name: "bool-compare ignore custom types",
			cfg: config.Config{
				DisableAll:      true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewBoolCompare().Name()},
				BoolCompare: config.BoolCompareConfig{
					IgnoreCustomTypes: true,
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewBoolCompare().SetIgnoreCustomTypes(true),
			},
			expAdvanced: []checkers.AdvancedChecker{},
		},
		{
			name: "expected-actual pattern defined",
			cfg: config.Config{
				DisableAll:      true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewExpectedActual().Name()},
				ExpectedActual: config.ExpectedActualConfig{
					ExpVarPattern: config.RegexpValue{Regexp: expVarPattern},
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewExpectedActual().SetExpVarPattern(expVarPattern),
			},
			expAdvanced: []checkers.AdvancedChecker{},
		},
		{
			name: "go-require ignore http handlers",
			cfg: config.Config{
				DisableAll:      true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewGoRequire().Name()},
				GoRequire: config.GoRequireConfig{
					IgnoreHTTPHandlers: true,
				},
			},
			expRegular: []checkers.RegularChecker{},
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewGoRequire().SetIgnoreHTTPHandlers(true),
			},
		},
		{
			name: "require-equal fn pattern defined",
			cfg: config.Config{
				DisableAll:      true,
				EnabledCheckers: config.KnownCheckersValue{checkers.NewRequireError().Name()},
				RequireError: config.RequireErrorConfig{
					FnPattern: config.RegexpValue{Regexp: fnPattern},
				},
			},
			expRegular: []checkers.RegularChecker{},
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewRequireError().SetFnPattern(fnPattern),
			},
		},
		{
			name: "suite-extra-assert-call mode defined",
			cfg: config.Config{
				DisableAll:      true,
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
				t.Fatalf("unexpected regular checkers:\n%v\n\n!=\n\n%v\n\ncheck the checkers config!",
					checkerNames(rc), checkerNames(tt.expRegular))
			}
			if !reflect.DeepEqual(tt.expAdvanced, ac) {
				t.Fatalf("unexpected expAdvanced checkers:\n%v\n\n!=\n\n%v\n\ncheck the checkers config!",
					checkerNames(ac), checkerNames(tt.expAdvanced))
			}
		})
	}
}

func Test_newCheckers_invalidConfig(t *testing.T) {
	_, _, err := newCheckers(config.Config{EnableAll: true, DisableAll: true})
	if nil == err {
		t.Fatal("no error but expected")
	}
}

func Test_newCheckers_unknownChecker(t *testing.T) {
	_, _, err := newCheckers(config.Config{EnabledCheckers: config.KnownCheckersValue{"unknown"}})
	if nil == err {
		t.Fatal("no error but expected")
	}
}

func filter[T checkers.Checker](in []T, exclude config.KnownCheckersValue) []T {
	result := make([]T, 0)
	for _, v := range in {
		if exclude.Contains(v.Name()) {
			continue
		}
		result = append(result, v)
	}
	return result
}

func replace[T checkers.Checker](in []T, with T) []T {
	result := make([]T, 0)
	for _, v := range in {
		if v.Name() == with.Name() {
			result = append(result, with)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func checkerNames[T checkers.Checker](in []T) string {
	var b strings.Builder
	for i, c := range in {
		b.WriteString(c.Name())
		if i != len(in)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
