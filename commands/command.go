package commands

import (
	"fmt"
	"strings"
)

type CommandArguments struct {
	arguments map[string]string
}

// Command - public interface
type Command interface {
	Execute() CommandResult
}

// CommandResult - is what command's Execute method returns
type CommandResult struct {
	data map[string]interface{}
}

// Get - looks up specific implementation of Command that matches user input
func Get(userInput string) (Command, error) {
	if strings.HasPrefix(userInput, "start") {
		cmd := Start{}
		cmd.arguments = make(map[string]string)
		cmd.arguments["main"] = stripCommandNameFromUserInput("start", userInput)

		return cmd, nil
	}
	return nil, fmt.Errorf("Failed to look up a command for `%s` name", userInput)
}

func stripCommandNameFromUserInput(commandName, userInput string) string {
	result := userInput[len(commandName):]
	return strings.TrimSpace(result)
}
