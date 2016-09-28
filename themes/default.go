package themes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"time"
)

// DefaultSlackMessageTheme - todo
type DefaultSlackMessageTheme struct {
	themeConfig
	ctx context.Context
}

var defaultThemeConfig = themeConfig{
	MarkdownEnabledFor:     []string{"text", "pretext"},
	SummaryAttachmentColor: "#FFFFFF",
	FooterIcon:             "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png",

	StartCommandThumbURL: "/assets/themes/default/ic_current.png",
	StartCommandColor:    "FB6E04",

	StopCommandThumbURL: "/assets/themes/default/ic_completed.png",
	StopCommandColor:    "#2779DA",

	StatusCommandThumbURL: "/assets/themes/default/ic_status.png",
	StatusCommandColor:    "#959150",
}

func NewDefaultSlackMessageTheme(ctx context.Context) *DefaultSlackMessageTheme {
	return &DefaultSlackMessageTheme{
		themeConfig: defaultThemeConfig,
		ctx:         ctx,
	}
}

func (t *DefaultSlackMessageTheme) FormatStatusCommand(data *models.StatusCommandReport) string {

	tpl := slackThemeTemplate{
		Text:        fmt.Sprintf("Your status for %s", data.PeriodName),
		Attachments: []slack.Attachment{},
	}

	summaryAttachmentVisible := len(data.Tasks) > 0 || data.AlreadyStartedTimer != nil

	if len(data.Tasks) > 0 {
		statusAttachment := t.defaultAttachment()
		statusAttachment.ThumbURL = t.asset(t.StatusCommandThumbURL)
		statusAttachment.Color = t.StatusCommandColor

		statusAttachment.Footer = "<http://www.foo.com|Edit tasks in Application>"
		statusAttachment.FooterIcon = t.FooterIcon
		var buffer bytes.Buffer
		for _, task := range data.Tasks {
			buffer.WriteString(fmt.Sprintf("•  *%s*  %s\n", utils.FormatDuration(time.Duration(int64(task.Minutes)*int64(time.Minute))), task.Name))
		}
		statusAttachment.AuthorName = "Completed:"
		statusAttachment.Text = buffer.String()
		tpl.Attachments = append(tpl.Attachments, statusAttachment)
	}

	if data.AlreadyStartedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("%s", data.AlreadyStartedTimer.TaskName),
			t.asset(t.StartCommandThumbURL),
			data.AlreadyStartedTimer,
			data.AlreadyStartedTimerTotalForToday)

		sa.Color = t.StartCommandColor
		sa.AuthorName = "Current:"

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if summaryAttachmentVisible {
		tpl.Attachments = append(tpl.Attachments, t.summaryAttachment(data.PeriodName, data.UserTotalForPeriod))
	} else {
		tpl.Text = fmt.Sprintf("You have no tasks completed %s", data.PeriodName)
	}

	result, err := json.Marshal(tpl)
	if err != nil {
		// todo return { "text": err.String() }
	}

	return string(result)
}

func (t *DefaultSlackMessageTheme) FormatStopCommand(data *models.StopCommandReport) string {
	tpl := slackThemeTemplate{
		Text:        "Timer has been updated",
		Attachments: []slack.Attachment{},
	}

	if data.StoppedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Stopped for: %s", data.StoppedTimer.TaskName),
			t.asset(t.StopCommandThumbURL),
			data.StoppedTimer,
			data.StoppedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	tpl.Attachments = append(tpl.Attachments, t.summaryAttachment("today", data.UserTotalForToday))

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
			fmt.Sprintf("Stopped for: %s", data.StoppedTimer.TaskName),
			t.asset(t.StopCommandThumbURL),
			data.StoppedTimer,
			data.StoppedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.StartedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Started for: %s", data.StartedTimer.TaskName),
			t.asset(t.StartCommandThumbURL),
			data.StartedTimer,
			data.StartedTaskTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.AlreadyStartedTimer != nil {
		sa := t.attachmentForTimer(
			fmt.Sprintf("Already started for: %s", data.AlreadyStartedTimer.TaskName),
			t.asset(t.StartCommandThumbURL),
			data.AlreadyStartedTimer,
			data.AlreadyStartedTimerTotalForToday)

		tpl.Attachments = append(tpl.Attachments, sa)
	}

	tpl.Attachments = append(tpl.Attachments, t.summaryAttachment("today", data.UserTotalForToday))

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

	sa.Fields = []slack.AttachmentField{}

	thisRoundField := slack.AttachmentField{
		Title: "This Round",
		Value: utils.FormatDuration(time.Duration(int64(timer.Minutes) * int64(time.Minute))),
		Short: true,
	}

	todayField := slack.AttachmentField{
		Title: "Today",
		Value: utils.FormatDuration(time.Duration(int64(totalForToday) * int64(time.Minute))),
		Short: true,
	}

	sa.Fields = append(sa.Fields, thisRoundField, todayField)
	return sa
}

func (t *DefaultSlackMessageTheme) summaryAttachment(period string, minutes int) slack.Attachment {
	result := slack.Attachment{}
	result.Text = fmt.Sprintf("*Your total for %s is %s*",
		period,
		utils.FormatDuration(time.Duration(int64(minutes)*int64(time.Minute))))

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

func (t *DefaultSlackMessageTheme) asset(assetPath string) string {
	return utils.GetSelfBaseURLFromContext(t.ctx) + assetPath
}
