// Code generated by testifylint/internal/testgen. DO NOT EDIT.

package suitebrokenparallel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteBrokenParallelCheckerSuite struct {
	suite.Suite
}

func TestSuiteBrokenParallelCheckerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SuiteBrokenParallelCheckerSuite))
}

func (s *SuiteBrokenParallelCheckerSuite) BeforeTest(_, _ string) {
	s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
}

func (s *SuiteBrokenParallelCheckerSuite) SetupTest() {
	s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
}

func (s *SuiteBrokenParallelCheckerSuite) TestAll() {
	s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"

	s.Run("", func() {
		s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
		s.Equal(1, 2)
	})

	s.T().Run("", func(t *testing.T) {
		t.Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
		assert.Equal(t, 1, 2)
	})

	s.Run("", func() {
		s.Equal(1, 2)

		s.Run("", func() {
			s.Equal(1, 2)

			s.Run("", func() {
				s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
				s.Equal(1, 2)
			})

			s.T().Run("", func(t *testing.T) {
				t.Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
				assert.Equal(t, 1, 2)
			})
		})
	})
}

func (s *SuiteBrokenParallelCheckerSuite) TestTable() {
	cases := []struct{ Name string }{}

	for _, tt := range cases {
		tt := tt
		s.T().Run(tt.Name, func(t *testing.T) {
			t.Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
			s.Equal(1, 2)
		})
	}

	for _, tt := range cases {
		tt := tt
		s.Run(tt.Name, func() {
			s.T().Parallel() // want "suite-broken-parallel: testify v1 does not support suite's parallel tests and subtests"
			s.Equal(1, 2)
		})
	}
}

func TestSimpleTable(t *testing.T) {
	t.Parallel()

	cases := []struct{ Name string }{}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, 1, 2)
		})
	}
}