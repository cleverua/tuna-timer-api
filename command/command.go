package command

import (
	"fmt"
	"strings"
)

// Command - public interface
type Command interface {
	Execute() string
}

// Get - looks up specific implementation of Command that matches user input
func Get(userInput string) (Command, error) {
	if strings.HasPrefix(userInput, "start") {
		return StartCommand{}, nil
	}
	return nil, fmt.Errorf("Failed to look up a command for `%s` name", userInput)
}
