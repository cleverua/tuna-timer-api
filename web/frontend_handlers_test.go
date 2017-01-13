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
	"github.com/justinas/alice"
	"github.com/dgrijalva/jwt-go"
)

func TestFrontendHandlers(t *testing.T) {
	gosuite.Run(t, &FrontendHandlersTestSuite{Is: is.New(t)})
}

func (s *FrontendHandlersTestSuite) TestAuthenticate(t *testing.T) {
	reqData := map[string]string{ "pid": "pass-for-jwt-generation" }
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(reqData)

	req, err := http.NewRequest("POST", "/api/v1/frontend/sessions", body)
	s.Nil(err)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Authenticate)
	handler.ServeHTTP(recorder, req)

	resp := JWTResponseBody{ResponseData: JwtToken{}}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	verificationToken, err := NewUserToken(s.user.ID.Hex(), s.session)
	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.ResponseData.Token, verificationToken)
}

func (s *FrontendHandlersTestSuite) TestAuthenticateWithWrongPid(t *testing.T) {
	reqData := map[string]string{ "pid": "gIkuvaNzQIHg97ATvDxqgjtO" }
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(reqData)

	req, err := http.NewRequest("POST", "/api/v1/frontend/session", body)
	req.Header.Set("Content-Type", "application/json")
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Authenticate)
	handler.ServeHTTP(recorder, req)

	resp := JWTResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)
	s.Equal(resp.ResponseStatus.UserMessage, "please login from slack application")
	s.Equal(resp.ResponseStatus.Status, "400")
	s.Equal(resp.ResponseData.Token, "")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
}

func (s *FrontendHandlersTestSuite) TestTimersData(t *testing.T)  {
	date := time.Now().Format("2006-1-2")
	url := "/api/v1/frontend/timers?startDate=" + date + "&endDate=" + date

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.TimersData).ServeHTTP(recorder, req)

	resp := TimersResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(len(resp.ResponseData), 1)
	s.Equal(resp.ResponseData[0].ID, s.timer.ID)
}

func (s *FrontendHandlersTestSuite) TestTimersDataWithoutDateRange(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/timers", nil)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.TimersData).ServeHTTP(recorder, req)

	resp := TimersResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "500")
	s.NotNil(resp.ResponseStatus.DeveloperMessage)
	s.Equal(resp.ResponseStatus.UserMessage, "")
	s.Len(resp.ResponseData, 0)
}

func (s *FrontendHandlersTestSuite) TestTimersDataWithNoExistingUser(t *testing.T)  {
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
	s.middlewareChain.ThenFunc(h.TimersData).ServeHTTP(recorder, req)

	resp := TimersResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	s.Equal(resp.ResponseStatus.Status, "400")
	s.Equal(resp.ResponseStatus.DeveloperMessage, mgo.ErrNotFound.Error())
	s.Equal(resp.ResponseStatus.UserMessage, "please login from slack application")
	s.Len(resp.ResponseData, 0)
}

func (s *FrontendHandlersTestSuite) TestProjectsData(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/projects", nil)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	s.Nil(err)

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.ProjectsData).ServeHTTP(recorder, req)

	resp := ProjectsResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(resp.ResponseData[0].ID, s.team.Projects[0].ID)
	s.Equal(resp.ResponseData[0].ExternalProjectID, s.team.Projects[0].ExternalProjectID)
	s.Equal(resp.ResponseData[0].ExternalProjectName, s.team.Projects[0].ExternalProjectName)
}

func (s *FrontendHandlersTestSuite) TestProjectsDataWithNoExistedUser(t *testing.T)  {
	req, err := http.NewRequest("GET", "/api/v1/frontend/projects", nil)
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
	s.middlewareChain.ThenFunc(h.ProjectsData).ServeHTTP(recorder, req)
	resp := ProjectsResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "400")
	s.Equal(resp.ResponseStatus.DeveloperMessage, mgo.ErrNotFound.Error())
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
}

func (s *FrontendHandlersTestSuite) TestCreateTimer(t *testing.T)  {
	// It should create new timer, stop active timer and return current day Timers for User
	newTimer := models.Timer{
		TaskName:   "New task name",
		TeamUserID: s.user.ID.Hex(),
		TeamID:     s.team.ID.Hex(),
		ProjectID:  bson.NewObjectId().Hex(),
		ProjectExternalID:  "external-project-id",
		ProjectExternalName: "external-project-name",
		Minutes: 30,
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(newTimer)

	req, err := http.NewRequest("POST", "/api/v1/frontend/timers", body)
	s.Nil(err)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.CreateTimer).ServeHTTP(recorder, req)

	resp := TimersResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)

	// Check response timers data
	s.Len(resp.ResponseData, 2)

	// Check for first timer completed
	s.Equal(resp.ResponseData[0].ID, s.timer.ID)
	s.Equal(resp.ResponseData[0].Minutes, 20)
	s.NotNil(resp.ResponseData[0])

	// Check new timers data
	s.Equal(resp.ResponseData[1].TaskName, newTimer.TaskName)
	s.Equal(resp.ResponseData[1].TeamUserID, newTimer.TeamUserID)
	s.Equal(resp.ResponseData[1].ProjectID, newTimer.ProjectID)
	s.Equal(resp.ResponseData[1].ProjectExternalID, newTimer.ProjectExternalID)
	s.Equal(resp.ResponseData[1].ProjectExternalName, newTimer.ProjectExternalName)
	s.Equal(resp.ResponseData[1].Minutes, 0)
	s.Equal(resp.ResponseData[1].TeamID, s.user.TeamID)
}

