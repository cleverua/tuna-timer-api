package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/cleverua/tuna-timer-api/models"
)

const (
	statusOK = "200"
	statusBadRequest = "400"
	statusInternalServerError = "500"
	userMessage = "please login from slack application"
)

// Handlers is a collection of net/http handlers to serve the API
type FrontendHandlers struct {
	env                   *utils.Environment
	mongoSession          *mgo.Session
	status                map[string]string
}

func (h *FrontendHandlers) jsonDecode(data *map[string]string, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

func (h *FrontendHandlers) getUserFromJWT(token string, session *mgo.Session) (*models.TeamUser, error) {
	jwtPayload := strings.Split(token, ".")[1]

	decodedPayload, err := jwt.DecodeSegment(jwtPayload)
	if err != nil {
		return nil, err
	}

	var userData struct { UserID string `json:"user_id"` }
	err = json.Unmarshal(decodedPayload, &userData)
	if err != nil {
		return nil, err
	}

	userService := data.NewUserService(session)
	return userService.FindByID(userData.UserID)
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

	requestData := map[string]string{}
	err := h.jsonDecode(&requestData, r)
	if err != nil {
		response.ResponseErrors["status"] = statusInternalServerError
		response.ResponseErrors["developerMessage"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	pid := requestData["pid"] // TODO: sanitize pid

	session := h.mongoSession.Clone()
	defer session.Close()

	passService := data.NewPassService(session)
	pass, err := passService.FindPassByToken(pid)

	if err == nil && pass == nil {
		response.ResponseErrors["status"] = statusBadRequest
		response.ResponseErrors["userMessage"] = userMessage
	} else if err != nil {
		response.ResponseErrors["status"] = statusInternalServerError
		response.ResponseErrors["developerMessage"] = err.Error()
	} else {
		jwtToken, jwtErr := NewUserToken(pass.TeamUserID, session)
		if jwtErr != nil {
			response.ResponseErrors["status"] = statusInternalServerError
			response.ResponseErrors["developerMessage"] = jwtErr.Error()
		} else {
			response.ResponseErrors["status"] = statusOK
			response.ResponseData.Token = jwtToken
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func(h *FrontendHandlers) UserTimersData(w http.ResponseWriter, r *http.Request) {
	//Get Query params (start and end date)
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")

	response := TasksResponseBody{
		ResponseData: nil,
		ResponseBody: ResponseBody{
			ResponseErrors: map[string]string{
				"status": statusOK,
			},
			AppInfo: h.status,
		},
	}

	session := h.mongoSession.Clone()
	defer session.Close()

	user, err := h.getUserFromJWT(r.Header.Get("Authorization"), session)
	if err != nil {
		response.ResponseErrors["status"] = statusBadRequest
		response.ResponseErrors["userMessage"] = userMessage
		response.ResponseErrors["developerMessage"] = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	timersService := data.NewTimerService(session)
	response.ResponseData, err = timersService.GetUserTasksByRange(startDate, endDate, user)
	if err != nil {
		response.ResponseErrors["status"] = statusInternalServerError
		response.ResponseErrors["developerMessage"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
