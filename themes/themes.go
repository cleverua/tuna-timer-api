package themes

import (
	"github.com/nlopes/slack"
	"github.com/pavlo/slack-time/models"
)

// SlackMessageTheme an interface each theme must to conform to
type SlackMessageTheme interface {
	FormatStartCommand(data *models.StartCommandReport) string
	FormatStopCommand(data *models.StopCommandReport) string
	FormatStatusCommand(data *models.StatusCommandReport) string
}

type slackThemeTemplate struct {
	Text        string             `json:"text"`
	Attachments []slack.Attachment `json:"attachments"`
}

// SlackMessageTheme has a bunch of settings for formatting slack messages that get back in reply to /timer command
type themeConfig struct {
	// 	// Summary  Attachment is an attachment that goes below the most of messages, holds something like "Your total for today is 08:42h"

	MarkdownEnabledFor     []string
	SummaryAttachmentColor string
	FooterIcon             string
	StartCommandThumbURL   string
	StartCommandColor      string
	StopCommandThumbURL    string
	StopCommandColor       string
	StatusCommandThumbURL  string
	StatusCommandColor     string
	// 	// Start Message
	// 	StartMessageAttachmentColor    string
	// 	StartMessageAttachmentThumbURL string

	// 	// Stop Message
	// 	StopMessageAttachmentColor    string
	// 	StopMessageAttachmentThumbURL string

	// 	// Resume Message
	// 	ResumeMessageAttachmentColor    string
	// 	ResumeMessageAttachmentThumbURL string
}

// Fooo ...
// type Fooo interface {
// 	// startTimerAttachment(timer *models.Task) *slack.Attachment
// 	// summaryAttachment(timer *models.Task) *slack.Attachment
// }

// DefaultSlackMessageTheme ...
// type DefaultSlackMessageTheme struct {
// 	SlackMessageTheme
// }

// DefaultSlackMessageTheme represents default set of settings
// var _DefaultSlackMessageTheme = SlackMessageTheme{
// 	MarkdownEnabledFor:     []string{"text", "pretext"},
// 	SummaryAttachmentColor: "#FFFFFF",
// 	FooterIcon:             "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png",

// 	StartMessageAttachmentColor:    "#FB6E04",
// 	StartMessageAttachmentThumbURL: "http://icons.iconarchive.com/icons/graphicloads/100-flat/128/new-icon.png",
// }

// func defaultAttachment() *slack.Attachment {
// 	result := &slack.Attachment{}
// 	result.MarkdownIn = DefaultSlackMessageTheme.MarkdownEnabledFor
// 	result.FooterIcon = DefaultSlackMessageTheme.FooterIcon
// 	return result
// }

// func createField(title, value string, short bool) slack.AttachmentField {
// 	return slack.AttachmentField{
// 		Short: short,
// 		Title: title,
// 		Value: value,
// 	}
// }

//func summaryAttachment(text string) *slack.Attachment {
//	return &slack.Attachment{
//		AuthorName: text,
//		Color:      DefaultSlackMessageTheme.SummaryAttachmentColor,
//	}
//}
