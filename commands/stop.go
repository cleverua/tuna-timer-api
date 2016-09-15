package commands

import (
	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Stop - starts timer for specific task
// If there is an other started task then it will be stopped
type Stop struct {
	CommandArguments
}

// Execute - implementation of Command interface
func (c Stop) Execute(env *utils.Environment) *CommandResult {
	result := &CommandResult{data: make(map[string]interface{})}

	team, user, project := CreateMainEntitiesIfNeeded(env, c.slackCommand)
	result.data["team"] = team
	result.data["user"] = user
	result.data["project"] = project

	dao := &data.Dao{DB: env.OrmDB}
	timerToFinish := dao.FindNotFinishedTimerForUser(user)

	if timerToFinish != nil {
		tasks := []data.Task{}
		env.OrmDB.Model(&timerToFinish).Association("Task").Find(&tasks)
		task := &tasks[0]

		MarkTimerAsFinished(task, timerToFinish)
		dao.DB.Save(&timerToFinish)
		dao.DB.Save(&task)

		result.data["finishedTimer"] = timerToFinish
		result.data["task"] = task
	}

	return result
}

// GetName return the name of this command
func (c Stop) GetName() string {
	return CommandNameStop
}
