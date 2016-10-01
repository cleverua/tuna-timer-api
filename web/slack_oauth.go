package web

import (
	"github.com/nlopes/slack"
)

type SlackOAuth interface {
	GetOAuthResponse(clientID, clientSecret, code string) (*slack.OAuthResponse, error)
}

type SlackOAuthImpl struct {
}

func NewSlackOAuth() *SlackOAuthImpl {
	return &SlackOAuthImpl{}
}

func (*SlackOAuthImpl) GetOAuthResponse(clientID, clientSecret, code string) (*slack.OAuthResponse, error) {
	s := slack.New("???")
	s.SetDebug(true) // if I do not do this, then Slack's logger does not get initialization and fails in mics.go parseResponseBody. Todo: check this out!
	return slack.GetOAuthResponse(clientID, clientSecret, code, "", true)
}
