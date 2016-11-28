package data

import (
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"testing"

	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"

	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"

)

func TestUserRepository(t *testing.T) {
	gosuite.Run(t, &UserRepositoryTestSuite{Is: is.New(t)})
}

func (s *UserRepositoryTestSuite) TestFindByExternalID(t *testing.T) {
	user := &models.TeamUser{
		TeamID:           "team-id",
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
		SlackUserInfo: &slack.User{
			IsAdmin: true,
		},
	}

	u, err := s.repository.save(user)

	s.Nil(err)
	s.NotNil(u)

	loadedUser, err := s.repository.FindByExternalID("ext-id")
	s.Nil(err)
	s.NotNil(loadedUser)
	s.Equal(loadedUser.ExternalUserID, "ext-id")
	s.True(loadedUser.SlackUserInfo.IsAdmin)
}

func (s *UserRepositoryTestSuite) TestSave(t *testing.T) {
	user := &models.TeamUser{
		TeamID:           "team-id",
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
		SlackUserInfo: &slack.User{
			IsAdmin: true,
		},
	}

	u, err := s.repository.save(user)
	s.Nil(err)

	u.SlackUserInfo.IsAdmin = false
	_, err = s.repository.save(u)
	s.Nil(err)

	loadedUser, err := s.repository.FindByExternalID("ext-id")
	s.Nil(err)
	s.NotNil(loadedUser)
	s.Equal(loadedUser.ExternalUserName, "ext-name")
	s.False(loadedUser.SlackUserInfo.IsAdmin)
	s.Equal(loadedUser.ModelVersion, models.ModelVersionTeamUser)
}

func (s *UserRepositoryTestSuite) TestFindByExternalIDNotExist(t *testing.T) {
	resultTeam, err := s.repository.FindByExternalID("external-id")
	s.Nil(err)
	s.Nil(resultTeam)
}

func (s *UserRepositoryTestSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repository = NewUserRepository(s.session)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *UserRepositoryTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *UserRepositoryTestSuite) TearDown() {}


type UserRepositoryTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	repository *UserRepository
}
