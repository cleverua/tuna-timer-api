package commands

import (
	"context"

	"fmt"
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
	passService  *data.PassService
	report       *models.StartCommandReport
	ctx          context.Context
	theme        themes.SlackMessageTheme
}

func NewStart(ctx context.Context) *Start {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Start{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		userService:  data.NewUserService(session),
		passService:  data.NewPassService(session),
		report:       &models.StartCommandReport{},
		ctx:          ctx,
		theme:        utils.GetThemeFromContext(ctx).(themes.SlackMessageTheme),
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

	if slackCommand.Text == "" {
		return c.errorResponse(
			fmt.Sprintf("Task name not provided! The correct command would look like: \n>`%s start My super exciting task`", slackCommand.Command),
		)
	}

	team, project, err := c.teamService.EnsureTeamSetUp(&slackCommand)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	teamUser, err := c.userService.EnsureUser(team, slackCommand.UserID)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	pass, err := c.passService.EnsurePass(team, teamUser, project)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	c.report.Team = team
	c.report.Project = project
	c.report.TeamUser = teamUser
	c.report.Pass = pass

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

	day := time.Now().Add(time.Duration(teamUser.SlackUserInfo.TZOffset) * time.Second)
	c.report.UserTotalForToday = c.timerService.TotalCompletedMinutesForDay(day.Year(), day.Month(), day.Day(), teamUser)

	return c.response()
}

func (c *Start) response() *ResponseToSlack {
	return &ResponseToSlack{
		Body: []byte(c.theme.FormatStartCommand(c.report)),
	}
}

func (c *Start) errorResponse(errorMessage string) *ResponseToSlack {
	return &ResponseToSlack{
		Body: []byte(c.theme.FormatError(errorMessage)),
	}
}
