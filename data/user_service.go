package data

import (
	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"time"
)

type UserService struct {
	repository *UserRepository
	slackAPI   userServerSlackAPI
}

func NewUserService(session *mgo.Session) *UserService {
	return &UserService{
		repository: NewUserRepository(session),
		slackAPI:   &userServiceSlackAPIImpl{},
	}
}

func (s *UserService) FindByID(id string) (*models.TeamUser, error){
	user, err := s.repository.FindByID(id)
	if err == mgo.ErrNotFound {
		return nil, err
	}
	return user, err
}

func (s *UserService) EnsureUser(team *models.Team, externalUserID string) (*models.TeamUser, error) {
	user, err := s.repository.FindByExternalID(externalUserID)
	if err != nil {
		return nil, err
	}

	if user == nil {

		slackUserData, err := s.slackAPI.GetUserInfo(team, externalUserID)
		if err != nil {
			return nil, err
		}

		user = &models.TeamUser{
			CreatedAt:        time.Now(),
			ExternalUserID:   externalUserID,
			ExternalUserName: slackUserData.Name,
			SlackUserInfo:    slackUserData,
			TeamID:           team.ID.Hex(),
		}

		user, err = s.repository.Save(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// UpdateSlackUserInfo - finds or creates TeamUser record with associated user data gathered from Slack
//func (s *UserService) UpdateSlackUserInfo(team *models.Team, user *models.TeamUser) (*models.TeamUser, error) {
//	info, err := s.slackAPI.GetUserInfo(team, user.ExternalUserID)
//	if err != nil {
//		return nil, err
//	}
//
//	user.SlackUserInfo = info
//	user, err = s.repository.save(user)
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}

// A wrapper around slack API used by this service. Unit tests will inject their own impl. of this to bypass network calls to Slack
type userServerSlackAPI interface {
	GetUserInfo(team *models.Team, externalUserID string) (*slack.User, error)
}

type userServiceSlackAPIImpl struct{}

func (u *userServiceSlackAPIImpl) GetUserInfo(team *models.Team, externalUserID string) (*slack.User, error) {
	slackAPI := slack.New(team.SlackOAuth.AccessToken)
	slackAPI.SetDebug(true)
	return slackAPI.GetUserInfo(externalUserID)
}
