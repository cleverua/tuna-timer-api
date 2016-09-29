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

// DefaultSlackMessageTheme - the basic UI theme for messages that go back from us to Slack users
type DefaultSlackMessageTheme struct {
	themeConfig
	ctx context.Context
}

var defaultThemeConfig = themeConfig{
	MarkdownEnabledFor:     []string{"text", "pretext"},
	SummaryAttachmentColor: "#000000",
	FooterIcon:             "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png",
	StartCommandThumbURL:   "/assets/themes/default/ic_current.png",
	StartCommandColor:      "F5A623",
	StopCommandThumbURL:    "/assets/themes/default/ic_completed.png",
	StopCommandColor:       "#4A90E2",
	StatusCommandThumbURL:  "/assets/themes/default/ic_status.png",
	StatusCommandColor:     "#9B9B9B",
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
		statusAttachment.ThumbURL = t.asset(t.StopCommandThumbURL)
		statusAttachment.Color = t.StopCommandColor
		statusAttachment.AuthorName = "Completed:"

		var buffer bytes.Buffer

		for _, task := range data.Tasks {

			displayProjectLink := task.ProjectExternalID != data.Project.ExternalProjectID

			if data.AlreadyStartedTimer == nil || data.AlreadyStartedTimer.TaskHash != task.TaskHash {
				if displayProjectLink {
					buffer.WriteString(t.taskWithProject(task.Name, task.Minutes, task.ProjectExternalID, task.ProjectExternalName))
				} else {
					buffer.WriteString(t.task(task.Name, task.Minutes))
				}
			}
		}

		if buffer.Len() > 0 {
			statusAttachment.Text = buffer.String()
			statusAttachment.Footer = "<http://www.foo.com|Open in Application>"
			tpl.Attachments = append(tpl.Attachments, statusAttachment)
		}
	}

	if data.AlreadyStartedTimer != nil {
		sa := t.attachmentForCurrentTask(data.AlreadyStartedTimer, data.AlreadyStartedTimerTotalForToday)
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
		Attachments: []slack.Attachment{},
	}

	if data.StoppedTimer != nil {
		sa := t.attachmentForStoppedTask(data.StoppedTimer, data.StoppedTaskTotalForToday)
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
		Attachments: []slack.Attachment{},
	}

	if data.StoppedTimer != nil {
		sa := t.attachmentForStoppedTask(data.StoppedTimer, data.StoppedTaskTotalForToday)
		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.StartedTimer != nil {
		sa := t.attachmentForNewTask(data.StartedTimer, data.StartedTaskTotalForToday)
		tpl.Attachments = append(tpl.Attachments, sa)
	}

	if data.AlreadyStartedTimer != nil {
		sa := t.attachmentForNewTask(data.AlreadyStartedTimer, data.AlreadyStartedTimerTotalForToday)
		tpl.Attachments = append(tpl.Attachments, sa)
	}

	tpl.Attachments = append(tpl.Attachments, t.summaryAttachment("today", data.UserTotalForToday))

	result, err := json.Marshal(tpl)
	if err != nil {
		// todo return { "text": err.String() }
	}

	return string(result)
}

func (t *DefaultSlackMessageTheme) attachmentForNewTask(timer *models.Timer, taskTotalForToday int) slack.Attachment {
	sa := t.defaultAttachment()
	sa.Text = t.task(timer.TaskName, taskTotalForToday) //fmt.Sprintf("•  *%s*  %s\n", utils.FormatDuration(time.Duration(int64(taskTotalForToday)*int64(time.Minute))), timer.TaskName)
	sa.ThumbURL = t.asset(t.StartCommandThumbURL)
	sa.Color = t.StartCommandColor
	sa.AuthorName = "Started:"

	sa.Footer = fmt.Sprintf(
		"Project: %s > Task: %s > <http://www.google.com|Edit in Application>", t.channelLinkForTimer(timer), timer.TaskHash)

	return sa
}

func (t *DefaultSlackMessageTheme) attachmentForCurrentTask(timer *models.Timer, totalForToday int) slack.Attachment {
	sa := t.defaultAttachment()
	sa.Text = t.task(timer.TaskName, totalForToday) //fmt.Sprintf("•  *%s*  %s\n", utils.FormatDuration(time.Duration(int64(totalForToday)*int64(time.Minute))), timer.TaskName)
	sa.ThumbURL = t.asset(t.StartCommandThumbURL)
	sa.Color = t.StartCommandColor
	sa.AuthorName = "Current:"

	sa.Footer = fmt.Sprintf(
		"Project: %s > Task: %s > <http://www.google.com|Open in Application>", t.channelLinkForTimer(timer), timer.TaskHash)

	sa.Fields = []slack.AttachmentField{}
	return sa
}

func (t *DefaultSlackMessageTheme) attachmentForStoppedTask(timer *models.Timer, totalForToday int) slack.Attachment {
	sa := t.defaultAttachment()
	sa.AuthorName = "Completed:"

	sa.Text = t.task(timer.TaskName, totalForToday) // fmt.Sprintf("•  *%s*  %s\n", utils.FormatDuration(time.Duration(int64(totalForToday)*int64(time.Minute))), timer.TaskName)
	sa.ThumbURL = t.asset(t.StopCommandThumbURL)
	sa.Color = t.StopCommandColor

	sa.Footer = fmt.Sprintf(
		"Project: %s > Task: %s > <http://www.google.com|Open in Application>", t.channelLinkForTimer(timer), timer.TaskHash)

	sa.Fields = []slack.AttachmentField{}
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
	return result
}

func (t *DefaultSlackMessageTheme) asset(assetPath string) string {
	return utils.GetSelfBaseURLFromContext(t.ctx) + assetPath
}

func (t *DefaultSlackMessageTheme) channelLinkForTimer(timer *models.Timer) string {
	return t.channelLink(timer.ProjectExternalID, timer.ProjectExternalName)
}

func (t *DefaultSlackMessageTheme) channelLink(channelID, channelName string) string {
	return fmt.Sprintf("<#%s|%s>", channelID, channelName)
}

func (t *DefaultSlackMessageTheme) task(text string, minutes int) string {
	return fmt.Sprintf("•  *%s*  %s\n", utils.FormatDuration(time.Duration(int64(minutes)*int64(time.Minute))), text)
}

func (t *DefaultSlackMessageTheme) taskWithProject(text string, minutes int, projectID, projectName string) string {
	return fmt.Sprintf("•  *%s  *%s  %s\n",
		utils.FormatDuration(time.Duration(int64(minutes)*int64(time.Minute))),
		t.channelLink(projectID, projectName),
		text)
}
