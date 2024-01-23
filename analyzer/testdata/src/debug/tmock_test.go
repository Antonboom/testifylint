package debug

import "github.com/stretchr/testify/assert"

var _ assert.TestingT = tMock{}

type tMock struct{}

func (t tMock) Errorf(format string, args ...interface{}) {}

type assertion func(t assert.TestingT, expected, actual any, msgAndArgs ...any) bool
