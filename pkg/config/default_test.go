package config_test

import (
	"testing"

	"github.com/Antonboom/testifylint/pkg/config"
)

func TestDefaultConfigIsValid(t *testing.T) {
	err := config.Validate(config.Default)
	if err != nil {
		t.Fatal(err.Error())
	}
}
