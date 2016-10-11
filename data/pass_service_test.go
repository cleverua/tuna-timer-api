package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
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

	pass, err := s.service.createPass(team, user, "project-id")
	c.Assert(err, IsNil)

	c.Assert(pass.Token, Not(Equals), "")
	c.Assert(pass.ClaimedAt, IsNil)
	c.Assert(pass.ExpiresAt.Sub(pass.CreatedAt), Equals, utils.PassExpiresInMinutes*time.Minute)
	c.Assert(pass.TeamID, Equals, teamID.Hex())
	c.Assert(pass.TeamUserID, Equals, userID.Hex())
	c.Assert(pass.ProjectID, Equals, "project-id")
	c.Assert(pass.ModelVersion, Equals, models.ModelVersionPass)
}

func (s *PassServiceTestSuite) TestEnsurePassNewPassCase(c *C) {
	teamID := bson.NewObjectId()
	team := &models.Team{
		ID: teamID,
	}

	userID := bson.NewObjectId()
	user := &models.TeamUser{
		ID: userID,
	}

	projectID := bson.NewObjectId()
	project := &models.Project{
		ID: projectID,
	}

	pass, err := s.service.EnsurePass(team, user, project)
	c.Assert(err, IsNil)
	c.Assert(pass, NotNil)

	c.Assert(pass.TeamUserID, Equals, userID.Hex())
	c.Assert(pass.ProjectID, Equals, projectID.Hex())
	c.Assert(pass.TeamID, Equals, teamID.Hex())
	c.Assert(pass.Token, NotNil)
	c.Assert(pass.CreatedAt, NotNil)
	c.Assert(pass.ExpiresAt.Sub(pass.CreatedAt), Equals, utils.PassExpiresInMinutes*time.Minute)
	c.Assert(pass.ClaimedAt, IsNil)
}

func (s *PassServiceTestSuite) TestEnsurePassExistingOne(c *C) {
	teamID := bson.NewObjectId()
	team := &models.Team{
		ID: teamID,
	}

	userID := bson.NewObjectId()
	user := &models.TeamUser{
		ID: userID,
	}

	projectID := bson.NewObjectId()
	project := &models.Project{
		ID: projectID,
	}

	pass, err := s.service.EnsurePass(team, user, project)
	c.Assert(err, IsNil)
	c.Assert(pass, NotNil)

	// let's change ExpiresAt for this timer a little bit to be able to assert it is prolonged
	pass.ExpiresAt = pass.ExpiresAt.Add(-3 * time.Minute)
	s.repository.update(pass)

	ensuredPass, err := s.service.EnsurePass(team, user, project)
	c.Assert(err, IsNil)
	c.Assert(ensuredPass, NotNil)

	c.Assert(ensuredPass.ID.Hex(), Equals, pass.ID.Hex())
	c.Assert(ensuredPass.TeamUserID, Equals, userID.Hex())
	c.Assert(ensuredPass.ProjectID, Equals, projectID.Hex())
	c.Assert(ensuredPass.TeamID, Equals, teamID.Hex())
	c.Assert(ensuredPass.Token, NotNil)
	c.Assert(ensuredPass.CreatedAt, NotNil)
	c.Assert(pass.ClaimedAt, IsNil)

	diffSeconds := utils.PassExpiresInMinutes*time.Minute.Seconds() - ensuredPass.ExpiresAt.Sub(time.Now()).Seconds()
	isZeroSeconds := diffSeconds < 0.001
	c.Assert(isZeroSeconds, Equals, true)
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
	s.repository = NewPassRepository(session)
}

func (s *PassServiceTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *PassServiceTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestPassService(t *testing.T) { TestingT(t) }

type PassServiceTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	service    *PassService
	repository *PassRepository
}

var _ = Suite(&PassServiceTestSuite{})
