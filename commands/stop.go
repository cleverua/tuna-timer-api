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

//Stop - handles the '/timer stop` command received from Slack
type Stop struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	userService  *data.UserService
	report       *models.StopCommandReport
	ctx          context.Context
}

func NewStop(ctx context.Context) *Stop {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Stop{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		userService:  data.NewUserService(session),
		report:       &models.StopCommandReport{},
		ctx:          ctx,
	}

	return start
}

// cases:
// 1. Successfully stopped a timer
// 2. No currently ticking timer existed
// 3. Any other errors

// Handle - SlackCustomCommandHandler interface
func (c *Stop) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack {

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
		c.timerService.StopTimer(timerToStop)
		c.report.StoppedTimer = timerToStop
		c.report.StoppedTaskTotalForToday = c.timerService.TotalMinutesForTaskToday(timerToStop)
	}

	day := time.Now().Add(time.Duration(teamUser.SlackUserInfo.TZOffset) * time.Second)
	c.report.UserTotalForToday = c.timerService.TotalUserMinutesForDay(day.Year(), day.Month(), day.Day(), teamUser)

	return c.response()
}

func (c *Stop) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = themes.NewDefaultSlackMessageTheme(c.ctx)
	content := theme.FormatStopCommand(c.report)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
