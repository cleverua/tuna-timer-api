package commands

import (
	"fmt"
	"strings"
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

type CommandArguments struct {
	slackCommand data.SlackCommand
	rawCommand string
}

// Command - public interface
type Command interface {
	Execute(env *utils.Environment) *CommandResult
}

// CommandResult - is what command's Execute method returns
type CommandResult struct {
	data map[string]interface{}
}

func Get(slackCommand data.SlackCommand) (Command, error) {
	userInput := slackCommand.Text
	if strings.HasPrefix(userInput, "start") {
		cmd := Start{CommandArguments: createCommandArguments(slackCommand, "start")}
		return cmd, nil
	}
	return nil, fmt.Errorf("Failed to look up a command for `%s` name", userInput)
}

func stripCommandNameFromUserInput(commandName, userInput string) string {
	result := userInput[len(commandName):]
	return strings.TrimSpace(result)
}

func createCommandArguments(slackCommand data.SlackCommand, commandName string) CommandArguments {
	return CommandArguments{
		slackCommand: slackCommand,
		rawCommand: stripCommandNameFromUserInput(commandName, slackCommand.Text),
	}
}
