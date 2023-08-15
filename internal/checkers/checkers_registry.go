package checkers

import (
	"fmt"
	"sort"
)

var (
	// checkersRegistry stores checkers in priority order.
	checkersRegistry = []struct {
		Checker
		enabledByDefault bool
	}{
		// Regular checkers.
		{Checker: NewBoolCompare(), enabledByDefault: true},
		{Checker: NewFloatCompare(), enabledByDefault: true},
		{Checker: NewEmpty(), enabledByDefault: true},
		{Checker: NewLen(), enabledByDefault: true},
		{Checker: NewCompares(), enabledByDefault: true},
		{Checker: NewError(), enabledByDefault: true},
		{Checker: NewErrorIs(), enabledByDefault: true},
		{Checker: NewRequireError(), enabledByDefault: true},
		{Checker: NewExpectedActual(), enabledByDefault: true},
		{Checker: NewSuiteNoExtraAssertCall(), enabledByDefault: false},
		{Checker: NewSuiteDontUsePkg(), enabledByDefault: true},
		// Advanced checkers.
		{Checker: NewSuiteTHelper(), enabledByDefault: false},
	}

	checkersByName = make(map[string]checkerMeta, len(checkersRegistry))
)

type checkerMeta struct {
	Checker
	enabledByDefault bool
	priority         int
}

func init() {
	for i, checker := range checkersRegistry {
		name := checker.Name()
		if name == "" {
			panic(fmt.Sprintf("checker with empty name: %T", checker))
		}
		if duplicate, ok := checkersByName[name]; ok {
			panic(fmt.Sprintf("duplicated checker %q: %T and %T", name, duplicate, checker))
		}

		checkersByName[name] = checkerMeta{
			Checker:          checker.Checker,
			enabledByDefault: checker.enabledByDefault,
			priority:         i,
		}
	}
}

// All returns all checkers names sorted by checker's priority.
func All() []string {
	result := make([]string, 0, len(checkersRegistry))
	for _, v := range checkersRegistry {
		result = append(result, v.Name())
	}
	return result
}

// EnabledByDefault returns checkers enabled by default sorted by checker's priority.
func EnabledByDefault() []string {
	result := make([]string, 0, len(checkersRegistry))
	for _, v := range checkersRegistry {
		if v.enabledByDefault {
			result = append(result, v.Name())
		}
	}
	return result
}

// Get returns checker by its name.
func Get(name string) (Checker, bool) {
	ch, ok := checkersByName[name]
	return ch.Checker, ok
}

// IsKnown checks if there is a checker with that name.
func IsKnown(name string) bool {
	_, ok := checkersByName[name]
	return ok
}

// IsEnabledByDefault returns true if a checker is enabled by default.
// Returns false if there is no such checker in the registry.
// For pre-validation use Get or IsKnown.
func IsEnabledByDefault(name string) bool {
	v, ok := checkersByName[name]
	return ok && v.enabledByDefault
}

// SortByPriority mutates the input checkers names by sorting them in checker priority order.
// Ignores unknown checkers. For pre-validation use Get or IsKnown.
func SortByPriority(checkers []string) {
	sort.Slice(checkers, func(i, j int) bool {
		lhs, rhs := checkers[i], checkers[j]
		return checkersByName[lhs].priority < checkersByName[rhs].priority
	})
}
