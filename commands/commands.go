package commands

import (
	"context"
	"fmt"
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
	subCommand := slackCommand.SubCommand

	if subCommand == CommandNameStart {
		cmd := NewStart(ctx)
		return cmd, nil
	} else if subCommand == CommandNameStop {
		cmd := NewStop(ctx)
		return cmd, nil
	} else if subCommand == CommandNameStatus {
		cmd := NewStatus(ctx)
		return cmd, nil
	}
	return nil, fmt.Errorf("Failed to look up a handler for `%s` name", subCommand)
}
