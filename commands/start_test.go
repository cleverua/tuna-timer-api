package commands

import (
	"fmt"
	"testing"
)

// simple command is a command that has only one (thus main) argument
func TestGetSimpleStartCommand(t *testing.T) {
	cmd, err := Get("start DDP-256 Add migration for user_id column")

	if err != nil {
		t.Error("Unexpected error thrown!")
	}

	commandType := fmt.Sprintf("%T", cmd)
	if ("commands.Start" != commandType) {
		t.Errorf("Expected Get method to return an instance of " +
			"commands.Start command, but it was %s", commandType)
	}
}
