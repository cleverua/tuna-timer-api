package web

import (
	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestSlackOauthHandlers(t *testing.T) {
	gosuite.Run(t, &SlackOauthHandlersSuite{Is: is.New(t)})
}

func (s *SlackOauthHandlersSuite) TestSlackOauth2RedirectEmptyCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/slack/oauth2redirect", nil)
	s.Nil(err)

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SlackOauth2Redirect)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusBadRequest)
	}
}

func (s *SlackOauthHandlersSuite) TestSlackOauth2Redirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/slack/oauth2redirect?code=2386021721.86286901378.a1666ad872&state=", nil)
	s.Nil(err)

	h := NewHandlers(s.env, s.session)
	h.slackOAuth = &TestSlackOAuth{
		result: &slack.OAuthResponse{
			AccessToken: "access-token",
			Scope:       "scope",
			TeamID:      "team-id",
			TeamName:    "team-name",
		},
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SlackOauth2Redirect)
	handler.ServeHTTP(recorder, req)

	s.Equal(recorder.Code, http.StatusOK)

	// should create a Team and set its accessToken and stuff
	teamRepo := data.NewTeamRepository(s.session)

	team, err := teamRepo.FindByExternalID("team-id")
	s.Nil(err)	
	s.NotNil(team)
}

type TestSlackOAuth struct {
	result *slack.OAuthResponse
	err    error
}

func (o *TestSlackOAuth) GetOAuthResponse(clientID, clientSecret, code string) (*slack.OAuthResponse, error) {
	return o.result, o.err
}

type SlackOauthHandlersSuite struct {
	env     *utils.Environment
	session *mgo.Session
	*is.Is
}

func (s *SlackOauthHandlersSuite) SetUpSuite() {

	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)
	s.env = e
	s.session = session.Clone()
}

func (s *SlackOauthHandlersSuite) TearDownSuite() {
	s.session.Close()
}

func (s *SlackOauthHandlersSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *SlackOauthHandlersSuite) TearDown() {}
