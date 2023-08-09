package analyzer //nolint:testpackage

import (
	"github.com/Antonboom/testifylint/internal/checkers"
	"reflect"
	"testing"

	"github.com/Antonboom/testifylint/config"
)

func Test_newCheckers(t *testing.T) {
	cases := []struct {
		name         string
		cfg          config.Config
		callCheckers []string
		advCheckers  []string
	}{
		{
			name: "enable one only",
			cfg: config.Config{
				EnabledCheckers: []string{checkers.NewBoolCompare().Name()},
			},
			callCheckers: []string{checkers.BoolCompareCheckerName},
			advCheckers:  []string{},
		},
		//{
		//	name: "no config",
		//	cfg:  config.Config{},
		//	callCheckers: []string{
		//		"bool-compare",
		//		"float-compare",
		//		"empty",
		//		"len",
		//		"compares",
		//		"error",
		//		"error-is",
		//		"require-error",
		//		"expected-actual",
		//		"suite-dont-use-pkg",
		//	},
		//	advCheckers: []string{},
		//},
		//{
		//	name: "disable all",
		//	cfg: config.Config{
		//		Checkers: config.CheckersConfig{
		//			DisableAll: true,
		//		},
		//	},
		//	callCheckers: []string{},
		//	advCheckers:  []string{},
		//},
		//{
		//	name: "disable all enable pair",
		//	cfg: config.Config{
		//		Checkers: config.CheckersConfig{
		//			DisableAll: true,
		//			Enable:     []string{"error", "error-is"},
		//		},
		//	},
		//	callCheckers: []string{"error", "error-is"},
		//	advCheckers:  []string{},
		//},
		//{
		//	name: "disable of enabled by default",
		//	cfg: config.Config{
		//		Checkers: config.CheckersConfig{
		//			DisableAll: false,
		//			Disable:    []string{"bool-compare", "len", "require-error"},
		//		},
		//	},
		//	callCheckers: []string{
		//		"float-compare",
		//		"empty",
		//		"compares",
		//		"error",
		//		"error-is",
		//		"expected-actual",
		//		"suite-dont-use-pkg",
		//	},
		//	advCheckers: []string{},
		//},
		//{
		//	name: "enable and disable after",
		//	cfg: config.Config{
		//		Checkers: config.CheckersConfig{
		//			DisableAll: true,
		//			Enable:     []string{"bool-compare", "len", "require-error"},
		//			Disable:    []string{"bool-compare", "len"},
		//		},
		//	},
		//	callCheckers: []string{"require-error"},
		//	advCheckers:  []string{},
		//},
		//{
		//	name: "enable advanced checkers",
		//	cfg: config.Config{
		//		Checkers: config.CheckersConfig{
		//			DisableAll: true,
		//			Enable:     []string{"bool-compare", "len", "suite-thelper"},
		//		},
		//	},
		//	callCheckers: []string{"bool-compare", "len"},
		//	advCheckers:  []string{"suite-thelper"},
		//},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cc, ac, err := newCheckers(tt.cfg)
			if err != nil {
				t.Fatal(err)
			}

			ccNames, acNames := checkersNames(cc), checkersNames(ac)

			if !reflect.DeepEqual(ccNames, tt.callCheckers) {
				t.Fatalf("unexpected call checkers: %#v", ccNames)
			}
			if !reflect.DeepEqual(acNames, tt.advCheckers) {
				t.Fatalf("unexpected advanced checkers: %#v", acNames)
			}
		})
	}
}

func checkersNames[T interface{ Name() string }](ch []T) []string {
	res := make([]string, len(ch))
	for i, c := range ch {
		res[i] = c.Name()
	}
	return res
}
