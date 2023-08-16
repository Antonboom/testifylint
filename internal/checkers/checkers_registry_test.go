package checkers_test

import (
	"sort"
	"testing"

	"slices"

	"github.com/Antonboom/testifylint/internal/checkers"
)

func TestAll(t *testing.T) {
	checkerList := checkers.All()
	if len(checkerList) == 0 {
		t.Fatalf("no known checkers: empty list")
	}

	// NOTE(a.telyshev): I don't use constants or checker's Name() method on purpose.
	expected := []string{
		"bool-compare",
		"float-compare",
		"empty",
		"len",
		"compares",
		"error",
		"error-is",
		"require-error",
		"expected-actual",
		"suite-no-extra-assert-call",
		"suite-dont-use-pkg",
		"suite-thelper",
	}
	if !slices.Equal(expected, checkerList) {
		t.Fatalf("unexpected list: %#v", checkerList)
	}
}

func TestEnabledByDefault(t *testing.T) {
	checkerList := checkers.EnabledByDefault()
	if len(checkerList) == 0 {
		t.Fatalf("no enabled by default checkers: empty list")
	}

	// NOTE(a.telyshev): I don't use constants or checker's Name() method on purpose.
	expected := []string{
		"bool-compare",
		"float-compare",
		"empty",
		"len",
		"compares",
		"error",
		"error-is",
		"require-error",
		"expected-actual",
		"suite-dont-use-pkg",
	}
	if !slices.Equal(expected, checkerList) {
		t.Fatalf("unexpected list: %#v", checkerList)
	}
}

func TestGet(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		for _, name := range checkers.All() {
			checker, ok := checkers.Get(name)
			if !ok || checker == nil {
				t.Fatalf("lost %q checker", name)
			}

			if checker.Name() != name {
				t.Fatalf("invalid mapping for %q checker", name)
			}
		}
	})

	t.Run("unknown checker", func(t *testing.T) {
		checker, ok := checkers.Get("unknown")
		if ok || checker != nil {
			t.Fatalf("unexpected checker")
		}
	})

	t.Run("checker types", func(t *testing.T) {
		{
			ch, _ := checkers.Get(checkers.NewLen().Name())
			_, ok := ch.(checkers.RegularChecker)
			if !ok {
				t.Fatal("satisfaction of checkers.RegularChecker interface was lost")
			}
		}
		{
			ch, _ := checkers.Get(checkers.NewSuiteTHelper().Name())
			_, ok := ch.(checkers.AdvancedChecker)
			if !ok {
				t.Fatal("satisfaction of checkers.AdvancedChecker interface was lost")
			}
		}
	})
}

func TestIsKnown(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		for _, name := range checkers.All() {
			ok := checkers.IsKnown(name)
			if !ok {
				t.Fatalf("lost %q checker", name)
			}
		}
	})

	t.Run("unknown checker", func(t *testing.T) {
		ok := checkers.IsKnown("unknown")
		if ok {
			t.Fatalf("unexpected checker")
		}
	})
}

func TestIsEnabledByDefault(t *testing.T) {
	if !checkers.IsEnabledByDefault(checkers.NewBoolCompare().Name()) {
		t.FailNow()
	}
	if checkers.IsEnabledByDefault(checkers.NewSuiteTHelper().Name()) {
		t.FailNow()
	}
	if checkers.IsEnabledByDefault("unknown") {
		t.FailNow()
	}
}

func TestSortByPriority(t *testing.T) {
	checkerList := checkers.All()
	sort.Strings(checkerList)
	if slices.Equal(checkerList, checkers.All()) {
		t.Fatal("precondition failed")
	}

	checkers.SortByPriority(checkerList)
	if !slices.Equal(checkerList, checkers.All()) {
		t.FailNow()
	}
}
