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

//Status - handles the '/timer status` command received from Slack
type Status struct {
	session      *mgo.Session
	teamService  *data.TeamService
	timerService *data.TimerService
	userService  *data.UserService
	report       *models.StatusCommandReport
	ctx          context.Context
}

func NewStatus(ctx context.Context) *Status {
	session := utils.GetMongoSessionFromContext(ctx)

	status := &Status{
		session:      session,
		teamService:  data.NewTeamService(session),
		timerService: data.NewTimerService(session),
		userService:  data.NewUserService(session),
		report:       &models.StatusCommandReport{},
		ctx:          ctx,
	}

	return status
}

// Handle - SlackCustomCommandHandler interface
// todo test it!
func (c *Status) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *ResponseToSlack {
	team, project, err := c.teamService.EnsureTeamSetUp(&slackCommand)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	teamUser, err := c.userService.EnsureUser(team, slackCommand.UserID)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}

	day := time.Now().Add(time.Duration(teamUser.SlackUserInfo.TZOffset) * time.Minute)
	c.report.PeriodName = "today"

	if slackCommand.Text == "yesterday" {
		day = day.AddDate(0, 0, -1)
		c.report.PeriodName = "yesterday"
	}

	c.report.Team = team
	c.report.Project = project
	c.report.TeamUser = teamUser

	tasks, err := c.timerService.GetCompletedTasksForDay(day.Year(), day.Month(), day.Day(), teamUser)
	if err != nil {
		// todo: format a decent Slack error message so user knows what's wrong and how to solve the issue
	}
	c.report.Tasks = tasks

	if c.report.PeriodName == "today" {
		alreadyStartedTimer, _ := c.timerService.GetActiveTimer(team.ID.Hex(), teamUser.ID.Hex())

		if alreadyStartedTimer != nil {
			alreadyStartedTimer.Minutes = c.timerService.CalculateMinutesForActiveTimer(alreadyStartedTimer)
			c.report.AlreadyStartedTimer = alreadyStartedTimer
			c.report.AlreadyStartedTimerTotalForToday = c.timerService.TotalMinutesForTaskToday(alreadyStartedTimer)
		}
	}

	c.report.UserTotalForPeriod = c.timerService.TotalUserMinutesForDay(teamUser.ID.Hex(), time.Now())

	return c.response()
}

func (c *Status) response() *ResponseToSlack {
	var theme themes.SlackMessageTheme = themes.NewDefaultSlackMessageTheme(c.ctx)
	content := theme.FormatStatusCommand(c.report)

	return &ResponseToSlack{
		Body: []byte(content),
	}
}
