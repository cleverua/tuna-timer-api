package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestPassService(t *testing.T) {
	gosuite.Run(t, &PassServiceTestSuite{Is: is.New(t)})
}

func (s *PassServiceTestSuite) TestCreatePass(t *testing.T) {
	teamID := bson.NewObjectId()
	team := &models.Team{
		ID: teamID,
	}

	userID := bson.NewObjectId()
	user := &models.TeamUser{
		ID: userID,
	}

	pass, err := s.service.createPass(team, user, "project-id")
	s.Nil(err)

	s.NotEqual("", pass.Token)
	s.Nil(pass.ClaimedAt)
	s.Equal(utils.PassExpiresInMinutes*time.Minute, pass.ExpiresAt.Sub(pass.CreatedAt))
	s.Equal(teamID.Hex(), pass.TeamID)
	s.Equal(userID.Hex(), pass.TeamUserID)
	s.Equal("project-id", pass.ProjectID)
	s.Equal(models.ModelVersionPass, pass.ModelVersion)
}

func (s *PassServiceTestSuite) TestEnsurePassNewPassCase(t *testing.T) {
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
	s.Nil(err)
	s.NotNil(pass)

	s.Equal(userID.Hex(), pass.TeamUserID)
	s.Equal(projectID.Hex(), pass.ProjectID)
	s.Equal(teamID.Hex(), pass.TeamID)
	s.NotEqual("", pass.Token)
	s.Equal(utils.PassExpiresInMinutes*time.Minute, pass.ExpiresAt.Sub(pass.CreatedAt))
	s.Nil(pass.ClaimedAt)
}

func (s *PassServiceTestSuite) TestEnsurePassExistingOne(t *testing.T) {
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
	s.Nil(err)
	s.NotNil(pass)

	// let's change ExpiresAt for this timer a little bit to be able to assert it is prolonged
	pass.ExpiresAt = pass.ExpiresAt.Add(-3 * time.Minute)
	s.repository.update(pass)

	ensuredPass, err := s.service.EnsurePass(team, user, project)
	s.Nil(err)

	s.NotNil(ensuredPass)
	s.Equal(pass.ID.Hex(), ensuredPass.ID.Hex())
	s.Equal(userID.Hex(), ensuredPass.TeamUserID)
	s.Equal(projectID.Hex(), ensuredPass.ProjectID)
	s.Equal(teamID.Hex(), ensuredPass.TeamID)
	s.NotEqual("", ensuredPass.Token)
	//s.IsType(time.Time{}, ensuredPass.CreatedAt) //todo, find a way to assert this using Is package
	s.Nil(ensuredPass.ClaimedAt)
	s.True(utils.PassExpiresInMinutes*time.Minute.Seconds() - ensuredPass.ExpiresAt.Sub(time.Now()).Seconds() < 0.01)
}

func (s *PassServiceTestSuite) TestRemoveStalePasses(t *testing.T) {

	now := time.Now()

	p1 := &models.Pass{ //should be removed as its expiresAt is in the past
		ID:         bson.NewObjectId(),
		Token:      "p1token",
		CreatedAt:  now.Add(-5 * time.Minute),
		ExpiresAt:  now.Add(-3 * time.Minute),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	p2 := &models.Pass{ //should NOT be removed as its expiresAt is in the future
		ID:         bson.NewObjectId(),
		Token:      "p2token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	claimedAt := now.Add(-15 * 60 * 24 * time.Minute)
	p3 := &models.Pass{ //should be removed as it is claimedAt is in a distant past
		ID:         bson.NewObjectId(),
		Token:      "p3token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  &claimedAt,
		TeamUserID: "user-id",
	}

	err := s.repository.Insert(p1)
	s.Nil(err)

	err = s.repository.Insert(p2)
	s.Nil(err)

	err = s.repository.Insert(p3)
	s.Nil(err)

	err = s.service.RemoveStalePasses()
	s.Nil(err)

	pass, err := s.repository.findByID(p1.ID.Hex())
	s.Nil(err)
	s.Nil(pass)

	pass, err = s.repository.findByID(p2.ID.Hex())
	s.Nil(err)
	s.NotNil(pass)

	pass, err = s.repository.findByID(p3.ID.Hex())
	s.Nil(err)
	s.Nil(pass)
}

type PassServiceTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	service    *PassService
	repository *PassRepository
}

func (s *PassServiceTestSuite) SetUpSuite() {
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
func (s *PassServiceTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *PassServiceTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *PassServiceTestSuite) TearDown() {}

