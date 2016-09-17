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

// FindOrCreateTeamUser finds or creates a new TeamUser for given Team
func (s *DataService) FindOrCreateTeamUser(db *gorm.DB, team *models.Team, slackUserID, slackUserName string) (*models.TeamUser, error) {
	result := &models.TeamUser{}

	request := db.Where("team_id = ? and user_id = ?", team.ID, slackUserID).First(&result)
	if request.RecordNotFound() {

		result.Name = slackUserName
		result.TeamID = team.ID
		result.SlackUserID = slackUserID
		db.Create(result)

		if db.Error != nil {
			return nil, db.Error
		}
	}
	return result, nil
}

// func (s *DataService) CreateTeamAndUserAndProject(db *gorm.DB, slackCommand models.SlackCustomCommand) (*models.Team, *models.TeamUser, *models.Project) {
// 	dao := &Dao{DB: db}
// 	team := dao.FindOrCreateTeamBySlackTeamID(slackCommand.TeamID)
// 	user := dao.FindOrCreateTeamUserBySlackUserID(team, slackCommand.UserID)
// 	project := dao.FindOrCreateProjectBySlackChannelID(team, slackCommand.ChannelID)
// 	return team, user, project
// }
