package debug

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// go test -race -run=TestSuiteWithParallelTests suite_broken_parallel_test.go
//
// DATARACE!
func TestSuiteWithParallelTests(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelTests))
}

type SuiteWithParallelTests struct {
	suite.Suite
}

func (s *SuiteWithParallelTests) TestOne() {
	s.T().Parallel()

	for i := 0; i < 100; i++ {
		s.GreaterOrEqual(i, 0)
	}
}

func (s *SuiteWithParallelTests) TestTwo() {
	s.T().Parallel()

	for i := 0; i < 100; i++ {
		s.GreaterOrEqual(i, 0)
	}
}

// go test -race -run=TestSuiteWithParallelSubTests suite_broken_parallel_test.go
//
// NO DATARACE, but s.T().Parallel() for subtests not working (18.407s):
//
//	TestOne: 0: I am started at 2024-06-06 09:18:56.72165
//	TestOne: 1: I am started at 2024-06-06 09:18:59.72348
//	TestOne: 2: I am started at 2024-06-06 09:19:02.72491
//	TestTwo: 1: I am started at 2024-06-06 09:19:08.72770
//	TestTwo: 2: I am started at 2024-06-06 09:19:11.72916
//	TestTwo: 0: I am started at 2024-06-06 09:19:05.72648
func TestSuiteWithParallelSubTests(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelSubTests))
}

type SuiteWithParallelSubTests struct {
	suite.Suite
}

func (s *SuiteWithParallelSubTests) TestOne() {
	for i := 0; i <= 2; i++ {
		i := i
		s.Run(fmt.Sprintf("%d", i), func() {
			s.T().Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestOne", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

func (s *SuiteWithParallelSubTests) TestTwo() {
	for i := 0; i <= 2; i++ {
		i := i
		s.Run(fmt.Sprintf("%d", i), func() {
			s.T().Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestTwo", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

// go test -race -run=TestSuiteWithParallelThroughSetupTest suite_broken_parallel_test.go
//
// DATARACE!
func TestSuiteWithParallelThroughSetupTest(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelThroughSetupTest))
}

type SuiteWithParallelThroughSetupTest struct {
	suite.Suite
}

func (s *SuiteWithParallelThroughSetupTest) SetupTest() {
	s.T().Parallel()
}

func (s *SuiteWithParallelThroughSetupTest) TestOne() {
	for i := 0; i <= 2; i++ {
		s.Run(fmt.Sprintf("%d", i), func() {
			s.T().Logf("%s: %d: I am started at %s", "TestOne", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

func (s *SuiteWithParallelThroughSetupTest) TestTwo() {
	for i := 0; i <= 2; i++ {
		s.Run(fmt.Sprintf("%d", i), func() {
			s.T().Logf("%s: %d: I am started at %s", "TestTwo", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

// go test -race -run=TestSuiteWithParallelRawTForRunSubTest suite_broken_parallel_test.go
//
// NO DATARACE, but:
//
//	panic: testing: t.Parallel called multiple times
func TestSuiteWithParallelRawTForRunSubTest(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelRawTForRunSubTest))
}

type SuiteWithParallelRawTForRunSubTest struct {
	suite.Suite
}

func (s *SuiteWithParallelRawTForRunSubTest) TestOne() {
	for i := 0; i <= 2; i++ {
		i := i
		s.T().Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s.T().Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestOne", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

func (s *SuiteWithParallelRawTForRunSubTest) TestTwo() {
	for i := 0; i <= 2; i++ {
		i := i
		s.T().Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s.T().Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestOne", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

// go test -race -run=TestSuiteWithParallelRawTForRunSubTestAndTParallel suite_broken_parallel_test.go
//
// NO DATARACE, and t.Parallel() working fine (6.279s):
//
//	AfterTest of TestOne: I am started at 2024-06-07 08:29:54.81907
//	TestOne: 0: I am started at 2024-06-07 08:29:54.81969
//	TestOne: 2: I am started at 2024-06-07 08:29:54.81982
//	TestOne: 1: I am started at 2024-06-07 08:29:54.81981
//	AfterTest of TestTwo: I am started at 2024-06-07 08:29:57.82127
//	TestTwo: 0: I am started at 2024-06-07 08:29:57.82130
//	TestTwo: 2: I am started at 2024-06-07 08:29:57.82132
//	TestTwo: 1: I am started at 2024-06-07 08:29:57.82132
//
// But AfterTest (and other hooks) don't work correctly.
func TestSuiteWithParallelRawTForRunSubTestAndTParallel(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelRawTForRunSubTestAndTParallel))
}

type SuiteWithParallelRawTForRunSubTestAndTParallel struct {
	suite.Suite
}

func (s *SuiteWithParallelRawTForRunSubTestAndTParallel) AfterTest(_, testName string) {
	s.T().Logf("AfterTest of %s: I am started at %s", testName, time.Now())
}

func (s *SuiteWithParallelRawTForRunSubTestAndTParallel) TestOne() {
	for i := 0; i <= 2; i++ {
		i := i
		s.T().Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestOne", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

func (s *SuiteWithParallelRawTForRunSubTestAndTParallel) TestTwo() {
	for i := 0; i <= 2; i++ {
		i := i
		s.T().Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			s.T().Logf("%s: %d: I am started at %s", "TestTwo", i, time.Now())
			s.GreaterOrEqual(i, 0)
			time.Sleep(3 * time.Second)
		})
	}
}

// go test -race -run=TestSuiteWithParallelInDifferentRunSubtest suite_broken_parallel_test.go
//
// DATARACE! But difficult to catch (try several times to see it).
func TestSuiteWithParallelInDifferentRunSubtest(t *testing.T) {
	suite.Run(t, new(SuiteWithParallelInDifferentRunSubtest))
}

type SuiteWithParallelInDifferentRunSubtest struct {
	suite.Suite
}

func (s *SuiteWithParallelInDifferentRunSubtest) TestOne() {
	s.T().Run("1", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 100; i++ {
			s.GreaterOrEqual(i, 0)
		}
	})

	s.Run("2", func() {
		s.T().Parallel()

		for i := 0; i < 100; i++ {
			s.GreaterOrEqual(i, 0)
		}
	})
}
