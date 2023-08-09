package suite_thelper

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GameRoomSuite struct {
	suite.Suite

	rooms   map[int]*Room
	players map[int]*Player
}

func TestGameRoomSuite(t *testing.T) {
	suite.Run(t, &GameRoomSuite{
		rooms:   map[int]*Room{},
		players: map[int]*Player{},
	})
}

func (s *GameRoomSuite) TestJoinRoom() {
	p := s.newPlayer(100, "DoomGuy")
	r := s.newRoom(200)
	s.joinRoom(p.ID, r.ID)
	s.assertPlayerNickName(100, "DoomGuy")
}

func (s *GameRoomSuite) newPlayer(id int, nickname string) *Player {
	p := &Player{ID: id, Nickname: nickname}
	s.players[p.ID] = p
	return p
}

func (s *GameRoomSuite) newRoom(id int) *Room {
	r := &Room{ID: id}
	s.rooms[r.ID] = r
	return r
}

func (s *GameRoomSuite) joinRoom(playerID, roomID int) { // want "suite-thelper: suite helper method should start with s\\.T\\(\\)\\.Helper\\(\\)"
	room, ok := s.rooms[roomID]
	s.Require().True(ok)

	player, ok := s.players[playerID]
	s.Require().True(ok)

	err := room.AddPlayer(player)
	s.Require().NoError(err)
}

func (s *GameRoomSuite) assertPlayerNickName(playerID int, expectedNickname string) {
	s.T().Helper()

	player, ok := s.players[playerID]
	s.Require().True(ok)

	s.Equal(expectedNickname, player.Nickname)
}

type Player struct {
	ID       int
	Nickname string
}

type Room struct {
	ID      int
	Players []*Player
}

func (r Room) AddPlayer(player *Player) error {
	r.Players = append(r.Players, player)
	return nil
}
