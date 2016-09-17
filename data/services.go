package data

import (
	"github.com/jinzhu/gorm"
	"github.com/pavlo/slack-time/models"
)

//DataService todo
type DataService struct {
}

//CreateDataService todo
func CreateDataService() *DataService {
	return &DataService{}
}

func (s *DataService) CreateTeamAndUserAndProject(db *gorm.DB, slackCommand models.SlackCustomCommand) (*models.Team, *models.TeamUser, *models.Project) {
	dao := &Dao{DB: db}
	team := dao.FindOrCreateTeamBySlackTeamID(slackCommand.TeamID)
	user := dao.FindOrCreateTeamUserBySlackUserID(team, slackCommand.UserID)
	project := dao.FindOrCreateProjectBySlackChannelID(team, slackCommand.ChannelID)
	return team, user, project
}
