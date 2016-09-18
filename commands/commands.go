package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/pavlo/slack-time/models"
)

const (
	// CommandNameStart hodls the name of 'start' command
	CommandNameStart = "start"

	// CommandNameStop hodls the name of 'stop' command
	CommandNameStop = "stop"

	// CommandNameStatus hodls the name of 'status' command
	CommandNameStatus = "status"
)

// SlackCustomCommandHandlerResult todo
type SlackCustomCommandHandlerResult struct {
	Body []byte
}

// SlackCustomCommandHandler todo
type SlackCustomCommandHandler interface {
	Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *SlackCustomCommandHandlerResult
}

// LookupHandler todo
func LookupHandler(slackCommand models.SlackCustomCommand) (SlackCustomCommandHandler, error) {
	userInput := slackCommand.Text
	if strings.HasPrefix(userInput, CommandNameStart) {
		cmd := &Stop{}
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
