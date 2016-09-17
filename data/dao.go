package data

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pavlo/slack-time/models"
)

// Dao - a tool for accessing persistent data
type Dao struct {
	DB *gorm.DB //should it be really exported?
}

// CreateTimer - method name is self-descriptive
func (dao *Dao) CreateTimer(user *models.TeamUser, task *models.Task) *models.Timer {
	result := &models.Timer{
		StartedAt:  time.Now(),
		TaskID:     task.ID,
		TeamUserID: user.ID,
	}
	dao.DB.Save(&result)

	return result
}

// FindNotFinishedTimerForUser - method name is self-descriptive
func (dao *Dao) FindNotFinishedTimerForUser(user *models.TeamUser) *models.Timer {
	result := &models.Timer{}
	request := dao.DB.Where("team_user_id = ? and finished_at is null and deleted_at is null", user.ID).First(&result)
	if !request.RecordNotFound() {
		return result
	}
	return nil
}

// FindOrCreateTaskByName - method name is self-descriptive
func (dao *Dao) FindOrCreateTaskByName(team *models.Team, project *models.Project, taskName string) *models.Task {
	result := &models.Task{}
	dao.DB.FirstOrCreate(&result, models.Task{ProjectID: project.ID, Name: taskName, TeamID: team.ID})
	if result.Hash == nil {
		dao.DB.Model(&result).Update("hash", taskSHA256(team, project, taskName))
	}
	return result
}

// FindOrCreateProjectBySlackChannelID - method name is self-descriptive
func (dao *Dao) FindOrCreateProjectBySlackChannelID(team *models.Team, slackChannelID string) *models.Project {
	result := &models.Project{}
	dao.DB.FirstOrCreate(&result, models.Project{TeamID: team.ID, SlackChannelID: slackChannelID})
	return result
}

// FindOrCreateTeamUserBySlackUserID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamUserBySlackUserID(team *models.Team, slackUserID string) *models.TeamUser {
	result := &models.TeamUser{}
	dao.DB.FirstOrCreate(&result, models.TeamUser{TeamID: team.ID, SlackUserID: slackUserID})
	return result
}

// FindOrCreateTeamBySlackTeamID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamBySlackTeamID(slackTeamID string) *models.Team {
	result := &models.Team{}
	dao.DB.FirstOrCreate(&result, models.Team{SlackTeamID: slackTeamID})
	return result
}

func taskSHA256(team *models.Team, project *models.Project, taskName string) string {
	hashSeed := fmt.Sprintf("%s%s%s", taskName, team.ID, project.ID)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hashSeed)))[0:8]
}
