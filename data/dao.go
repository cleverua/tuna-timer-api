package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Dao - a tool for accessing persistent data
type Dao struct {
	DB *gorm.DB
}

// CreateTimer - method name is self-descriptive
func (dao *Dao) CreateTimer(user *TeamUser, task *Task) *Timer {
	result := &Timer{
		StartedAt:  time.Now(),
		TaskID:     task.ID,
		TeamUserID: user.ID,
	}
	dao.DB.Save(&result)
	return result
}

// FindNotFinishedTimerForUser - method name is self-descriptive
func (dao *Dao) FindNotFinishedTimerForUser(user *TeamUser) *Timer {
	result := &Timer{}
	request := dao.DB.Where("team_user_id = ? and finished_at is null and deleted_at is null", user.ID).First(&result)
	if !request.RecordNotFound() {
		return result
	}
	return nil
}

// FindOrCreateTaskByName - method name is self-descriptive
func (dao *Dao) FindOrCreateTaskByName(team *Team, project *Project, taskName string) *Task {
	result := &Task{}
	dao.DB.FirstOrCreate(&result,
		Task{ProjectID: project.ID, Name: taskName, TeamID: team.ID})
	return result
}

// FindOrCreateProjectBySlackChannelID - method name is self-descriptive
func (dao *Dao) FindOrCreateProjectBySlackChannelID(team *Team, slackChannelID string) *Project {
	result := &Project{}
	dao.DB.FirstOrCreate(&result, Project{TeamID: team.ID, SlackChannelID: slackChannelID})
	return result
}

// FindOrCreateTeamUserBySlackUserID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamUserBySlackUserID(team *Team, slackUserID string) *TeamUser {
	result := &TeamUser{}
	dao.DB.FirstOrCreate(&result, TeamUser{TeamID: team.ID, SlackUserID: slackUserID})
	return result
}

// FindOrCreateTeamBySlackTeamID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamBySlackTeamID(slackTeamID string) *Team {
	result := &Team{}
	dao.DB.FirstOrCreate(&result, Team{SlackTeamID: slackTeamID})
	return result
}
