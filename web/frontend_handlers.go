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
	userLoginMessage = "please login from slack application"
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

func (h *FrontendHandlers) getUserFromJWT(token string, session *mgo.Session, status *ResponseStatus) (*models.TeamUser, bool) {
	jwtPayload := strings.Split(token, ".")[1]

	setErrors := func(err error) {
		status.Status = statusBadRequest
		status.UserMessage = userLoginMessage
		status.DeveloperMessage = err.Error()
	}

	decodedPayload, err := jwt.DecodeSegment(jwtPayload)
	if err != nil {
		setErrors(err)
		return nil, false
	}

	var userData struct { UserID string `json:"user_id"` }
	err = json.Unmarshal(decodedPayload, &userData)
	if err != nil {
		setErrors(err)
		return nil, false
	}

	userService := data.NewUserService(session)
	user, err := userService.FindByID(userData.UserID)
	if err != nil {
		setErrors(err)
		return nil, false
	}
	return user, true
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
			ResponseStatus: &ResponseStatus{},
			AppInfo: h.status,
		},
	}

	requestData := map[string]string{}
	err := h.jsonDecode(&requestData, r)
	if err != nil {
		response.ResponseStatus.Status = statusInternalServerError
		response.ResponseStatus.DeveloperMessage = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	pid := requestData["pid"] // TODO: sanitize pid

	session := h.mongoSession.Clone()
	defer session.Close()

	passService := data.NewPassService(session)
	pass, err := passService.FindPassByToken(pid)

	if err == nil && pass == nil {
		response.ResponseStatus.Status = statusBadRequest
		response.ResponseStatus.UserMessage = userLoginMessage
	} else if err != nil {
		response.ResponseStatus.Status = statusInternalServerError
		response.ResponseStatus.DeveloperMessage = err.Error()
	} else {
		jwtToken, jwtErr := NewUserToken(pass.TeamUserID, session)
		if jwtErr != nil {
			response.ResponseStatus.Status = statusInternalServerError
			response.ResponseStatus.DeveloperMessage = jwtErr.Error()
		} else {
			response.ResponseStatus.Status = statusOK
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
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: h.status,
		},
	}

	session := h.mongoSession.Clone()
	defer session.Close()

	user, ok := h.getUserFromJWT(r.Header.Get("Authorization"), session, response.ResponseStatus)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	timersService := data.NewTimerService(session)
	tasks, err := timersService.GetUserTasksByRange(startDate, endDate, user)
	if err != nil {
		response.ResponseStatus.Status = statusInternalServerError
		response.ResponseStatus.DeveloperMessage = err.Error()
	} else {
		response.ResponseData = tasks
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func(h *FrontendHandlers) UserProjectsData(w http.ResponseWriter, r *http.Request) {
	response := ProjectsResponseBody{
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: h.status,
		},
	}

	session := h.mongoSession.Clone()
	defer session.Close()

	user, ok := h.getUserFromJWT(r.Header.Get("Authorization"), session, response.ResponseStatus)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	teamsService := data.NewTeamService(session)
	team, err := teamsService.FindByID(user.TeamID)
	if err != nil {
		response.ResponseStatus.Status = statusInternalServerError
		response.ResponseStatus.DeveloperMessage = err.Error()
	}else {
		response.ResponseData = team.Projects
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
