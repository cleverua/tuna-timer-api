package commands

import (
	"context"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
	"github.com/pavlo/slack-time/themes"
	"gopkg.in/mgo.v2"
)

//Start - handles the '/timer stop` command received from Slack
type Start struct {
	session *mgo.Session
	teamService *data.TeamService
	timerService *data.TimerService
	inventory *models.StartCommandInventory
}

func NewStart(ctx context.Context) *Start {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Start{
		session: session,
		teamService: data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		inventory: &models.StartCommandInventory{},
	}

	return start
}

// cases:
// * Has nothing to do if requested task is already started
// * Successfully started a new timer while there were no existing running
// * Successfully started a new timer and stopped the previous one
// * Timer is already in progress
// * The started one has the same taskName thus the task is actually resumed
// * Any other errors

// Handle - SlackCustomCommandHandler interface
func (c *Start) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack {
	team, project, teamUser, err := c.teamService.EnsureTeamSetUp(&slackCommand)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	c.inventory.Team = team
	c.inventory.Project = project
	c.inventory.TeamUser = teamUser

	timerToStop, err := c.timerService.GetActiveTimer(team.ID.Hex(), teamUser.ID.Hex())
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	if timerToStop != nil {
		if timerToStop.TaskName == slackCommand.Text {
			c.inventory.AlreadyStartedTimer = timerToStop
			c.inventory.AlreadyStartedTimerTotalForToday = c.timerService.TaskTotalMinutesForToday(timerToStop)
		} else {
			c.timerService.StopTimer(timerToStop)
			c.inventory.StoppedTimer = timerToStop
			c.inventory.StoppedTaskTotalForToday = c.timerService.TaskTotalMinutesForToday(timerToStop)
		}
	}

	if c.inventory.AlreadyStartedTimer != nil {
		startedTimer, err := c.timerService.StartTimer(team.ID.Hex(), project.ID.Hex(), teamUser.ID.Hex(), slackCommand.Text)
		if err != nil {
			// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
		}
		c.inventory.StartedTimer = startedTimer
		c.inventory.StartedTaskTotalForToday = c.timerService.TaskTotalMinutesForToday(c.inventory.StartedTimer)
	}

	return c.response()
}

func (c *Start) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = &themes.DefaultSlackMessageTheme{}
	content := theme.FormatStartCommand(c.inventory)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
