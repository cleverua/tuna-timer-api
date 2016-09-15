package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
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

	cmd, _ := h.commandLookupFunction(slackCommand)
	cmd.Execute(h.env)

}

// DumpSlackCommand is a helper that logs an incoming POST from Slack. It is a temporary piece
func (h *Handlers) DumpSlackCommand(w http.ResponseWriter, r *http.Request) {

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

	log.Println("-----------------------------------------------")
	log.Printf("%+v\n", slackCommand)
	log.Println("-----------------------------------------------")
	dumpRequest(r)
	log.Println("-----------------------------------------------")

	text := map[string]string{
		"text": slackCommand.Text,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(text)
}

// Health handles a call for app health request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.env.CreatedAt)
	h.status["uptime"] = uptime.String() //is it good or not if I modify the map here?
	json.NewEncoder(w).Encode(h.status)
}

func dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println("Fail to dump the request!")
	}
	log.Println("Dumping the request:")
	log.Printf("%q", dump)
}
