package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteNoTestingTTestsGenerator struct{}

func (SuiteNoTestingTTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteNoTestingT()
}

func (g SuiteNoTestingTTestsGenerator) TemplateData() any {
	var (
		name   = g.Checker().Name()
		report = QuoteReport(name + ": suite method must not include a testing.T parameter")
	)

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(name),
		Report:      report,
	}
}

func (SuiteNoTestingTTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteNoTestingTTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteNoTestingTTestTmpl))
}

func (SuiteNoTestingTTestsGenerator) GoldenTemplate() Executor {
	golden := strings.ReplaceAll(suiteNoTestingTTestTmpl,
		") joinRoom(t *testing.T, roomID int)",
		") joinRoom(roomID int)",
	)
	golden = strings.ReplaceAll(golden,
		") createRoom(t *testing.T)",
		") createRoom()",
	)
	return template.Must(template.New("SuiteNoTestingTTestsGenerator.GoldenTemplate").Funcs(fm).Parse(golden))
}

const suiteNoTestingTTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChatSessionSuite struct {
	suite.Suite
}

func TestChatSessionSuite(t *testing.T) {
	suite.Run(t, &ChatSessionSuite{})
}

func (s *ChatSessionSuite) TestCreateRoom() {
	s.createRoom(&testing.T{})
}

func (s *ChatSessionSuite) TestJoinRoom() {
	s.joinRoom(&testing.T{}, 123)
}

func (suite *ChatSessionSuite) createRoom(t *testing.T) { // want {{ $.Report }}
	suite.Require().True(true)
}

func (suite *ChatSessionSuite) joinRoom(t *testing.T, roomID int) { // want {{ $.Report }}
	suite.Require().Equal(roomID, 123)
}
`
