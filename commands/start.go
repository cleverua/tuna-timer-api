package commands

import (
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
	"time"
)

// Start - starts timer for specific task
// If there is an other started task then it will be stopped
type Start struct {
	CommandArguments
}

func (c Start) Execute(env *utils.Environment) *CommandResult {

	dao := &data.Dao{DB: env.OrmDB}

	// Create team if needed
	team := dao.FindOrCreateTeamBySlackTeamID(c.slackCommand.TeamID)

	// Create user if needed
	user := dao.FindOrCreateTeamUserBySlackUserID(team, c.slackCommand.UserID)

	// Create project if needed
	project := dao.FindOrCreateProjectBySlackChannelID(team, c.slackCommand.ChannelID)

	// Create project if needed
	task := dao.FindOrCreateTaskByName(team, project, c.rawCommand)

	finishedTimer := dao.FindNotFinishedTimerForUser(user)
	now := time.Now()
	if finishedTimer != nil {
		finishedTimer.FinishedAt = &now
		dao.DB.Save(&finishedTimer)
	}

	// Check to see if there's already a timer for this user started for this task
	timer := dao.CreateTimer(user, task)

	m := make(map[string]interface{})
	m["team"] 		= team
	m["user"]		= user
	m["project"]		= project
	m["task"]		= task
	m["finishedTimer"] 	= finishedTimer
	m["startedTimer"]	= timer

	result := CommandResult{
		data: m,
	}

	return &result
}
