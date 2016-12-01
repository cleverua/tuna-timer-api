package data

import (
	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/tylerb/is.v1"
	"log"
	"testing"
	"github.com/pavlo/gosuite"
)

func TestUserService(t *testing.T) {
	gosuite.Run(t, &UserServiceTestSuite{Is: is.New(t)})
}

type userServiceSlackAPIImplTest struct {
	err  error
	user *slack.User
}

func (u *userServiceSlackAPIImplTest) GetUserInfo(team *models.Team, externalUserID string) (*slack.User, error) {
	return u.user, u.err
}

func newUserServiceSlackAPIImplTest(user *slack.User, err error) *userServiceSlackAPIImplTest {
	return &userServiceSlackAPIImplTest{
		user: user,
		err:  err,
	}
}

func (s *UserServiceTestSuite) TestEnsureUserNew(t *testing.T) {

	teamID := bson.NewObjectId()
	team := &models.Team{
		ID: teamID,
	}
	service := NewUserService(s.session)
	service.slackAPI = newUserServiceSlackAPIImplTest(
		&slack.User{
			IsAdmin:  true,
			Name:     "test-user",
			TZOffset: -2000,
		},
		nil,
	)

	user, err := service.EnsureUser(team, "ext-id")
	s.Nil(err)
	s.NotNil(user)
	s.Equal(user.TeamID, teamID.Hex())
	s.Equal(user.TeamID, teamID.Hex())
	s.Equal(user.SlackUserInfo.TZOffset, -2000)
	s.True(user.SlackUserInfo.IsAdmin)
}

func (s *UserServiceTestSuite) TestEnsureUserExisting(t *testing.T) {

	service := NewUserService(s.session)
	service.repository.save(&models.TeamUser{
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
	})

	user, err := service.EnsureUser(nil, "ext-id")
	s.Nil(err)
	s.NotNil(user)
	s.Equal(user.ExternalUserName, "ext-name")
}

func (s *UserServiceTestSuite) TestFindByID(t *testing.T) {
	service := NewUserService(s.session)
	u, err := s.repository.save(&models.TeamUser{
		ExternalUserID:   "ext-id",
		ExternalUserName: "user-name",
	})
	s.Nil(err)

	user, err := service.FindByID(u.ID.Hex())
	s.Nil(err)
	s.NotNil(user)
	s.Equal(user.ExternalUserName, "user-name")
}

func (s *UserServiceTestSuite) TestFindByIDNotExist(t *testing.T) {
	service := NewUserService(s.session)
	_, err := s.repository.save(&models.TeamUser{
		ExternalUserID:   "ext-id",
		ExternalUserName: "user-name",
	})
	s.Nil(err)

	random_id := bson.NewObjectId().Hex()
	user, err := service.FindByID(random_id)
	s.Nil(err)
	s.Nil(user)
}

type UserServiceTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	repository *UserRepository
}

func (s *UserServiceTestSuite) SetUpSuite() {
	log.Println("UserServiceTestSuite#SetUpSuite")

	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repository = NewUserRepository(session)
}

func (s *UserServiceTestSuite) TearDownSuite() {
	log.Println("UserServiceTestSuite#TearDownSuite")
	s.session.Close()
}

func (s *UserServiceTestSuite) SetUp() {
	log.Println("UserServiceTestSuite#SetUp")
	utils.TruncateTables(s.session)
}

func (s *UserServiceTestSuite) TearDown() {
	log.Println("UserServiceTestSuite#TearDown")
}
