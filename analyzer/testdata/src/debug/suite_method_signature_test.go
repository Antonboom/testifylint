package debug

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBrokenSuite(t *testing.T) {
	suite.Run(t, new(BrokenSuite))
}

type BrokenSuite struct {
	suite.Suite
}

func (s *BrokenSuite) SetupTest() {
	var _ suite.TearDownSubTest
}

func (s *BrokenSuite) SetT() { // *BrokenSuite does not implement suite.TestingSuite (wrong type for method SetT)
}

func (s *BrokenSuite) TestTypo(_ *testing.T) { // value.go:424: test panicked: reflect: Call with too few input arguments
	s.True(true)
}

/*
1) Flagged by govet (copylocks).
2) Related to https://github.com/go-critic/go-critic/issues/331
*/
func (s BrokenSuite) TestValueReceiver() { //  and revive
	s.True(true)
}
