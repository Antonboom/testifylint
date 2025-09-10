package main

func (s *UnitTestSuite) TestX() {
	s.Assert().Equal(42, 43) // want "suite-extra-assert-call: need to simplify the assertion to s\\.Equal"
	s.Equal(42, 43)
}
