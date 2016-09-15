package commands

import (
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Start - starts timer for specific task
// If there is an other started task then it will be stopped
type Start struct {
	CommandArguments
}

// Execute - implementation of Command interface
func (c Start) Execute(env *utils.Environment) *CommandResult {
	result := &CommandResult{data: make(map[string]interface{})}

	team, user, project := CreateMainEntitiesIfNeeded(env, c.slackCommand)
	result.data["team"] = team
	result.data["user"] = user
	result.data["project"] = project

	dao := &data.Dao{DB: env.OrmDB}
	task := dao.FindOrCreateTaskByName(team, project, c.rawCommand)
	result.AffectedTask = task

	timerToFinish := dao.FindNotFinishedTimerForUser(user)
	if timerToFinish != nil {
		MarkTimerAsFinished(task, timerToFinish)
		dao.DB.Save(&timerToFinish)
		dao.DB.Save(&task)
		result.data["finishedTimer"] = timerToFinish
	}

	// Check to see if there's already a timer for this user started for this task
	timer := dao.CreateTimer(user, task)
	result.data["startedTimer"] = timer

	return result
}

// GetName return the name of this command
func (c Start) GetName() string {
	return CommandNameStart
}
