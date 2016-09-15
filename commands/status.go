package commands

import (
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Status - starts timer for specific task
// If there is an other started task then it will be stopped
type Status struct {
	CommandArguments
}

// Execute - implementation of Command interface
func (c Status) Execute(env *utils.Environment) *CommandResult {
	result := &CommandResult{data: make(map[string]interface{})}

	team, user, project := CreateMainEntitiesIfNeeded(env, c.slackCommand)
	result.data["team"] = team
	result.data["user"] = user
	result.data["project"] = project

	dao := &data.Dao{DB: env.OrmDB}
	timer := dao.FindNotFinishedTimerForUser(user)

	if timer != nil {
		timer.Minutes = GetMinutesTimerRun(timer)
		result.data["timer"] = timer

		tasks := []data.Task{}
		env.OrmDB.Model(&timer).Association("Task").Find(&tasks)
		task := &tasks[0]
		result.data["task"] = task
	}

	return result
}

// GetName return the name of this command
func (c Status) GetName() string {
	return CommandNameStatus
}
