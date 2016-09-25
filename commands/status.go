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

//Status - handles the '/timer status` command received from Slack
type Status struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	report       *models.StatusCommandReport
}

func NewStatus(ctx context.Context) *Status {
	session := utils.GetMongoSessionFromContext(ctx)

	status := &Status{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		report:       &models.StatusCommandReport{},
	}

	return status
}

// Handle - SlackCustomCommandHandler interface
// todo test it!
func (c *Status) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack {
	team, project, teamUser, err := c.teamService.EnsureTeamSetUp(&slackCommand)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	day := time.Now()
	c.report.PeriodName = "today"

	if slackCommand.Text == "yesterday" {
		day = day.AddDate(0, 0, -1)
		c.report.PeriodName = "yesterday"
	}

	c.report.Team = team
	c.report.Project = project
	c.report.TeamUser = teamUser

	tasks, err := c.timerService.GetCompletedTasksForDay(teamUser.ID.Hex(), day)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	if c.report.PeriodName == "today" {
		alreadyStartedTimer, _ := c.timerService.GetActiveTimer(team.ID.Hex(), teamUser.ID.Hex())

		if alreadyStartedTimer != nil {
			alreadyStartedTimer.Minutes = c.timerService.CalculateMinutesForActiveTimer(alreadyStartedTimer)
			c.report.AlreadyStartedTimer = alreadyStartedTimer
			c.report.AlreadyStartedTimerTotalForToday = c.timerService.TotalMinutesForTaskToday(alreadyStartedTimer)
		}
	}

	c.report.Tasks = tasks
	c.report.UserTotalForPeriod = c.timerService.TotalUserMinutesForDay(teamUser.ID.Hex(), time.Now())

	return c.response()
}

func (c *Status) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = themes.NewDefaultSlackMessageTheme()
	content := theme.FormatStatusCommand(c.report)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
