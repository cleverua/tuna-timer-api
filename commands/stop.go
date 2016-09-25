package commands

import (
	"context"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/themes"
	"github.com/pavlo/slack-time/utils"
	"gopkg.in/mgo.v2"
	"time"
)

//Stop - handles the '/timer stop` command received from Slack
type Stop struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	report       *models.StopCommandReport
}

func NewStop(ctx context.Context) *Stop {
	session := utils.GetMongoSessionFromContext(ctx)

	start := &Stop{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		report:       &models.StopCommandReport{},
	}

	return start
}

// cases:
// 1. Successfully stopped a timer
// 2. No currently ticking timer existed
// 3. Any other errors

// Handle - SlackCustomCommandHandler interface
func (c *Stop) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack {

	team, project, teamUser, err := c.teamService.EnsureTeamSetUp(&slackCommand)
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

	c.report.UserTotalForToday = c.timerService.TotalUserMinutesForDay(teamUser.ID.Hex(), time.Now())

	return c.response()
}

func (c *Stop) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = themes.NewDefaultSlackMessageTheme()
	content := theme.FormatStopCommand(c.report)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
