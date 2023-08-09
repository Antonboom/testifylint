package checkers

import (
	"fmt"
	"sort"

	"slices"
)

var (
	// checkersRegistry stores checkers in priority order.
	checkersRegistry = []struct {
		Checker
		enabledByDefault bool
	}{
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

	checkersByName   = make(map[string]Checker, len(checkersRegistry))
	prioritiesByName = make(map[string]int, len(checkersRegistry))

	allCheckers      = make([]string, 0, len(checkersRegistry))
	enabledByDefault = make([]string, 0, len(checkersRegistry))
)

func init() {
	for i, checker := range checkersRegistry {
		name := checker.Name()
		if name == "" {
			panic(fmt.Sprintf("checker with empty name: %T", checker))
		}
		if duplicate, ok := checkersByName[name]; ok {
			panic(fmt.Sprintf("duplicated checker %q: %T and %T", name, duplicate, checker))
		}

		checkersByName[name] = checker.Checker
		prioritiesByName[name] = i

		allCheckers = append(allCheckers, checker.Name())
		if checker.enabledByDefault {
			enabledByDefault = append(enabledByDefault, checker.Name())
		}
	}

	sort.Strings(allCheckers)
	sort.Strings(enabledByDefault)
}

// All returns all checkers names sorted by name.
func All() []string {
	return slices.Clone(allCheckers)
}

// EnabledByDefault returns checkers enabled by default sorted by name.
func EnabledByDefault() []string {
	return slices.Clone(enabledByDefault)
}

// Get returns checker by its name.
func Get(name string) (Checker, bool) {
	ch, ok := checkersByName[name]
	return ch, ok
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
	return slices.Contains(enabledByDefault, name)
}

// SortByPriority mutates the input checkers by sorting them in order of priority.
func SortByPriority[T Checker](checkers []T) {
	sort.Slice(checkers, func(i, j int) bool {
		lhs, rhs := checkers[i], checkers[j]
		return prioritiesByName[lhs.Name()] < prioritiesByName[rhs.Name()]
	})
}
