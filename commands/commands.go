package commands

import (
	"context"
	"fmt"
	"github.com/cleverua/tuna-timer-api/models"
)

const (
	CommandNameStart  = "start"
	CommandNameStop   = "stop"
	CommandNameStatus = "status"
)

type ResponseToSlack struct {
	Body []byte
}

type SlackCustomCommandHandler interface {
	Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack
}

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
	return nil, fmt.Errorf("Unknown command `%s`!", subCommand)
}
