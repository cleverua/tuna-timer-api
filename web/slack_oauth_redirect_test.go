package web

import (
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/data"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *SlackOauthHandlersSuite) TestSlackOauth2RedirectEmptyCode(c *C) {
	req, err := http.NewRequest("GET", "/api/v1/slack/oauth2redirect", nil)
	if err != nil {
		c.Fatal(err)
	}

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SlackOauth2Redirect)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		c.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusBadRequest)
	}
}

func (s *SlackOauthHandlersSuite) TestSlackOauth2Redirect(c *C) {
	req, err := http.NewRequest("GET", "/api/v1/slack/oauth2redirect?code=2386021721.86286901378.a1666ad872&state=", nil)
	if err != nil {
		c.Fatal(err)
	}

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

	c.Assert(recorder.Code, Equals, http.StatusOK)

	// should create a Team and set its accessToken and stuff

	teamRepo := data.NewTeamRepository(s.session)

	team, err := teamRepo.FindByExternalID("team-id")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)
}

type TestSlackOAuth struct {
	result *slack.OAuthResponse
	err    error
}

func (o *TestSlackOAuth) GetOAuthResponse(clientID, clientSecret, code string) (*slack.OAuthResponse, error) {
	return o.result, o.err
}

func TestSlackOauthHandlers(t *testing.T) { TestingT(t) }

type SlackOauthHandlersSuite struct {
	env     *utils.Environment
	session *mgo.Session
}

var _ = Suite(&SlackOauthHandlersSuite{})

// Suite lifecycle and callbacks
func (s *SlackOauthHandlersSuite) SetUpSuite(c *C) {

	log.Println("SlackOauthHandlersSuite:SetUpSuite")

	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)
	s.env = e
	s.session = session.Clone()
}

func (s *SlackOauthHandlersSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *SlackOauthHandlersSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}
