package main

import "fmt"

import "github.com/pavlo/slack-time/commands"

func main() {
	command, err := commands.Get("start DDP-256 Add migration for user_id column")
	if err != nil {
		fmt.Print("Failed to look up a command!")
	} else {
		fmt.Print(command.Execute())
	}
}
