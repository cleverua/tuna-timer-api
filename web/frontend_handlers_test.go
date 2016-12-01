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
)

func (s *TestHandlersSuite) TestUserAuthentication(t *testing.T) {
	_, user_err := utils.Create(&models.TeamUser{}, s.session)
	pass, pass_err := utils.Create(&models.Pass{}, s.session)
	s.Nil(user_err)
	s.Nil(pass_err)

	v := url.Values{}
	v.Set("pid", "pass-for-jwt-generation")
	req, err := http.NewRequest("POST", "/frontend/sessions", bytes.NewBufferString(v.Encode()))
	s.Nil(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := JwtResponseBody{ResponseData: JwtToken{}}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)

	verification_token, err := NewUserToken(pass.(models.Pass).TeamUserID, s.session)
	s.Nil(err)
	s.Equal(resp.ResponseErrors, make(map[string]string))
	s.Equal(resp.ResponseData.Token, verification_token)
}

func (s *TestHandlersSuite) TestUserAuthenticationWithWrongPid(t *testing.T) {
	v := url.Values{}
	v.Set("pid", "gIkuvaNzQIHg97ATvDxqgjtO")

	req, err := http.NewRequest("POST", "/api/v1/frontend/session", bytes.NewBufferString(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	s.Nil(err)

	h := NewHandlers(s.env, s.session)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UserAuthentication)
	handler.ServeHTTP(recorder, req)

	resp := JwtResponseBody{}
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	s.Nil(err)
	s.Equal(resp.ResponseErrors["userMessage"], "please login from slack application")
	s.Equal(resp.ResponseData.Token, "")
	s.Equal(resp.AppInfo["env"], utils.TestEnv)
	s.Equal(resp.AppInfo["version"], s.env.AppVersion)
}
