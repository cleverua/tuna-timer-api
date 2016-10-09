package data

import (
	"log"
	"testing"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

func (s *PassServiceTestSuite) TestCreatePass(c *C) {

	teamID := bson.NewObjectId()
	team := &models.Team{
		ID: teamID,
	}

	userID := bson.NewObjectId()
	user := &models.TeamUser{
		ID: userID,
	}

	pass, err := s.service.CreatePass(team, user, "project-id")
	c.Assert(err, IsNil)

	c.Assert(pass.Token, Not(Equals), "")
	c.Assert(pass.ClaimedAt, IsNil)
	c.Assert(pass.ExpiresAt.Sub(pass.CreatedAt), Equals, utils.PassExpiresInMinutes * time.Minute)
	c.Assert(pass.TeamID, Equals, teamID.Hex())
	c.Assert(pass.TeamUserID, Equals, userID.Hex())
	c.Assert(pass.ProjectID, Equals, "project-id")
	c.Assert(pass.ModelVersion, Equals, models.ModelVersionPass)
}

func (s *PassServiceTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.service = NewPassService(s.session)
}

func (s *PassServiceTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *PassServiceTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestPassService(t *testing.T) { TestingT(t) }

type PassServiceTestSuite struct {
	env     *utils.Environment
	session *mgo.Session
	service *PassService
}

var _ = Suite(&PassServiceTestSuite{})
