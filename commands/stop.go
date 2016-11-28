package commands

import (
	"context"

	"github.com/cleverua/tuna-timer-api/data"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/themes"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"time"
)

//Stop - handles the '/timer stop` command received from Slack
type Stop struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	userService  *data.UserService
	passService  *data.PassService
	report       *models.StopCommandReport
	ctx          context.Context
	theme        themes.SlackMessageTheme
}

func NewStop(ctx context.Context) *Stop {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Stop{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		userService:  data.NewUserService(session),
		passService:  data.NewPassService(session),
		report:       &models.StopCommandReport{},
		ctx:          ctx,
		theme:        utils.GetThemeFromContext(ctx).(themes.SlackMessageTheme),
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
		c.timerService.StopTimer(timerToStop)
		c.report.StoppedTimer = timerToStop
		c.report.StoppedTaskTotalForToday = c.timerService.TotalMinutesForTaskToday(timerToStop)
	}

	day := time.Now().Add(time.Duration(teamUser.SlackUserInfo.TZOffset) * time.Second)
	c.report.UserTotalForToday = c.timerService.TotalCompletedMinutesForDay(day.Year(), day.Month(), day.Day(), teamUser)

	return c.response()
}

func (c *Stop) response() *ResponseToSlack {
	return &ResponseToSlack{
		Body: []byte(c.theme.FormatStopCommand(c.report)),
	}
}
