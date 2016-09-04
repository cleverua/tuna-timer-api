package data

import "github.com/jinzhu/gorm"

type Dao struct {
	db *gorm.DB
}

func (dao *Dao) FindOrCreateTeamBySlackTeamId(slackTeamId string) *Team {
	result := &Team{}
	dao.db.FirstOrCreate(&result, Team{SlackTeamId: slackTeamId})
	return result
}
