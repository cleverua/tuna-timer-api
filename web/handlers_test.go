package web

import (
	"bytes"

	"encoding/json"
	"errors"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/cleverua/tuna-timer-api/commands"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"github.com/cleverua/tuna-timer-api/themes"
	"gopkg.in/tylerb/is.v1"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/pavlo/gosuite"
)


func TestHandlers(t *testing.T) {
	gosuite.Run(t, &TestHandlersSuite{Is: is.New(t)})
}

func (s *TestHandlersSuite) TestHandlersTimer(t *testing.T) {
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
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	mockCmd := &mockCommand{executed: false}
	h := NewHandlers(s.env, s.session)

	h.commandLookupFunction = func(ctx context.Context, slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error) {
		s.Equal(slackCommand.ChannelID, "C2147483705")
		s.Equal(slackCommand.ChannelName, "test")
		s.Equal(slackCommand.Command, "/timer")
		s.Equal(slackCommand.ResponseURL, "https://hooks.slack.com/commands/1234/5678")
		s.Equal(slackCommand.TeamDomain, "example")
		s.Equal(slackCommand.TeamID, "T0001")
		s.Equal(slackCommand.Text, "Convert the logotype to PNG")
		s.Equal(slackCommand.Token, "gIkuvaNzQIHg97ATvDxqgjtO")
		s.Equal(slackCommand.UserID, "U2147483697")
		s.Equal(slackCommand.UserName, "Steve")
		return mockCmd, nil
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Timer)

	handler.ServeHTTP(recorder, req)
	s.Equal(mockCmd.executed, true)
}

func (s *TestHandlersSuite) TestHandlersTimerCommandLookupFailure(t *testing.T) {
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
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	h := NewHandlers(s.env, s.session)

	h.commandLookupFunction = func(ctx context.Context, slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error) {
		return nil, errors.New("Simulated failure")
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Timer)

	handler.ServeHTTP(recorder, req)

	jsonResponse := recorder.Body.String()
	message := themes.SlackThemeTemplate{}
	json.Unmarshal([]byte(jsonResponse), &message)

	s.Equal(message.Attachments[0].Text, "Simulated failure")
}

func (s *TestHandlersSuite) TestHealth(t *testing.T) {

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Health)
	handler.ServeHTTP(recorder, req)

	data := make(map[string]interface{})
	err = json.Unmarshal(recorder.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	s.Equal(data["env"].(string), utils.TestEnv)
	s.NotNil(data["uptime"].(string))
	s.Equal(data["version"].(string), s.env.AppVersion)
}

type mockCommand struct {
	executed bool
}

func (cmd *mockCommand) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *commands.ResponseToSlack {
	cmd.executed = true
	return &commands.ResponseToSlack{
		Body: []byte("OK"),
	}
}

func (cmd *mockCommand) GetName() string {
	return "mockCmd"
}

type TestHandlersSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
}

func (s *TestHandlersSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)
	s.env = e
	s.session = session.Clone()
}

func (s *TestHandlersSuite) TearDownSuite() {
	s.session.Close()
}

func (s *TestHandlersSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *TestHandlersSuite) TearDown() {}
