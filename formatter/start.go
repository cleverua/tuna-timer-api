package formatter

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/utils"
)

// Start defined the struct to hold the Format method for start command
type Start struct {
}

// Format - ResponseFormatter interface implementation
func (s *Start) Format(env *utils.Environment, commandResult *commands.CommandResult) interface{} {

	// team, user, project, task, finishedTimer?, startedTimer

	task := commandResult.AffectedTask

	attachments := []*slack.Attachment{}
	attachments = append(attachments, &slack.Attachment{
		Text:     fmt.Sprintf("Started for: %s", task.Name),
		ThumbURL: "http://icons.iconarchive.com/icons/graphicloads/100-flat/128/new-icon.png",
		Footer:   fmt.Sprintf("Task ID: %s > <http://www.disney.com|Open in Application>", *task.Hash),
		Color:    "#FB6E04",
		Fields: []slack.AttachmentField{
			createField("This Round", "00:01h", true),
			createField("Today", "02:12h", true),
		},
	})
	attachments = append(attachments, todayTotalAttachment("Your total for today is 05:57 hours"))

	message := &SlackMessage{
		Text:        "The Tuna Timer has been",
		Attachments: attachments,
	}

	return message
}
