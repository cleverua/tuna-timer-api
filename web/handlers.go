package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	"context"

	"github.com/tuna-timer/tuna-timer-api/commands"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// Handlers is a collection of net/http handlers to serve the API
type Handlers struct {
	env                   *utils.Environment
	mongoSession          *mgo.Session
	status                map[string]string
	commandLookupFunction func(ctx context.Context, slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error)
}

// NewHandlers constructs a Handlers collection
func NewHandlers(env *utils.Environment, mongoSession *mgo.Session) *Handlers {
	return &Handlers{
		env:          env,
		mongoSession: mongoSession,
		status: map[string]string{
			"env":     env.Name,
			"version": env.AppVersion,
		},
		commandLookupFunction: commands.LookupHandler,
	}
}

// Timer handles Slack /timer command
func (h *Handlers) Timer(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	slackCommand := models.SlackCustomCommand{
		ChannelID:   r.PostFormValue("channel_id"),
		ChannelName: r.PostFormValue("channel_name"),
		Command:     r.PostFormValue("command"),
		ResponseURL: r.PostFormValue("response_url"),
		TeamDomain:  r.PostFormValue("team_domain"),
		TeamID:      r.PostFormValue("team_id"),
		Text:        r.PostFormValue("text"),
		Token:       r.PostFormValue("token"),
		UserID:      r.PostFormValue("user_id"),
		UserName:    r.PostFormValue("user_name"),
	}

	slackCommand = utils.NormalizeSlackCustomCommand(slackCommand)

	session := h.mongoSession.Clone()
	defer session.Close()

	ctx := utils.PutMongoSessionInContext(r.Context(), session)

	selfBaseURL := utils.GetSelfURLFromRequest(r)
	ctx = utils.PutSelfBaseURLInContext(ctx, selfBaseURL)

	command, err := h.commandLookupFunction(ctx, slackCommand)
	if err != nil { //todo it is going to be a nicely formatted slack message sent back to user
		w.Write([]byte(fmt.Sprintf("Unknown command: %s!", slackCommand.SubCommand)))
		return
	}

	result := command.Handle(ctx, slackCommand)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result.Body)

	//todo: rather defer it
	log.Printf("Timer command took %s", time.Since(now).String())
}

// Health handles a call for app health request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.env.CreatedAt)
	h.status["uptime"] = uptime.String() //is it good or not if I modify the map here?
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.status)
}

// ClearAllData - is supposed to be called by the QA team during early testing stage
func (h *Handlers) ClearAllData(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	utils.TruncateTables(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"success": true})
}
