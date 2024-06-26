// Code generated by testifylint/internal/testgen. DO NOT EDIT.

package suitesubtestrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteSubtestRunCheckerSuite struct {
	suite.Suite
}

func TestSuiteSubtestRunCheckerSuite(t *testing.T) {
	suite.Run(t, new(SuiteSubtestRunCheckerSuite))
}

func (s *SuiteSubtestRunCheckerSuite) BeforeTest(suiteName, testName string) {
	s.Equal(1, 2)

	s.T().Run("init 1", func(t *testing.T) { // want "suite-subtest-run: use s\\.Run to run subtest"
		s.Require().Equal(1, 2)
		assert.Equal(t, 1, 2)
	})

	s.Run("init 2", func() {
		s.Require().Equal(2, 1)
	})
}

func (s *SuiteSubtestRunCheckerSuite) TestOne() {
	s.Equal(1, 2)

	s.T().Run("init 1", func(t *testing.T) { // want "suite-subtest-run: use s\\.Run to run subtest"
		s.Require().Equal(11, 22)

		s.Run("init 1.1", func() {
			s.T().Run("init 1.2", func(t *testing.T) { // want "suite-subtest-run: use s\\.Run to run subtest"
				s.Equal(111, 222)
				assert.Equal(t, 111, 222)
			})
		})
	})

	s.Run("init 2", func() {
		s.Equal(2, 1)
	})
}

func (suite *SuiteSubtestRunCheckerSuite) TestTwo() {
	suite.Equal(1, 2)

	suite.Run("init 1", func() {
		suite.Equal(1, 2)
	})

	suite.Run("init 2", func() {
		suite.Equal(2, 1)

		suite.T().Run("init 2.1", func(t *testing.T) { // want "suite-subtest-run: use suite\\.Run to run subtest"
			suite.Run("init 2.2", func() {
				suite.Require().Equal(222, 111)
			})
			assert.Equal(t, 22, 11)
		})
	})

	cases := []struct{ Name string }{}
	for _, tt := range cases {
		suite.T().Run(tt.Name, func(t *testing.T) { // want "suite-subtest-run: use suite\\.Run to run subtest"
			suite.Equal(1, 2)
		})
	}
}
