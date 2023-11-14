package debug

import "github.com/stretchr/testify/assert"

var _ assert.TestingT = tMock{}

type tMock struct{}

func (t tMock) Errorf(format string, args ...interface{}) {}
