package themes

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
	"time"
)

// DefaultSlackMessageTheme - todo
type DefaultSlackMessageTheme struct {
	themeConfig
}

var defaultThemeConfig = themeConfig{
	MarkdownEnabledFor:     []string{"text", "pretext"},
	SummaryAttachmentColor: "#FFFFFF",
	FooterIcon:             "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png",
	StartCommandThumbURL:   "http://icons.iconarchive.com/icons/graphicloads/100-flat/128/new-icon.png",
	StartCommandColor:      "FB6E04",
	StopCommandThumbURL:    "http://icons.iconarchive.com/icons/graphicloads/100-flat/128/pause-icon.png",
	StopCommandColor:       "#2779DA",
}

func NewDefaultSlackMessageTheme() *DefaultSlackMessageTheme {
	return &DefaultSlackMessageTheme{
		themeConfig: defaultThemeConfig,
	}
}

func (t *DefaultSlackMessageTheme) FormatStopCommand(data *models.StopCommandReport) string {
	tpl := slackThemeTemplate{
		Text:        "Timer has been updated",
		Attachments: []slack.Attachment{},
	}

	if data.StoppedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Stopped for: %s", data.StoppedTimer.TaskName),
			t.StartCommandThumbURL,
			data.StoppedTimer,
			data.StoppedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	tpl.Attachments = append(tpl.Attachments, t.summaryAttachment(data.UserTotalForToday))

	result, err := json.Marshal(tpl)
	if err != nil {
		// todo return { "text": err.String() }
	}

	return string(result)
}

func (t *DefaultSlackMessageTheme) FormatStartCommand(data *models.StartCommandReport) string {
	tpl := slackThemeTemplate{
		Text:        "Timer has been updated",
		Attachments: []slack.Attachment{},
	}

	if data.StoppedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Stopped for: %s", data.StartedTimer.TaskName),
			t.StartCommandThumbURL,
			data.StoppedTimer,
			data.StoppedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.StartedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Started for: %s", data.StartedTimer.TaskName),
			t.StartCommandThumbURL,
			data.StartedTimer,
			data.StartedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.AlreadyStartedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Already started for: %s", data.AlreadyStartedTimer.TaskName),
			t.StartCommandThumbURL,
			data.AlreadyStartedTimer,
			data.AlreadyStartedTimerTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	tpl.Attachments = append(tpl.Attachments, t.summaryAttachment(data.UserTotalForToday))

	result, err := json.Marshal(tpl)
	if err != nil {
		// todo return { "text": err.String() }
	}

	return string(result)
}

func (t *DefaultSlackMessageTheme) attachmentForTimer(text string, thumbURL string, timer *models.Timer, totalForToday int) slack.Attachment {
	sa := t.defaultAttachment()
	sa.Text = text
	sa.ThumbURL = thumbURL
	sa.Footer = fmt.Sprintf(
		"Task ID: %s > <http://www.google.com|Open in Application>", timer.TaskHash)
	sa.Color = t.StartCommandColor

	sa.Fields = []slack.AttachmentField{}

	thisRoundField := slack.AttachmentField{
		Title: "This Round",
		Value: utils.FormatDuration(time.Duration(timer.Minutes * int(time.Minute))),
		Short: true,
	}

	todayField := slack.AttachmentField{
		Title: "Today",
		Value: utils.FormatDuration(time.Duration(totalForToday * int(time.Minute))),
		Short: true,
	}

	sa.Fields = append(sa.Fields, thisRoundField, todayField)
	return sa
}

func (t *DefaultSlackMessageTheme) summaryAttachment(todayTotalMinutes int) slack.Attachment {
	result := slack.Attachment{}
	result.Text = fmt.Sprintf("*Your total for today is %s*",
		utils.FormatDuration(time.Duration(todayTotalMinutes*int(time.Minute))))

	result.Color = t.SummaryAttachmentColor
	result.MarkdownIn = t.MarkdownEnabledFor
	return result
}

func (t *DefaultSlackMessageTheme) defaultAttachment() slack.Attachment {
	result := slack.Attachment{}
	result.MarkdownIn = t.MarkdownEnabledFor
	result.FooterIcon = t.FooterIcon
	return result
}