func (s *FrontendHandlersTestSuite) TestUpdateTimer(t *testing.T)  {
	timersData := models.Timer{
		ID: s.timer.ID,
		TaskName:   "New task name",
		ProjectID:  bson.NewObjectId().Hex(),
		ProjectExternalID:  "external-project-id",
		ProjectExternalName: "external-project-name",
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(timersData)

	req, err := http.NewRequest("PUT", "/api/v1/frontend/timers", body)
	s.Nil(err)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.UpdateTimer).ServeHTTP(recorder, req)

	resp := TimerResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(resp.ResponseData.ID, s.timer.ID)
	s.Equal(resp.ResponseData.TaskName, timersData.TaskName)
	s.Equal(resp.ResponseData.TeamUserID, s.timer.TeamUserID)
	s.Equal(resp.ResponseData.ProjectID, timersData.ProjectID)
	s.Equal(resp.ResponseData.ProjectExternalID, timersData.ProjectExternalID)
	s.Equal(resp.ResponseData.ProjectExternalName, timersData.ProjectExternalName)
	s.Equal(resp.ResponseData.Minutes, 20)
	s.Equal(resp.ResponseData.TeamID, s.user.TeamID)
}

func (s *FrontendHandlersTestSuite) TestUpdateTimerWithNoExistingTimer(t *testing.T)  {
	timersData := models.Timer{
		ID: bson.NewObjectId(),
		TaskName:   "New task name",
		ProjectID:  bson.NewObjectId().Hex(),
		ProjectExternalID:  "external-project-id",
		ProjectExternalName: "external-project-name",
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(timersData)

	req, err := http.NewRequest("PUT", "/api/v1/frontend/timers", body)
	s.Nil(err)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.UpdateTimer).ServeHTTP(recorder, req)

	resp := TimerResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "500")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(resp.ResponseStatus.DeveloperMessage, mgo.ErrNotFound.Error())
}

func (s *FrontendHandlersTestSuite) TestUpdateTimerStopAction(t *testing.T)  {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(models.Timer{})

	req, err := http.NewRequest("PUT", "/api/v1/frontend/timers?stop_timer=true", body)
	s.Nil(err)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.UpdateTimer).ServeHTTP(recorder, req)

	resp := TimerResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "200")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
	s.Equal(resp.ResponseData.ID, s.timer.ID)
	s.Equal(resp.ResponseData.Minutes, 20)
	s.NotNil(resp.ResponseData.FinishedAt)
}

func (s *FrontendHandlersTestSuite) TestUpdateTimerStopAlreadyStoppedTimer(t *testing.T)  {
	timerService := data.NewTimerService(s.session)
	err := timerService.StopTimer(s.timer)
	s.Nil(err)

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(models.Timer{})

	req, err := http.NewRequest("PUT", "/api/v1/frontend/timers?stop_timer=true", body)
	s.Nil(err)
	req.Header.Set("Authorization", "Bearer " + s.userJwt)
	req.Header.Set("Content-Type", "application/json")

	h := NewFrontendHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	s.middlewareChain.ThenFunc(h.UpdateTimer).ServeHTTP(recorder, req)

	resp := TimerResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)

	s.Nil(err)
	s.Equal(resp.ResponseStatus.Status, "400")
	s.Equal(resp.ResponseStatus.DeveloperMessage, mgo.ErrNotFound.Error())
	s.Equal(resp.ResponseStatus.UserMessage, "already stopped")
}

// =================== TEST setup =================== //
type FrontendHandlersTestSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
	user    *models.TeamUser
	pass    *models.Pass
	team    *models.Team
	timer 	*models.Timer
	userJwt string
	secureCTX *SecureContext
	middlewareChain alice.Chain
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

	s.secureCTX = &SecureContext{
		Origin:  s.env.Config.UString("origin.url"),
		Session: s.session,
		Env: 	 s.env,
	}
	s.middlewareChain = alice.New(
		s.secureCTX.CorsMiddleware,
		JWTMiddleware,
		s.secureCTX.CurrentUserMiddleware)
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

	//Create team with project
	s.team, err = teamRepository.CreateTeam("ExtTeamID", "ExtTeamName")
	s.Nil(err)
	err = teamRepository.AddProject(s.team, "external-project-id", "external-project-name")
	s.Nil(err)
	s.team, _ = teamRepository.FindByID(s.team.ID.Hex())

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
			CreatedAt:  time.Now().Add(-20 * time.Minute),
			Minutes:    20,
	})
	s.Nil(err)

	//Generate user JWT
	s.userJwt, err = NewUserToken(s.user.ID.Hex(), s.session)
	s.Nil(err)
}

func (s *FrontendHandlersTestSuite) TearDown() {}
