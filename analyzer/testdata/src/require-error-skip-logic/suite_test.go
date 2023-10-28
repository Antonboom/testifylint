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
	s.NoError(err)
	s.NoError(err)
}

func (s *RequireErrorSkipSuite) TearDownTest() {
	var err error
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), err)
}

func (s *RequireErrorSkipSuite) AfterTest(suiteName, testName string) {
	var err error
	assert.NoError(s.T(), err)
	s.NoError(err)
}

func (s *RequireErrorSkipSuite) TearDownSuite() {
	var err error
	s.NoError(err)
	assert.NoError(s.T(), err)
}

func (s *RequireErrorSkipSuite) HandleStats(suiteName string, stats *suite.SuiteInformation) {
	var err error
	s.NoError(err)
	s.Require().NoError(err)
	s.NoError(err)
}

func (s *RequireErrorSkipSuite) TestAll() {
	var err error

	if assert.NoError(s.T(), err) {
		assert.NoError(s.T(), err)
		assert.NoError(s.T(), err)
	} else {
		s.NoError(err)
		s.NoError(err)
	}

	s.Run("", func() {
		s.NoError(err) // want "require-error: for error assertions use require"

		if s.Error(err) {
			s.NoError(err)
			s.NoError(err)
		}

		if ok := s.NoError(err); ok {
			s.NoError(err)
			s.NoError(err)
		} else {
			s.T().Run("", func(t *testing.T) {
				s.NoError(err) // want "require-error: for error assertions use require"
				s.NoError(err)
			})
		}

		s.T().Run("", func(t *testing.T) {
			s.NoError(err) // want "require-error: for error assertions use require"
			s.NoError(err)
		})
	})

	for range []struct{}{} {
		s.Require().NoError(err)
		s.NoError(err) // want "require-error: for error assertions use require"
		s.NoError(err)
	}
}

func (s *RequireErrorSkipSuite) TestCleanup() {
	var err error

	s.T().Cleanup(func() {
		assert.NoError(s.T(), err)
		assert.NoError(s.T(), err)

		s.T().Cleanup(func() {
			assert.NoError(s.T(), err)
			assert.NoError(s.T(), err)
		})
	})
}

func (s *RequireErrorSkipSuite) TestGoroutine() {
	var err error

	go func() {
		assert.NoError(s.T(), err)
		assert.NoError(s.T(), err)

		s.Run("", func() {
			assert.NoError(s.T(), err) // want "require-error: for error assertions use require"
			assert.NoError(s.T(), err)
		})
	}()
}
