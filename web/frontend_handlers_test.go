package web

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"net/url"
	"bytes"
	"github.com/cleverua/tuna-timer-api/utils"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (s *TestHandlersSuite) TestUserAuthentication(t *testing.T) {
	v := url.Values{}
	v.Set("pid", "pass-for-jwt-generation")

	passCollection := s.session.DB("").C(utils.MongoCollectionPasses)
	userCollection := s.session.DB("").C(utils.MongoCollectionTeamUsers)

	user := &models.TeamUser{
		ID:               bson.NewObjectId(),
		TeamID:           "team-id",
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
		SlackUserInfo: &slack.User{
			IsAdmin: true,
		},
	}

	pass := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "pass-for-jwt-generation",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
		TeamUserID:   user.ID.Hex(),
	}

	pass_err := passCollection.Insert(pass)
	user_err := userCollection.Insert(user)

	s.Nil(user_err)
	s.Nil(pass_err)

	req, err := http.NewRequest("POST", "/frontend/sessions", bytes.NewBufferString(v.Encode()))
	s.Nil(err)

	req.Header = map[string][]string{"Content-Type" : {"application/x-www-form-urlencoded"}}

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := ResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	jwt_token, err := NewToken(pass, s.session)
	s.Nil(err)

	s.Equal(resp.ResponseErrors, make(map[string]string))
	s.Equal(resp.ResponseData.Token, jwt_token)
}

func (s *TestHandlersSuite) TestUserAuthenticationWithWrongPid(t *testing.T) {
	v := url.Values{}
	v.Set("pid", "gIkuvaNzQIHg97ATvDxqgjtO")

	req, err := http.NewRequest("POST", "/api/v1/frontend/session", bytes.NewBufferString(v.Encode()))
	req.Header = map[string][]string{"Content-Type" : {"application/x-www-form-urlencoded"}}
	s.Nil(err)

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := ResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)
	s.Equal(resp.ResponseErrors["userMessage"], "please login from slack application")
	s.Equal(resp.ResponseData.Token, "")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
}
