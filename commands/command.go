package commands

import (
	"fmt"
	"strings"
)

// Command - public interface
type Command interface {
	Execute() CommandResult
}

// CommandResult - is what command's Execute method returns
type CommandResult struct {
	data map[string] interface{}
}

// Get - looks up specific implementation of Command that matches user input
func Get(userInput string) (Command, error) {
	if strings.HasPrefix(userInput, "start") {
		return Start{}, nil
	}
	return nil, fmt.Errorf("Failed to look up a command for `%s` name", userInput)
}
