package config

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/Antonboom/testifylint/internal/checkers"
)

var (
	_ flag.Value = (*KnownCheckersValue)(nil)
	_ flag.Value = (*RegexpValue)(nil)
)

// KnownCheckersValue implements comma separated list of testify checkers.
type KnownCheckersValue []string

func (kcv KnownCheckersValue) String() string {
	return strings.Join(kcv, ",")
}

func (kcv *KnownCheckersValue) Set(v string) error {
	chckrs := strings.Split(v, ",")
	for _, checkerName := range chckrs {
		if ok := checkers.IsKnown(checkerName); !ok {
			return fmt.Errorf("unknown checker %q", checkerName)
		}
	}

	*kcv = chckrs
	return nil
}

// RegexpValue is a special wrapper for support of flag.FlagSet over regexp.Regexp.
// Original regexp is available through Regexp.Regexp.
type RegexpValue struct {
	*regexp.Regexp
}

func (rv RegexpValue) String() string {
	if rv.Regexp == nil {
		return ""
	}
	return rv.Regexp.String()
}

func (rv *RegexpValue) Set(v string) error {
	compiled, err := regexp.Compile(v)
	if err != nil {
		return err
	}

	rv.Regexp = compiled
	return nil
}
