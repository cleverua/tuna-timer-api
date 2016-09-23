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
		sa := t.defaultAttachment()
		sa.Text = fmt.Sprintf("Started for: %s", data.StartedTimer.TaskName)
		sa.ThumbURL = t.StartCommandThumbURL
		sa.Footer = fmt.Sprintf(
			"Task ID: %s > <http://www.google.com|Open in Application>", data.StartedTimer.TaskHash)
		sa.Color = t.StartCommandColor

		sa.Fields = []slack.AttachmentField{}

		thisRoundField := slack.AttachmentField{
			Title: "This Round",
			Value: utils.FormatDuration(time.Duration(1 * time.Minute)),
			Short: true,
		}

		todayField := slack.AttachmentField{
			Title: "Today",
			Value: utils.FormatDuration(time.Duration(data.StartedTaskTotalForToday * int(time.Minute))),
			Short: true,
		}

		sa.Fields = append(sa.Fields, thisRoundField, todayField)
		tpl.Attachments = append(tpl.Attachments, sa)
	}

	return &tpl
}

func (t *DefaultSlackMessageTheme) defaultAttachment() slack.Attachment {
	result := slack.Attachment{}
	result.MarkdownIn = t.MarkdownEnabledFor
	result.FooterIcon = t.FooterIcon
	return result
}
