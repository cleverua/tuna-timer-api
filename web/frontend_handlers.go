package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/gorilla/context"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
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

// jsonDecode decodes json request body
func jsonDecode(data interface{}, r *http.Request, status *ResponseStatus) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		writeError(status, statusInternalServerError, err.Error(), "")
		return false
	}
	return true
}

// writeError write error messages to the response body
// code - status code, dm - developer message, um - user message
func writeError(rs *ResponseStatus, code, dm, um string ) {
	rs.DeveloperMessage = dm
	rs.UserMessage = um
	rs.Status = code
}

// encodeResponse encodes response body to JSON
func encodeResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *FrontendHandlers) Authenticate(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	resp := NewJWTResponse(h.status)
	defer encodeResponse(w, resp)

	requestData := map[string]string{}
	if ok := jsonDecode(&requestData, r, resp.ResponseStatus); !ok {
		return
	}

	pid := requestData["pid"] // TODO: sanitize pid

	passService := data.NewPassService(session)
	pass, err := passService.FindPassByToken(pid)

	if err == nil && pass == nil {
		writeError(resp.ResponseStatus, statusBadRequest, "", userLoginMessage)
	} else if err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
	} else {
		jwtToken, err := NewUserToken(pass.TeamUserID, session)
		if err != nil {
			writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
		}
		resp.ResponseData.Token = jwtToken
	}
}

func(h *FrontendHandlers) TimersData(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)
	resp := NewTimersResponse(h.status)
	defer encodeResponse(w, resp)

	//Get Query params (start and end date)
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")

	timersService := data.NewTimerService(session)
	tasks, err := timersService.GetUserTimersByRange(startDate, endDate, user)
	if err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
	}

	resp.ResponseData = tasks
}

func(h *FrontendHandlers) ProjectsData(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)

	resp := NewProjectsResponse(h.status)
	defer encodeResponse(w, resp)

	teamsService := data.NewTeamService(session)
	team, err := teamsService.FindByID(user.TeamID)
	if err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
	}

	resp.ResponseData = team.Projects
}

func(h *FrontendHandlers) CreateTimer(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)

	resp := NewTimersResponse(h.status)
	defer encodeResponse(w, resp)

	//Decode response data
	newTimer := &models.Timer{}
	if ok := jsonDecode(&newTimer, r, resp.ResponseStatus); !ok {
		return
	}

	timerService := data.NewTimerService(session)
	project := &models.Project{
		ID:                  bson.ObjectIdHex(newTimer.ProjectID),
		ExternalProjectName: newTimer.ProjectExternalName,
		ExternalProjectID:   newTimer.ProjectExternalID,
	}

	//Find and stop previous timer
	if activeTimer, _ := timerService.GetActiveTimer(user.TeamID, user.ID.Hex()); activeTimer != nil {
		if err := timerService.StopTimer(activeTimer); err != nil {
			writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
			return
		}
	}

	if _, err := timerService.StartTimer(user.TeamID, project, user, newTimer.TaskName); err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
		return
	}

	date := time.Now().Format("2006-1-2")
	timers, _ := timerService.GetUserTimersByRange(date, date, user)

	resp.ResponseData = timers
}

func(h *FrontendHandlers) UpdateTimer(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)

	resp := NewTimerResponse(h.status)
	defer encodeResponse(w, resp)

	// Decode response data
	newTimerData := &models.Timer{}
	if ok := jsonDecode(&newTimerData, r, resp.ResponseStatus); !ok {
		return
	}

	// STOP or update timer
	timerService := data.NewTimerService(session)
	var timer *models.Timer
	var err error
	if r.URL.Query().Get("stop_timer") != "" {
		timer, err = timerService.GetActiveTimer(user.TeamID, user.ID.Hex())
		if err != nil {
			writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
			return
		}

		if timer == nil {
			writeError(resp.ResponseStatus, statusBadRequest, mgo.ErrNotFound.Error(), "already stopped")
			return
		}
		timerService.StopTimer(timer)
	} else {
		timer, _ = timerService.FindByID(newTimerData.ID.Hex())

		if err = timerService.UpdateUserTimer(user, timer, newTimerData); err != nil {
			writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
			return
		}
	}

	resp.ResponseData = *timer
}

func (h *FrontendHandlers) DeleteTimer(w http.ResponseWriter, r *http.Request)  {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)

	resp := NewResponseBody(h.status)
	resp.ResponseStatus.UserMessage = "successfully deleted"
	defer encodeResponse(w, resp)

	timerID := mux.Vars(r)["id"]

	timerService := data.NewTimerService(session)
	timer, err := timerService.FindByID(timerID)
	if err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
		return
	}

	if err = timerService.DeleteUserTimer(user, timer); err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
	}
}

func (h *FrontendHandlers) MonthStatistic(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	user := context.Get(r, "user").(*models.TeamUser)

	resp := NewUserStatisticsResponse(h.status)
	defer encodeResponse(w, resp)

	date := r.URL.Query().Get("date")

	timerService := data.NewTimerService(session)
	monthStatistic, err := timerService.UserMonthStatistics(user, date)
	if err != nil {
		writeError(resp.ResponseStatus, statusInternalServerError, err.Error(), "")
		return
	}
	resp.ResponseData = monthStatistic
}
