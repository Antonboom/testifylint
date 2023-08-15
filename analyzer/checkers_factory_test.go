package analyzer //nolint:testpackage

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/Antonboom/testifylint/config"
	"github.com/Antonboom/testifylint/internal/checkers"
)

func Test_newCheckers(t *testing.T) {
	pattern := regexp.MustCompile(`^expected[A-Z0-9].*`)

	enabledByDefaultRegularCheckers := []checkers.RegularChecker{
		checkers.NewBoolCompare(),
		checkers.NewFloatCompare(),
		checkers.NewEmpty(),
		checkers.NewLen(),
		checkers.NewCompares(),
		checkers.NewError(),
		checkers.NewErrorIs(),
		checkers.NewRequireError(),
		checkers.NewExpectedActual(),
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
			name: "enabled checkers defined",
			cfg: config.Config{
				EnabledCheckers: []string{
					checkers.NewSuiteTHelper().Name(),
					checkers.NewRequireError().Name(),
					checkers.NewSuiteNoExtraAssertCall().Name(),
					checkers.NewLen().Name(),
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewLen(),
				checkers.NewRequireError(),
				checkers.NewSuiteNoExtraAssertCall(),
			},
			expAdvanced: []checkers.AdvancedChecker{
				checkers.NewSuiteTHelper(),
			},
		},
		{
			name: "expected-actual pattern defined",
			cfg: config.Config{
				EnabledCheckers: []string{checkers.NewExpectedActual().Name()},
				ExpectedActual: config.ExpectedActualConfig{
					Pattern: config.Regexp{pattern},
				},
			},
			expRegular: []checkers.RegularChecker{
				checkers.NewExpectedActual().SetExpPattern(pattern),
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
	_, _, err := newCheckers(config.Config{EnabledCheckers: []string{"unknown"}})
	if nil == err {
		t.Fatal("no error but expected")
	}
}
