package models

import "github.com/nlopes/slack"

// // SlackMessageTheme holds
// type SlackMessageTheme struct {
// }

// SlackMessageTheme has a bunch of settings for formatting slack messages
type SlackMessageTheme struct {
	// Summary  Attachment is an attachment that goes below the most of messages, holds something like "Your total for today is 08:42h"
	SummaryAttachmentColor string
	MarkdownEnabledFor     []string
	FooterIcon             string

	// Start Message
	StartMessageAttachmentColor    string
	StartMessageAttachmentThumbURL string

	// Stop Message
	StopMessageAttachmentColor    string
	StopMessageAttachmentThumbURL string

	// Resume Message
	ResumeMessageAttachmentColor    string
	ResumeMessageAttachmentThumbURL string
}

// Fooo ...
type Fooo interface {
	startTimerAttachment(timer *Task) *slack.Attachment
	summaryAttachment(timer *Task) *slack.Attachment
}

// DefaultSlackMessageTheme ...
type DefaultSlackMessageTheme struct {
	SlackMessageTheme
}

// DefaultSlackMessageTheme represents default set of settings
var _DefaultSlackMessageTheme = SlackMessageTheme{
	MarkdownEnabledFor:     []string{"text", "pretext"},
	SummaryAttachmentColor: "#FFFFFF",
	FooterIcon:             "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png",

	StartMessageAttachmentColor:    "#FB6E04",
	StartMessageAttachmentThumbURL: "http://icons.iconarchive.com/icons/graphicloads/100-flat/128/new-icon.png",
}

func (t *DefaultSlackMessageTheme) startTimerAttachment(timer *Task) *slack.Attachment {
	return nil
}

func (t *DefaultSlackMessageTheme) summaryAttachment(timer *Task) *slack.Attachment {
	return nil
}

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

// func todayTotalAttachment(text string) *slack.Attachment {
// 	return &slack.Attachment{
// 		AuthorName: text,
// 		Color:      DefaultSlackMessageTheme.SummaryAttachmentColor,
// 	}
// }
