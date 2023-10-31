package requireerrorskiplogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRequireErrorSkipSuite(t *testing.T) {
	suite.Run(t, new(RequireErrorSkipSuite))
}

type RequireErrorSkipSuite struct {
	suite.Suite
}

func (s *RequireErrorSkipSuite) TearDownSubTest() {
	var err error
	s.Error(err)
	s.Error(err)
}

func (s *RequireErrorSkipSuite) TearDownTest() {
	var err error
	assert.Error(s.T(), err)
	assert.Error(s.T(), err)
}

func (s *RequireErrorSkipSuite) AfterTest(suiteName, testName string) {
	var err error
	assert.Error(s.T(), err)
	s.Error(err)
}

func (s *RequireErrorSkipSuite) TearDownSuite() {
	var err error
	s.Error(err)
	assert.Error(s.T(), err)
}

func (s *RequireErrorSkipSuite) HandleStats(suiteName string, stats *suite.SuiteInformation) {
	var err error
	s.Error(err)
	s.Require().Error(err)
	s.Error(err)
}

func (s *RequireErrorSkipSuite) TestAll() {
	var err error

	if assert.Error(s.T(), err) {
		assert.Error(s.T(), err)
		assert.Error(s.T(), err)
	} else {
		s.Error(err)
		s.Error(err)
	}

	s.Run("", func() {
		s.Error(err) // want "require-error: for error assertions use require"

		if s.Error(err) {
			s.Error(err)
			s.Error(err)
		}

		if ok := s.Error(err); ok {
			s.Error(err)
			s.Error(err)
		} else {
			s.T().Run("", func(t *testing.T) {
				s.Error(err) // want "require-error: for error assertions use require"
				s.Error(err)
			})
		}

		s.T().Run("", func(t *testing.T) {
			s.Error(err) // want "require-error: for error assertions use require"
			s.Error(err)
		})
	})

	for range []struct{}{} {
		s.Require().Error(err)
		s.Error(err) // want "require-error: for error assertions use require"
		s.Error(err)
	}
}

func (s *RequireErrorSkipSuite) TestCleanup() {
	var err error

	s.T().Cleanup(func() {
		assert.Error(s.T(), err)
		assert.Error(s.T(), err)

		s.T().Cleanup(func() {
			assert.Error(s.T(), err)
			assert.Error(s.T(), err)
		})
	})
}

func (s *RequireErrorSkipSuite) TestGoroutine() {
	var err error

	go func() {
		assert.Error(s.T(), err)
		assert.Error(s.T(), err)

		s.Run("", func() {
			assert.Error(s.T(), err) // want "require-error: for error assertions use require"
			assert.Error(s.T(), err)
		})
	}()
}

func (s *RequireErrorSkipSuite) TestNoErrorGroup() {
	var err error
	s.Require().NoError(err)
	s.NoError(err)
	s.NoErrorf(err, "boom")
	assert.NoError(s.T(), err)
	s.Require().NoErrorf(err, "boom")
}
