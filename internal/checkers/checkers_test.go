package checkers_test

import (
	"fmt"
	"github.com/Antonboom/testifylint/internal/checkers"
	"reflect"
	"strings"
	"testing"
)

func TestAllCheckers(t *testing.T) {
	checkers := checkers.AllCheckers()
	if len(checkers) == 0 {
		t.Fatalf("no known checkers: empty list")
	}

	// todo: коммент, что специально не исопльзуются константы, чтобы в них не напутать
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
		"suite-thelper",
		"suite-dont-use-pkg",
	}
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestEnabledByDefaultCheckers(t *testing.T) {
	checkers := checkers.EnabledByDefaultCheckers()
	if len(checkers) == 0 {
		t.Fatalf("no enabled checkers: empty list")
	}

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
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestDisabledByDefaultCheckers(t *testing.T) {
	checkers := checkers.DisabledByDefaultCheckers()

	expected := []string{
		"suite-no-extra-assert-call",
		"suite-thelper",
	}
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestIsKnown(t *testing.T) {
	checkers := checkers.AllCheckers()

	for _, ch := range checkers {
		if !checkers.IsKnown(ch) {
			t.Fatalf("checker %v is unknown but mustn't be", ch)
		}
	}

	for _, ch := range []string{"", "lenlen", "bool-cmp"} {
		if checkers.IsKnown(ch) {
			t.Fatalf("checker %v is known but mustn't be", ch)
		}
	}
}

func TestGet(t *testing.T) {
	checkers := checkers.AllCheckers()
	checkerTypes := make(map[string]struct{}, len(checkers))

	for _, name := range checkers {
		ch, ok := checkers.Get(name)
		if !ok || ch == nil {
			t.Fatalf("lost checker: %v", name)
		}
		if ch.Name() != name {
			t.Fatalf("invalid checkers mapping: %v", name)
		}

		chType := strings.TrimLeft(fmt.Sprintf("%T", ch), "*")
		if _, ok := checkerTypes[chType]; ok {
			t.Fatalf("duplicated checker instance: %v (%v)", name, chType)
		}
		checkerTypes[chType] = struct{}{}
	}

	for _, name := range []string{"", "lenlen", "bool-cmp"} {
		ch, ok := checkers.Get(name)
		if ok || ch != nil {
			t.Fatalf("unexpected checker: %v", name)
		}
	}
}
