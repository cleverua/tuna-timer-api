package themes

import (
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
)

// SlackMessageTheme an interface each theme must to conform to
type SlackMessageTheme interface {
	FormatStartCommand(data *models.StartCommandReport) string
	FormatStopCommand(data *models.StopCommandReport) string
	FormatStatusCommand(data *models.StatusCommandReport) string
	FormatError(errorMessage string) string
}

type slackThemeTemplate struct {
	Text        string             `json:"text"`
	Attachments []slack.Attachment `json:"attachments"`
}

// themeConfig has a bunch of settings for formatting slack messages that get back in reply to /timer command
type themeConfig struct {
	MarkdownEnabledFor     []string
	SummaryAttachmentColor string
	FooterIcon             string
	StartCommandThumbURL   string
	StartCommandColor      string
	StopCommandThumbURL    string
	StopCommandColor       string
	StatusCommandThumbURL  string
	StatusCommandColor     string
	ErrorIcon              string
	ErrorColor             string
}
