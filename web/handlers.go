package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
)

// Handlers is a collection of net/http handlers to serve the API
type Handlers struct {
	env                   *utils.Environment
	mongoSession          *mgo.Session
	status                map[string]string
	commandLookupFunction func(slackCommand models.SlackCustomCommand) (commands.SlackCustomCommandHandler, error)
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

	command, err := h.commandLookupFunction(slackCommand)
	if err != nil {
		log.Fatal("Failed to look up handler!")
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

func dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println("Fail to dump the request!")
	}
	log.Println("Dumping the request:")
	log.Printf("%q", dump)
}
