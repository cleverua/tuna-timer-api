package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/gorilla/context"
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

func (h *FrontendHandlers) jsonDecode(data interface{}, r *http.Request, status *ResponseStatus) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		status.Status = statusInternalServerError
		status.DeveloperMessage = err.Error()
		return false
	}
	return true
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
	ok := h.jsonDecode(&requestData, r, response.ResponseStatus)
	if !ok {
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

	user := context.Get(r, "user").(*models.TeamUser)

	timersService := data.NewTimerService(session)
	tasks, err := timersService.GetUserTimersByRange(startDate, endDate, user)
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

	user := context.Get(r, "user").(*models.TeamUser)

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

func(h *FrontendHandlers) CreateUserTimer(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)
	response := NewTaskResponseBody(h.status)

	//Decode response data
	newTimerData := &models.Timer{}
	ok := h.jsonDecode(&newTimerData, r, response.ResponseStatus)
	if !ok {
		json.NewEncoder(w).Encode(response)
		return
	}

	timerService := data.NewTimerService(session)
	project := &models.Project{
		ID:                  newTimerData.ID,
		ExternalProjectName: newTimerData.ProjectExternalName,
		ExternalProjectID:   newTimerData.ProjectExternalID,
	}

	//Find and stop previous timer
	timerToStop, _ := timerService.GetActiveTimer(user.TeamID, user.ID.Hex())
	if timerToStop != nil {
		err := timerService.StopTimer(timerToStop)
		if err != nil {
			response.ResponseStatus.Status = statusInternalServerError
			response.ResponseStatus.DeveloperMessage = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	timer, err := timerService.StartTimer(user.TeamID, project, user, newTimerData.TaskName)
	if err != nil {
		response.ResponseStatus.Status = statusInternalServerError
		response.ResponseStatus.DeveloperMessage = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// TODO return all user timers for current day
	response.ResponseData = *timer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func(h *FrontendHandlers) UpdateUserTimer(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	response := NewTaskResponseBody(h.status)
	user := context.Get(r, "user").(*models.TeamUser)

	// Decode response data
	newTimerData := &models.Timer{}
	ok := h.jsonDecode(&newTimerData, r, response.ResponseStatus)
	if !ok {
		json.NewEncoder(w).Encode(response)
		return
	}

	//Get timer for update
	timerService := data.NewTimerService(session)
	var timer *models.Timer
	var err error

	if r.URL.Query().Get("stop_timer") != "" {
		timer, err = timerService.GetActiveTimer(user.TeamID, user.ID.Hex())
		if err != nil {
			response.ResponseStatus.Status = statusInternalServerError
			response.ResponseStatus.DeveloperMessage = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		} else {
			timerService.StopTimer(timer)
		}
	} else {
		timer, _ = timerService.FindByID(newTimerData.ID.Hex())
		err = timerService.UpdateUserTimer(user, timer, newTimerData)

		if err != nil {
			response.ResponseStatus.Status = statusInternalServerError
			response.ResponseStatus.DeveloperMessage = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response.ResponseData = *timer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
