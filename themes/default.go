package themes

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pavlo/slack-time/models"
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
}

func NewDefaultSlackMessageTheme() *DefaultSlackMessageTheme {
	return &DefaultSlackMessageTheme{
		themeConfig: defaultThemeConfig,
	}
}

func (t *DefaultSlackMessageTheme) FormatStartCommand(data *models.StartCommandReport) string {
	tpl := t.format(data)
	result, err := json.Marshal(tpl)
	if err != nil {
		// todo return { "text": err.String() }
	}

	return string(result)
}

func (t *DefaultSlackMessageTheme) format(data *models.StartCommandReport) *slackThemeTemplate {
	tpl := slackThemeTemplate{
		Text:        "Your Tuna Timer has been updated",
		Attachments: []slack.Attachment{},
	}

	if data.StartedTimer != nil {
		startTimerAttachment := t.defaultAttachment()
		startTimerAttachment.ThumbURL = t.StartCommandThumbURL
		startTimerAttachment.Footer = fmt.Sprintf(
			"Task ID: %s > <http://www.google.com|Open in Application>", data.StartedTimer.TaskHash)
		startTimerAttachment.Color = t.StartCommandColor

		startTimerAttachment.Fields = []slack.AttachmentField{}

		thisRoundField := slack.AttachmentField{
			Title: "This Round",
			Value: "00:01h",
			Short: true,
		}

		todayField := slack.AttachmentField{
			Title: "Today",
			Value: string(data.StartedTaskTotalForToday),
			Short: true,
		}

		startTimerAttachment.Fields = append(startTimerAttachment.Fields, thisRoundField, todayField)
		tpl.Attachments = append(tpl.Attachments, startTimerAttachment)
	}

	return &tpl
}

func (t *DefaultSlackMessageTheme) defaultAttachment() slack.Attachment {
	result := slack.Attachment{}
	result.MarkdownIn = t.MarkdownEnabledFor
	result.FooterIcon = t.FooterIcon
	return result
}
