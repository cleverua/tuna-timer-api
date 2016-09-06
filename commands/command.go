package commands

import (
	"fmt"
	"math"
	"strings"
	"time"

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
	} else if strings.HasPrefix(userInput, "stop") {
		cmd := Stop{CommandArguments: createCommandArguments(slackCommand, "stop")}
		return cmd, nil
	} else if strings.HasPrefix(userInput, "status") {
		cmd := Status{CommandArguments: createCommandArguments(slackCommand, "status")}
		return cmd, nil
	}
	return nil, fmt.Errorf("Failed to look up a command for `%s` name", userInput)
}

// MarkTimerAsFinished performs housekeeping when a timer is finishing:
// 1. Calculates how many minutes the task was started
// 2. Adds the minutes to task's TotalMinutes
// Note, it does not update objects in DB!
func MarkTimerAsFinished(task *data.Task, t *data.Timer) {
	now := time.Now()
	duration := time.Since(t.StartedAt)
	t.FinishedAt = &now
	t.Minutes = int(math.Floor(duration.Minutes()))
	task.TotalMinutes = task.TotalMinutes + t.Minutes
}

// CreateMainEntitiesIfNeeded lazily creates Team, TeamUser and Project
func CreateMainEntitiesIfNeeded(env *utils.Environment, slackCommand data.SlackCommand) (*data.Team, *data.TeamUser, *data.Project) {
	dao := &data.Dao{DB: env.OrmDB}
	team := dao.FindOrCreateTeamBySlackTeamID(slackCommand.TeamID)
	user := dao.FindOrCreateTeamUserBySlackUserID(team, slackCommand.UserID)
	project := dao.FindOrCreateProjectBySlackChannelID(team, slackCommand.ChannelID)
	return team, user, project
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
