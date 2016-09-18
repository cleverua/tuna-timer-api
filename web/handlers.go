package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Handlers is a collection of net/http handlers to serve the API
type Handlers struct {
	env                   *utils.Environment
	status                map[string]string
	commandLookupFunction func(slackCommand data.SlackCommand) (commands.Command, error)
}

// NewHandlers constructs a Handlers collection
func NewHandlers(env *utils.Environment) *Handlers {
	return &Handlers{
		env: env,
		status: map[string]string{
			"env":     env.Name,
			"version": env.AppVersion,
		},
		commandLookupFunction: commands.Get,
	}
}

// Timer handles Slack /timer command
func (h *Handlers) Timer(w http.ResponseWriter, r *http.Request) {
	slackCommand := data.SlackCommand{
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

	w.Header().Set("Content-Type", "application/json")

	command, err := h.commandLookupFunction(slackCommand)
	if err != nil { //todo it is going to be a nicely formatted slack message sent back to user
		w.Write([]byte(fmt.Sprintf("Unknown command: %s!", slackCommand.Text)))
		return
	}

	session := h.mongoSession.Clone()
	defer session.Close()

	ctx := utils.PutMongoSessionInContext(r.Context(), session)
	result := command.Handle(ctx, slackCommand)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result.Body)
}

// Health handles a call for app health request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.env.CreatedAt)
	h.status["uptime"] = uptime.String() //is it good or not if I modify the map here?
	json.NewEncoder(w).Encode(h.status)
}
