package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteTHelperTestsGenerator struct{}

func (SuiteTHelperTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteTHelper()
}

func (g SuiteTHelperTestsGenerator) TemplateData() any {
	var (
		name    = g.Checker().Name()
		report1 = QuoteReport(name + ": suite helper method must start with suite.T().Helper()")
		report2 = QuoteReport(name + ": suite helper method must start with s.T().Helper()")
	)

	return struct {
		CheckerName      CheckerName
		Report1, Report2 string
	}{
		CheckerName: CheckerName(name),
		Report1:     report1,
		Report2:     report2,
	}
}

func (SuiteTHelperTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteTHelperTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteTHelperTestTmpl))
}

func (SuiteTHelperTestsGenerator) GoldenTemplate() Executor {
	golden := strings.ReplaceAll(suiteTHelperTestTmpl,
		"\troom, ok := suite.rooms[roomID]",
		"\tsuite.T().Helper()\n\n\troom, ok := suite.rooms[roomID]",
	)
	golden = strings.ReplaceAll(golden,
		"\troom, ok := s.rooms[roomID]",
		"\ts.T().Helper()\n\n\troom, ok := s.rooms[roomID]",
	)
	return template.Must(template.New("SuiteTHelperTestsGenerator.GoldenTemplate").Funcs(fm).Parse(golden))
}

const suiteTHelperTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

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

func (s *GameRoomSuite) SetupTest() {
	s.Require().Empty(s.rooms)
	s.Require().Empty(s.players)
}

func (s *GameRoomSuite) TearDownTest() {
	s.rooms = map[int]*Room{}
	s.players = map[int]*Player{}
}

func (s *GameRoomSuite) TestJoinRoom() {
	p := s.newPlayer(100, "DoomGuy")
	r := s.newRoom(200)
	s.joinRoom(p.ID, r.ID)
	s.assertPlayerNickName(100, "DoomGuy")
	s.AssertPlayersNumber(200, 1)
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

func (suite *GameRoomSuite) joinRoom(playerID, roomID int) { // want {{ $.Report1 }}
	room, ok := suite.rooms[roomID]
	suite.Require().True(ok)

	player, ok := suite.players[playerID]
	suite.Require().True(ok)

	err := room.AddPlayer(player)
	suite.Require().NoError(err)
}

func (s *GameRoomSuite) assertPlayerNickName(playerID int, expectedNickname string) {
	s.T().Helper()

	player, ok := s.players[playerID]
	s.Require().True(ok)

	s.Assert().Equal(playerID, player.ID)
	s.Equal(expectedNickname, player.Nickname)
}

func (s *GameRoomSuite) AssertPlayersNumber(roomID int, playersNum int) { // want {{ $.Report2 }}
	room, ok := s.rooms[roomID]
	s.Require().True(ok)
	s.Len(room.Players, playersNum)
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
`
