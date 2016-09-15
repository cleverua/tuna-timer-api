package formatter

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/pavlo/slack-time/commands"
	"github.com/pavlo/slack-time/utils"
)

// ResponseFormatter - responsble for getting CommandResult to convert it to something that can be marshalled
type ResponseFormatter interface {
	Format(env *utils.Environment, cmd *commands.CommandResult) interface{}
}

// SlackMessage - a message that would be sent back in response to SlackCommand
type SlackMessage struct {
	Text        string              `json:"text,omitempty"`
	Attachments []*slack.Attachment `json:"attachments,omitempty"`
}

// Lookup finds a formatter for given slack command
func Lookup(cmdName string) (ResponseFormatter, error) {
	if cmdName == commands.CommandNameStart {
		return &Start{}, nil
	}

	return nil, fmt.Errorf("Failed to look up a formatter for `%s` command", cmdName)
}

func defaultAttachment() *slack.Attachment {
	result := &slack.Attachment{}
	result.MarkdownIn = []string{"text", "pretext"}
	result.FooterIcon = "http://icons.iconarchive.com/icons/martin-berube/flat-animal/48/tuna-icon.png"
	return result
}

func createField(title, value string, short bool) slack.AttachmentField {
	return slack.AttachmentField{
		Short: short,
		Title: title,
		Value: value,
	}
}

func todayTotalAttachment(text string) *slack.Attachment {
	return &slack.Attachment{
		AuthorName: text,
		Color:      "#FFFFFF",
	}
}
