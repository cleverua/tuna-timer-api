package web

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

func TestHandlers(t *testing.T) { TestingT(t) }

type TestHandlersSuite struct {
	env     *utils.Environment
	session *mgo.Session
}

var _ = Suite(&TestHandlersSuite{})

func (s *TestHandlersSuite) TestTimer(c *C) {
	v := url.Values{}
	v.Set("token", "gIkuvaNzQIHg97ATvDxqgjtO")
	v.Set("team_id", "T0001")
	v.Set("team_domain", "example")
	v.Set("channel_id", "C2147483705")
	v.Set("channel_name", "test")
	v.Set("user_id", "U2147483697")
	v.Set("user_name", "Steve")
	v.Set("command", "/timer")
	v.Set("text", "start Convert the logotype to PNG")
	v.Set("response_url", "https://hooks.slack.com/commands/1234/5678")

	req, err := http.NewRequest("POST", "/timer", bytes.NewBufferString(v.Encode()))
	if err != nil {
		c.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	mockCmd := &mockCommand{executed: false}
	h := NewHandlers(s.env, s.session)

	h.commandLookupFunction = func(slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error) {
		c.Assert(slackCommand.ChannelID, Equals, "C2147483705")
		c.Assert(slackCommand.ChannelName, Equals, "test")
		c.Assert(slackCommand.Command, Equals, "/timer")
		c.Assert(slackCommand.ResponseURL, Equals, "https://hooks.slack.com/commands/1234/5678")
		c.Assert(slackCommand.TeamDomain, Equals, "example")
		c.Assert(slackCommand.TeamID, Equals, "T0001")
		c.Assert(slackCommand.Text, Equals, "start Convert the logotype to PNG")
		c.Assert(slackCommand.Token, Equals, "gIkuvaNzQIHg97ATvDxqgjtO")
		c.Assert(slackCommand.UserID, Equals, "U2147483697")
		c.Assert(slackCommand.UserName, Equals, "Steve")
		return mockCmd, nil
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Timer)

	handler.ServeHTTP(recorder, req)
	c.Assert(mockCmd.executed, Equals, true)
}

type mockCommand struct {
	executed bool
}

func (cmd *mockCommand) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *commands.SlackCustomCommandHandlerResult {
	cmd.executed = true
	return &commands.SlackCustomCommandHandlerResult{
		Body: []byte("OK"),
	}
}

func (cmd *mockCommand) GetName() string {
	return "mockCmd"
}

// Suite lifecycle and callbacks
func (s *TestHandlersSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)
	s.env = e
	s.session = session.Clone()
}

func (s *TestHandlersSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TestHandlersSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}
