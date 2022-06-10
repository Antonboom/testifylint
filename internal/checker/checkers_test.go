package checker_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Antonboom/testifylint/internal/checker"
)

func TestAllCheckers(t *testing.T) {
	checkers := checker.AllCheckers()
	if len(checkers) == 0 {
		t.Fatalf("no known checkers: empty list")
	}

	expected := []string{
		"bool-compare",
		"suite-dont-use-pkg",
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
	}
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestEnabledByDefaultCheckers(t *testing.T) {
	checkers := checker.EnabledByDefaultCheckers()
	if len(checkers) == 0 {
		t.Fatalf("no enabled checkers: empty list")
	}

	expected := []string{
		"bool-compare",
		"suite-dont-use-pkg",
		"float-compare",
		"empty",
		"len",
		"compares",
		"error",
		"error-is",
		"require-error",
		"expected-actual",
	}
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestDisabledByDefaultCheckers(t *testing.T) {
	checkers := checker.DisabledByDefaultCheckers()

	expected := []string{
		"suite-no-extra-assert-call",
		"suite-thelper",
	}
	if !reflect.DeepEqual(checkers, expected) {
		t.Fatalf("unexpected list: %#v", checkers)
	}
}

func TestIsKnown(t *testing.T) {
	checkers := checker.AllCheckers()

	for _, ch := range checkers {
		if !checker.IsKnown(ch) {
			t.Fatalf("checker %v is unknown but mustn't be", ch)
		}
	}

	for _, ch := range []string{"", "lenlen", "bool-cmp"} {
		if checker.IsKnown(ch) {
			t.Fatalf("checker %v is known but mustn't be", ch)
		}
	}
}

func TestGet(t *testing.T) {
	checkers := checker.AllCheckers()
	checkerTypes := make(map[string]struct{}, len(checkers))

	for _, name := range checkers {
		ch, ok := checker.Get(name)
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
		ch, ok := checker.Get(name)
		if ok || ch != nil {
			t.Fatalf("unexpected checker: %v", name)
		}
	}
}
