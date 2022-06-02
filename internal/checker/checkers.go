package checker

import "fmt"

func IsKnown(name string) bool {
	_, ok := checkersByName[name]
	return ok
}

func KnownCheckers() []string {
	checkers := make([]string, 0, len(allCheckers))
	for _, ch := range allCheckers {
		checkers = append(checkers, ch.Name())
	}
	return checkers
}

func Get(name string) (Checker, bool) {
	ch, ok := checkersByName[name]
	return ch, ok
}

var allCheckers = []Checker{
	NewBoolCompare(),
	NewCompares(),
	NewEmpty(),
	NewError(),
	NewErrorIs(),
	NewExpectedActual(),
	NewFloatCompare(),
	NewLen(),
	NewRequireError(),
}

var checkersByName = make(map[string]Checker, len(allCheckers))

func init() {
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
