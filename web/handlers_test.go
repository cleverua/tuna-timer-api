package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"

	. "gopkg.in/check.v1"
)

func TestHandlers(t *testing.T) { TestingT(t) }

type TestHandlersSuite struct {
	env *utils.Environment
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
	handlers := NewHandlers(s.env)
	handlers.commandLookupFunction = func(slackCommand data.SlackCommand) (commands.Command, error) {
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
	handler := http.HandlerFunc(handlers.Timer)

	handler.ServeHTTP(recorder, req)
	c.Assert(mockCmd.executed, Equals, true)
}

func (s *TestHandlersSuite) TestTimerCommandLookupFailure(c *C) {
	v := url.Values{}
	v.Set("text", "foobar")
	v.Set("token", "gIkuvaNzQIHg97ATvDxqgjtO")
	v.Set("team_id", "T0001")
	v.Set("team_domain", "example")
	v.Set("channel_id", "C2147483705")
	v.Set("channel_name", "test")
	v.Set("user_id", "U2147483697")
	v.Set("user_name", "Steve")
	v.Set("command", "/timer")
	v.Set("response_url", "https://hooks.slack.com/commands/1234/5678")

	req, err := http.NewRequest("POST", "/timer", bytes.NewBufferString(v.Encode()))
	if err != nil {
		c.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	h := NewHandlers(s.env, s.session)

	h.commandLookupFunction = func(slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error) {
		return nil, errors.New("Simulated failure")
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Timer)

	handler.ServeHTTP(recorder, req)

	c.Assert(recorder.Body.String(), Equals, "Unknown command: foobar!")
}

func (s *TestHandlersSuite) TestHealth(c *C) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		c.Fatal(err)
	}

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Health)
	handler.ServeHTTP(recorder, req)

	data := make(map[string]interface{})
	err = json.Unmarshal(recorder.Body.Bytes(), &data)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(data["env"].(string), Equals, utils.TestEnv)
	c.Assert(data["uptime"].(string), NotNil)
	c.Assert(data["version"].(string), Equals, s.env.AppVersion)
}

type mockCommand struct {
	executed bool
}

func (cmd *mockCommand) Execute(env *utils.Environment) *commands.CommandResult {
	cmd.executed = true
	return nil
}

// Suite lifecycle and callbacks
func (s *TestHandlersSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	e.MigrateDatabase()

	s.env = e
}

func (s *TestHandlersSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *TestHandlersSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.env)
}
