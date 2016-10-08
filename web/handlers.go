package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	"context"

	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/commands"
	"github.com/tuna-timer/tuna-timer-api/data"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/themes"
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
	slackOAuth            SlackOAuth
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
		slackOAuth:            NewSlackOAuth(),
	}
}

// Timer handles Slack /timer command
func (h *Handlers) Timer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
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

	theme := themes.NewDefaultSlackMessageTheme(ctx)
	ctx = utils.PutThemeInContext(ctx, theme)

	command, err := h.commandLookupFunction(ctx, slackCommand)
	if err != nil {
		w.Write([]byte(theme.FormatError(err.Error())))
		return
	}

	result := command.Handle(ctx, slackCommand)
	w.Write(result.Body)

	//todo: rather defer it
	log.Printf("Timer command took %s", time.Since(now).String())
}

// SlackOauth2Redirect handles the OAuth2 redirect from Slack and exchanges the `code` with `accessToken`
// https://api.slack.com/methods/oauth.access
func (h *Handlers) SlackOauth2Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	clientID := h.env.Config.UString("slack.client_id")
	clientSecret := h.env.Config.UString("slack.client_secret")

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("`code` parameter is either missed or blank!"))
		return
	}

	oauthResponse, err := h.slackOAuth.GetOAuthResponse(clientID, clientSecret, code)
	if err != nil {
		msg := fmt.Sprintf("Got a failure during getting an access token from Slack: %s", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	teamService := data.NewTeamService(h.mongoSession)

	err = teamService.CreateOrUpdateWithSlackOAuthResponse(oauthResponse)
	if err != nil {
		msg := fmt.Sprintf("Got a failure during creating a team: %s", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Health handles a call for app health request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.env.CreatedAt)
	h.status["uptime"] = uptime.String() //is it good or not if I modify the map here?
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.status)
}

func (h *Handlers) SendSampleMessageFromBot(w http.ResponseWriter, r *http.Request) {

	teamRepo := data.NewTeamRepository(h.mongoSession)
	team, _ := teamRepo.FindByExternalID("T02BC0MM7")

	accessToken := team.SlackOAuth.Bot.BotAccessToken
	slackAPI := slack.New(accessToken)

	slackAPI.PostMessage("U02BC0MM9", "You're about stopping a timer...", slack.PostMessageParameters{
		AsUser: true,
		Attachments: []slack.Attachment{
			{
				Text:       "Would you like to stop the timer?",
				AuthorName: "Pavlo",
				Actions: []slack.AttachmentAction{
					{
						Text:  "Yes, I'd like to stop it",
						Name:  "yes",
						Type:  "button",
						Style: "danger",
						Confirm: []slack.ConfirmationField{
							{
								Text:        "Are you sure?",
								DismissText: "Cancel",
								OkText:      "Yes!",
								Title:       "Are you sure you want to stop the timer?",
							},
						},
					},
					{
						Text:  "I am not sure yet",
						Name:  "not sure",
						Type:  "button",
						Style: "default",
					},
					{
						Text:  "No, let's keep it!",
						Name:  "no",
						Type:  "button",
						Style: "primary",
					},
				},
			},
		},
	})
}

// ClearAllData - is supposed to be called by the QA team during early testing stage
func (h *Handlers) ClearAllData(w http.ResponseWriter, r *http.Request) {
	session := h.mongoSession.Clone()
	defer session.Close()
	utils.TruncateTables(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"success": true})
}
