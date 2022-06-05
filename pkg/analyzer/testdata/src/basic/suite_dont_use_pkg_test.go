package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SuiteDontUsePackageSuite struct {
	suite.Suite
}

func TestSuiteDontUsePackageSuite(t *testing.T) {
	suite.Run(t, new(SuiteDontUsePackageSuite))
}

func (s *SuiteDontUsePackageSuite) TestAll() {
	var b bool

	assert.True(s.T(), b)                         // want "suite-dont-use-pkg: use s\\.True"
	assert.True(s.T(), b, "msg")                  // want "suite-dont-use-pkg: use s\\.True"
	assert.True(s.T(), b, "msg with arg %d", 42)  // want "suite-dont-use-pkg: use s\\.True"
	assert.Truef(s.T(), b, "msg")                 // want "suite-dont-use-pkg: use s\\.Truef"
	assert.Truef(s.T(), b, "msg with arg %d", 42) // want "suite-dont-use-pkg: use s\\.Truef"

	require.False(s.T(), b)                         // want "suite-dont-use-pkg: use s\\.Require\\(\\)\\.False"
	require.False(s.T(), b, "msg")                  // want "suite-dont-use-pkg: use s\\.Require\\(\\)\\.False"
	require.False(s.T(), b, "msg with arg %d", 42)  // want "suite-dont-use-pkg: use s\\.Require\\(\\)\\.False"
	require.Falsef(s.T(), b, "msg")                 // want "suite-dont-use-pkg: use s\\.Require\\(\\)\\.Falsef"
	require.Falsef(s.T(), b, "msg with arg %d", 42) // want "suite-dont-use-pkg: use s\\.Require\\(\\)\\.Falsef"

	// Negative cases are covered by neighbour files.
}
