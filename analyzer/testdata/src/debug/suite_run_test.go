package debug

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SomeSuite struct {
	suite.Suite
}

func TestSomeSuite(t *testing.T) {
	suite.Run(t, new(SomeSuite))
}

func (s *SomeSuite) BeforeTest(suiteName, testName string) {
	s.Equal(1, 2)

	s.T().Run("init 1", func(t *testing.T) {
		s.Require().Equal(1, 2)
	})

	s.Run("init 2", func() {
		s.Require().Equal(2, 1)
	})
}

func (s *SomeSuite) TestSomething() {
	s.T().Parallel()

	s.Run("sum", func() {
		dummy := 3 + 1
		s.Equal(4, dummy)
	})

	s.Run("mul", func() {
		dummy := 3 * 1
		s.Equal(3, dummy)
	})
}

func (s *SomeSuite) TestSomething_ThroughT() {
	s.T().Parallel()

	s.T().Run("sum", func(t *testing.T) {
		t.Parallel()

		dummy := 3 + 1
		s.Equal(4, dummy)
	})

	s.T().Run("mul", func(t *testing.T) {
		t.Parallel()

		dummy := 3 * 1
		s.Equal(3, dummy)
	})
}
