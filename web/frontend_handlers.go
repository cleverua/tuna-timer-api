package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
)

// Handlers is a collection of net/http handlers to serve the API
type FrontendHandlers struct {
	env                   *utils.Environment
	mongoSession          *mgo.Session
	status                map[string]string
}

// NewHandlers constructs a FrontendHandler collection
func NewFrontendHandlers(env *utils.Environment, mongoSession *mgo.Session) *FrontendHandlers {
	return &FrontendHandlers{
		env:          env,
		mongoSession: mongoSession,
		status: map[string]string{
			"env":     env.Name,
			"version": env.AppVersion,
		},
	}
}

func (h *FrontendHandlers) UserAuthentication(w http.ResponseWriter, r *http.Request) {
	response := JwtResponseBody{
		ResponseData: JwtToken{},
		ResponseBody: ResponseBody{
			ResponseErrors: map[string]string{},
			AppInfo: h.status,
		},
	}
	pid := r.PostFormValue("pid") // TODO: sanitize pid
	session := h.mongoSession.Clone()
	defer session.Close()

	passService := data.NewPassService(session)
	pass, err := passService.FindPassByToken(pid)

	w.Header().Set("Content-Type", "application/json")
	if err == nil && pass == nil {
		w.WriteHeader(http.StatusBadRequest)
		response.ResponseErrors["userMessage"] = "please login from slack application"
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ResponseErrors["developerMessage"] = err.Error()
	} else {
		jwtToken, jwtErr := NewUserToken(pass.TeamUserID, session)
		if jwtErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseErrors["developerMessage"] = jwtErr.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			response.ResponseData.Token = jwtToken
		}
	}

	json.NewEncoder(w).Encode(response)
}
