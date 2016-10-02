package data

import (
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
)

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

func (s *UserServiceTestSuite) TestEnsureUserNew(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(user, NotNil)

	c.Assert(user.TeamID, Equals, teamID.Hex())
	c.Assert(user.ExternalUserName, Equals, "test-user")
	c.Assert(user.SlackUserInfo.TZOffset, Equals, -2000)
	c.Assert(user.SlackUserInfo.IsAdmin, Equals, true)
}

func (s *UserServiceTestSuite) TestEnsureUserExisting(c *C) {

	service := NewUserService(s.session)
	service.repository.save(&models.TeamUser{
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
	})

	user, err := service.EnsureUser(nil, "ext-id")

	c.Assert(err, IsNil)
	c.Assert(user, NotNil)
	c.Assert(user.ExternalUserName, Equals, "ext-name")
}

func (s *UserServiceTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	//s.service = NewUserService(s.session)
	s.repository = NewUserRepository(session)
}

func (s *UserServiceTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *UserServiceTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestUserService(t *testing.T) { TestingT(t) }

type UserServiceTestSuite struct {
	env     *utils.Environment
	session *mgo.Session
	//service    *UserService
	repository *UserRepository
}

var _ = Suite(&UserServiceTestSuite{})
