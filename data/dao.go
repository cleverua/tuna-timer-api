package data

import "github.com/jinzhu/gorm"

// Dao - a tool for accessing persisten data
type Dao struct {
	db *gorm.DB
}

// FindOrCreateProjectBySlackChannelID - method name is self-descriptive
func (dao *Dao) FindOrCreateProjectBySlackChannelID(team *Team, slackChannelID string) *Project {
	result := &Project{}
	dao.db.FirstOrCreate(&result, Project{TeamID: team.ID, SlackChannelID: slackChannelID})
	return result
}

// FindOrCreateTeamUserBySlackUserID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamUserBySlackUserID(team *Team, slackUserID string) *TeamUser {
	result := &TeamUser{}
	dao.db.FirstOrCreate(&result, TeamUser{TeamID: team.ID, SlackUserID: slackUserID})
	return result
}

// FindOrCreateTeamBySlackTeamID - method name is self-descriptive
func (dao *Dao) FindOrCreateTeamBySlackTeamID(slackTeamID string) *Team {
	result := &Team{}
	dao.db.FirstOrCreate(&result, Team{SlackTeamID: slackTeamID})
	return result
}
