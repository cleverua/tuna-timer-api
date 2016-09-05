package commands

import (
	"fmt"
	"strings"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Command - public interface
type Command interface {
	Execute(env *utils.Environment) *CommandResult
}

// CommandArguments - this is what a Command needs to operate
type CommandArguments struct {
	slackCommand data.SlackCommand
	rawCommand   string
}

// CommandResult - is what command's Execute method returns
type CommandResult struct {
	data map[string]interface{}
}

// Get - looks up a Command that would serve the corresponding Slack Command
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
		rawCommand:   stripCommandNameFromUserInput(commandName, slackCommand.Text),
	}
}
