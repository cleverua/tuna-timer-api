package data

import (
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"testing"

	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	. "gopkg.in/check.v1"
)

func (s *UserRepositoryTestSuite) TestFindByExternalID(c *C) {
	user := &models.TeamUser{
		TeamID:           "team-id",
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
		SlackUserInfo: &slack.User{
			IsAdmin: true,
		},
	}

	u, err := s.repository.save(user)
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)

	loadedUser, err := s.repository.FindByExternalID("ext-id")
	c.Assert(err, IsNil)
	c.Assert(loadedUser, NotNil)

	c.Assert(loadedUser.ExternalUserID, Equals, "ext-id")
	c.Assert(loadedUser.ExternalUserName, Equals, "ext-name")
	c.Assert(loadedUser.SlackUserInfo.IsAdmin, Equals, true)
}

func (s *UserRepositoryTestSuite) TestSave(c *C) {
	user := &models.TeamUser{
		TeamID:           "team-id",
		ExternalUserID:   "ext-id",
		ExternalUserName: "ext-name",
		SlackUserInfo: &slack.User{
			IsAdmin: true,
		},
	}

	u, err := s.repository.save(user)
	c.Assert(err, IsNil)

	u.SlackUserInfo.IsAdmin = false
	_, err = s.repository.save(u)
	c.Assert(err, IsNil)

	loadedUser, err := s.repository.FindByExternalID("ext-id")
	c.Assert(err, IsNil)
	c.Assert(loadedUser, NotNil)

	c.Assert(loadedUser.ExternalUserName, Equals, "ext-name")
	c.Assert(loadedUser.SlackUserInfo.IsAdmin, Equals, false)
}

func (s *UserRepositoryTestSuite) TestFindByExternalIDNotExist(c *C) {
	resultTeam, err := s.repository.FindByExternalID("external-id")
	c.Assert(err, IsNil)
	c.Assert(resultTeam, IsNil)
}

func (s *UserRepositoryTestSuite) SetUpSuite(c *C) {
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

func (s *UserRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *UserRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestUserRepository(t *testing.T) { TestingT(t) }

type UserRepositoryTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	repository *UserRepository
}

var _ = Suite(&UserRepositoryTestSuite{})
