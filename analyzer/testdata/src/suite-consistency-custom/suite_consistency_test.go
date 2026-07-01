package suite_consistency_custom

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCustomSuite(t *testing.T) {
	suite.Run(t, new(CustomSuite))
}

func (s *CustomSuite) TestCustomName() { // want "suite-consistency: suite receiver name s does not match configured name suite" "suite-consistency: suite test method name TestCustomName does not match pattern \\^Test_\\[A-Z\\]\\[a-zA-Z0-9\\]\\*\\$"
	s.Equal(1, 1)
}

func (suite *CustomSuite) Test_CustomName() {}
