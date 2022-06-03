package basic

import "github.com/stretchr/testify/assert"

func (s *GameRoomSuite) TestExample() {
	assert.Equal(s.T(), 1, 2) // "use s.Equal(1, 2)"
}
