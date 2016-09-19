package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/pavlo/slack-time/models"
)

const (
	// CommandNameStart holds the name of 'start' command
	CommandNameStart = "start"

	// CommandNameStop holds the name of 'stop' command
	CommandNameStop = "stop"

	// CommandNameStatus holds the name of 'status' command
	CommandNameStatus = "status"
)

// SlackCustomCommandHandlerResult todo
type ResponseToSlack struct {
	Body []byte
}

// SlackCustomCommandHandler todo
type SlackCustomCommandHandler interface {
	Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack
}

// LookupHandler todo
func LookupHandler(ctx context.Context, slackCommand models.SlackCustomCommand) (SlackCustomCommandHandler, error) {
	userInput := slackCommand.Text
	if strings.HasPrefix(userInput, CommandNameStart) {
		cmd := NewStart(ctx)
		return cmd, nil
	} else if strings.HasPrefix(userInput, CommandNameStop) {
		cmd := &Stop{}
		return cmd, nil
	} else if strings.HasPrefix(userInput, CommandNameStatus) {
		cmd := &Stop{}
		return cmd, nil
	}
	return nil, fmt.Errorf("Failed to look up a handler for `%s` name", userInput)
}
