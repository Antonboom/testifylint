package checker

import (
	"fmt"
	"sort"
)

func IsKnown(name string) bool {
	_, ok := checkersByName[name]
	return ok
}

// AllCheckers returns all checkers names sorted by Checker.Priority.
func AllCheckers() []string {
	return checkersNames(allCheckers)
}

func EnabledByDefaultCheckers() []string {
	return checkersNames(enabledByDefault)
}

func DisabledByDefaultCheckers() []string {
	return checkersNames(disabledByDefault)
}

func checkersNames(in []Checker) []string {
	checkers := make([]string, 0, len(in))
	for _, ch := range in {
		checkers = append(checkers, ch.Name())
	}
	return checkers
}

func Get(name string) (Checker, bool) {
	ch, ok := checkersByName[name]
	return ch, ok
}

var (
	allCheckers = []Checker{
		NewBoolCompare(),
		NewCompares(),
		NewEmpty(),
		NewError(),
		NewErrorIs(),
		NewExpectedActual(),
		NewFloatCompare(),
		NewLen(),
		NewRequireError(),
		NewSuiteDontUsePkg(),
		NewSuiteNoExtraAssertCall(),
	}
	checkersByName = make(map[string]Checker, len(allCheckers))

	enabledByDefault  []Checker
	disabledByDefault []Checker
)

func init() {
	sort.SliceStable(allCheckers, func(i, j int) bool {
		return allCheckers[i].Priority() < allCheckers[j].Priority()
	})

	buildCheckersByName()
	buildEnabledByDefault()
}

func buildCheckersByName() {
	for _, ch := range allCheckers {
		name := ch.Name()
		if name == "" {
			panic(fmt.Sprintf("checker with empty name: %T", ch))
		}

		if _, ok := checkersByName[name]; ok {
			panic("duplicated checker: " + name)
		}
		checkersByName[name] = ch
	}
}

type disabler interface {
	DisabledByDefault() bool
}

func buildEnabledByDefault() {
	enabledByDefault = make([]Checker, 0, len(allCheckers))
	disabledByDefault = make([]Checker, 0, len(allCheckers))

	for _, ch := range allCheckers {
		if v, ok := ch.(disabler); ok && v.DisabledByDefault() {
			disabledByDefault = append(disabledByDefault, ch)
		} else {
			enabledByDefault = append(enabledByDefault, ch)
		}
	}
}
