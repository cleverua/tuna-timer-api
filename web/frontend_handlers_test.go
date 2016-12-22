package web

import (
	"testing"
	"github.com/cleverua/tuna-timer-api/utils"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2/bson"
	"log"
	"gopkg.in/tylerb/is.v1"
	"gopkg.in/mgo.v2"
	"github.com/pavlo/gosuite"
	"net/http"
	"bytes"
	"net/http/httptest"
	"encoding/json"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

func TestFrontendHandlers(t *testing.T) {
	gosuite.Run(t, &FrontendHandlersTestSuite{Is: is.New(t)})
}

func (s *FrontendHandlersTestSuite) TestUserAuthentication(t *testing.T) {
	reqData := map[string]string{ "pid": "pass-for-jwt-generation" }
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(reqData)

	req, err := http.NewRequest("POST", "/api/v1/frontend/sessions", body)
	s.Nil(err)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := JwtResponseBody{ResponseData: JwtToken{}}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	verificationToken, err := NewUserToken(s.user.ID.Hex(), s.session)
	s.Nil(err)
	s.Equal(resp.ResponseErrors["status"], "200")
	s.Equal(resp.ResponseData.Token, verificationToken)
}

func (s *FrontendHandlersTestSuite) TestUserAuthenticationWithWrongPid(t *testing.T) {
	reqData := map[string]string{ "pid": "gIkuvaNzQIHg97ATvDxqgjtO" }
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(reqData)

	req, err := http.NewRequest("POST", "/api/v1/frontend/session", body)
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := JwtResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)
	s.Equal(resp.ResponseErrors["userMessage"], "please login from slack application")
	s.Equal(resp.ResponseErrors["status"], "400")
	s.Equal(resp.ResponseData.Token, "")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
}

func (s *FrontendHandlersTestSuite) TestGetUserFromJWT(t *testing.T) {
	//Should return user
	h := NewFrontendHandlers(s.env, s.session)
	user, err := h.getUserFromJWT(s.userJwt, s.session)

	s.Nil(err)
	s.Equal(user.ID, s.user.ID)
	s.Equal(user.ExternalUserID, s.user.ExternalUserID)
	s.Equal(user.ExternalUserName, s.user.ExternalUserName)

	//Should return error with corrupted payload
	jwtParts := strings.Split(s.userJwt, ".")
	jwtParts[1] += "==corruptedString"
	token := strings.Join(jwtParts, ".")

	user, err = h.getUserFromJWT(token, s.session)
	s.Err(err)
	s.Equal(err.Error(), "illegal base64 data at input byte 248")
	s.Nil(user)

	//Should return error with corrupted json data
	jwtParts = strings.Split(s.userJwt, ".")
	jwtParts[1] += "x"
	token = strings.Join(jwtParts, ".")

	user, err = h.getUserFromJWT(token, s.session)
	s.Err(err)
	s.Equal(err.Error(), "invalid character '\\f' after top-level value")
	s.Nil(user)

	//Should not return user with wrong ID, and return not found error
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":	 bson.NewObjectId(),
		"is_team_admin": s.user.SlackUserInfo.IsAdmin,
		"name":		 s.user.ExternalUserName,
		"image48":	 s.user.SlackUserInfo.Profile.Image48,
		"team_id":	 s.team.ID,
		"ext_team_id":	 s.team.ExternalTeamID,
		"ext_team_name": s.team.ExternalTeamName,

	})
	token, err = newToken.SignedString([]byte("TODO: Extract me in config/env"))

	user, err = h.getUserFromJWT(token, s.session)
	s.Err(err)
	s.Equal(err.Error(), "not found")
	s.Nil(user)
}

func (s *FrontendHandlersTestSuite) TestUserTimersData(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/timers?startDate=2016-12-20&endDate=2016-12-22", nil)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserTimersData)
	handler.ServeHTTP(recorder, req)

	resp := TasksResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseErrors["status"], "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(resp.ResponseData[0].ID, s.timer.ID)
}

func (s *FrontendHandlersTestSuite) TestUserTimersDataWithoutDateRange(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/timers", nil)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserTimersData)
	handler.ServeHTTP(recorder, req)

	resp := TasksResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseErrors["status"], "500")
	s.NotNil(resp.ResponseErrors["developerMessage"])
	s.Equal(resp.ResponseErrors["userMessage"], "")
	s.Len(resp.ResponseData, 0)
}

func (s *FrontendHandlersTestSuite) TestUserTimersDataWithNoExistingUser(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/timers?startDate=2016-12-20&endDate=2016-12-22", nil)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":	 bson.NewObjectId(),
		"is_team_admin": s.user.SlackUserInfo.IsAdmin,
		"name":		 s.user.ExternalUserName,
		"image48":	 s.user.SlackUserInfo.Profile.Image48,
		"team_id":	 s.team.ID,
		"ext_team_id":	 s.team.ExternalTeamID,
		"ext_team_name": s.team.ExternalTeamName,
	})
	token, err := newToken.SignedString([]byte("TODO: Extract me in config/env"))

	req.Header.Set("Authorization", "Bearer " + token)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserTimersData)
	handler.ServeHTTP(recorder, req)

	resp := TasksResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	s.Equal(resp.ResponseErrors["status"], "400")
	s.Equal(resp.ResponseErrors["developerMessage"], "not found")
	s.Equal(resp.ResponseErrors["userMessage"], "please login from slack application")
	s.Len(resp.ResponseData, 0)
}

type FrontendHandlersTestSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
	user    *models.TeamUser
	pass    *models.Pass
	team    *models.Team
	timer 	*models.Timer
	userJwt string
}

func (s *FrontendHandlersTestSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	s.session = session.Clone()
	e.MigrateDatabase(session)
	s.env = e
}

func (s *FrontendHandlersTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *FrontendHandlersTestSuite) SetUp() {
	//Clear Database
	utils.TruncateTables(s.session)

	//Seed Database
	passRepository := data.NewPassRepository(s.session)
	userRepository := data.NewUserRepository(s.session)
	teamRepository := data.NewTeamRepository(s.session)
	timerRepository := data.NewTimerRepository(s.session)

	var err error

	//Create team
	s.team, err = teamRepository.CreateTeam("ExtTeamID", "ExtTeamName")
	s.Nil(err)

	//Create user
	s.user = &models.TeamUser{
		TeamID:           s.team.ID.Hex(),
		ExternalUserID:   "ext-user-id",
		ExternalUserName: "user-name",
		SlackUserInfo:    &slack.User{
			IsAdmin: true,
		},
	}
	_, err = userRepository.Save(s.user)
	s.Nil(err)

	//Create pass
	s.pass = &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "pass-for-jwt-generation",
		TeamUserID:   s.user.ID.Hex(),
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	err = passRepository.Insert(s.pass)
	s.Nil(err)

	//Create timer
	s.timer, err = timerRepository.CreateTimer(
		&models.Timer{
			ID:         bson.NewObjectId(),
			TeamID:     s.team.ID.Hex(),
			ProjectID:  "project",
			TeamUserID: s.user.ID.Hex(),
			CreatedAt:  utils.PT("2016 Dec 21 00:00:00"),
			Minutes:    20,
	})
	s.Nil(err)

	//Generate user JWT
	s.userJwt, err = NewUserToken(s.user.ID.Hex(), s.session)
	s.Nil(err)
}

func (s *FrontendHandlersTestSuite) TearDown() {}
