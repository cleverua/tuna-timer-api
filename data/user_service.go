package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserService struct {
	repository *UserRepository
}

func NewUserService(session *mgo.Session) *UserService {
	return &UserService{
		repository: NewUserRepository(session),
	}
}

func (s *UserService) EnsureUser(team *models.Team, externalUserID string) (*models.TeamUser, error) {
	slackAPI := slack.New(team.SlackOAuthResponse.AccessToken)
	slackAPI.SetDebug(true)

	slackUserData, err := slackAPI.GetUserInfo(externalUserID)
	if err != nil {
		return nil, err
	}

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		CreatedAt: time.Now(),
		ExternalUserID:externalUserID,
		ExternalUserName:slackUserData.Name,
		SlackUserData:slackUserData,
		TeamID: team.ID.Hex(),
	}

	user, err = s.repository.save(user);
	if err != nil {
		return nil, err
	}

	return user, nil
}
