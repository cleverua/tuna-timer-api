package commands

import (
	"context"

	"github.com/tuna-timer/tuna-timer-api/data"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/themes"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"time"
)

//Start - handles the '/timer stop` command received from Slack
type Start struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	userService  *data.UserService
	report       *models.StartCommandReport
	ctx          context.Context
}

func NewStart(ctx context.Context) *Start {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Start{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		userService:  data.NewUserService(session),
		report:       &models.StartCommandReport{},
		ctx:          ctx,
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
	team, project, err := c.teamService.EnsureTeamSetUp(&slackCommand)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	teamUser, err := c.userService.EnsureUser(team, slackCommand.UserID)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	c.report.Team = team
	c.report.Project = project
	c.report.TeamUser = teamUser

	timerToStop, err := c.timerService.GetActiveTimer(team.ID.Hex(), teamUser.ID.Hex())
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	if timerToStop != nil {
		if timerToStop.TaskName == slackCommand.Text && timerToStop.ProjectID == project.ID.Hex() {
			c.report.AlreadyStartedTimer = timerToStop
			c.report.AlreadyStartedTimerTotalForToday = c.timerService.TotalMinutesForTaskToday(timerToStop)
		} else {
			c.timerService.StopTimer(timerToStop)
			c.report.StoppedTimer = timerToStop
			c.report.StoppedTaskTotalForToday = c.timerService.TotalMinutesForTaskToday(timerToStop)
		}
	}
	if c.report.AlreadyStartedTimer == nil {
		startedTimer, err := c.timerService.StartTimer(team.ID.Hex(), project, teamUser, slackCommand.Text)
		if err != nil {
			// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
		}
		c.report.StartedTimer = startedTimer
		c.report.StartedTaskTotalForToday = c.timerService.TotalMinutesForTaskToday(c.report.StartedTimer)
	}

	c.report.UserTotalForToday = c.timerService.TotalUserMinutesForDay(teamUser.ID.Hex(), time.Now())

	return c.response()
}

func (c *Start) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = themes.NewDefaultSlackMessageTheme(c.ctx)
	content := theme.FormatStartCommand(c.report)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
