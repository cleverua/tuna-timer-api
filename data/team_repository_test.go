package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestTeamRepository(t *testing.T) {
	gosuite.Run(t, &TeamRepositoryTestSuite{Is: is.New(t)})
}

func (s *TeamRepositoryTestSuite) TestAddProject(t *testing.T) {
	team, err := s.repository.CreateTeam("external-id", "external-name")
	s.Nil(err)
	s.NotNil(team)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	s.Nil(err)

	reloadedTeam, _ := s.repository.FindByExternalID("external-id")
	s.Equal(1, len(reloadedTeam.Projects))
	testProject := reloadedTeam.Projects[0]

	s.NotNil(testProject.ID)
	s.NotNil(testProject.CreatedAt) //todo this does not make sense?
	s.Equal("external-project-id", testProject.ExternalProjectID)
	s.Equal("external-project-name", testProject.ExternalProjectName)
}

func (s *TeamRepositoryTestSuite) TestAddProjectExists(t *testing.T) {
	team, err := s.repository.CreateTeam("external-id", "external-name")
	s.Nil(err)
	s.NotNil(team)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	s.Nil(err)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	s.Nil(err)

	reloadedTeam, _ := s.repository.FindByExternalID("external-id")
	s.Equal(1, len(reloadedTeam.Projects))
}

// Find By External ID
func (s *TeamRepositoryTestSuite) TestFindByExternalID(t *testing.T) {
	team, err := s.repository.CreateTeam("external-id", "external-name")

	s.Nil(err)
	s.NotNil(team)

	resultTeam, err := s.repository.FindByExternalID("external-id")
	s.Nil(err)

	s.NotNil(resultTeam)
	s.Equal(team.ID, resultTeam.ID)
}

func (s *TeamRepositoryTestSuite) TestFindByExternalIDNotExist(t *testing.T) {
	resultTeam, err := s.repository.FindByExternalID("external-id")
	s.Nil(err)
	s.Nil(resultTeam)
}

// CREATE TEAM
func (s *TeamRepositoryTestSuite) TestCreateTeam(t *testing.T) {
	team, err := s.repository.CreateTeam("external-id", "external-name")
	s.Nil(err)
	s.NotNil(team)
	s.NotNil(team.ID)
	s.Equal("external-id", team.ExternalTeamID)
	s.Equal("external-name", team.ExternalTeamName)
	s.NotNil(team.CreatedAt) // todo - this can not be nil ever, check type rather
	s.Equal(0, len(team.Projects))
	s.Equal(models.ModelVersionTeam, team.ModelVersion)
}

func (s *TeamRepositoryTestSuite) TestCreateTeamWhenAlreadyExists(t *testing.T) {
	_, err := s.repository.CreateTeam("external-id", "external-name")
	s.Nil(err)
	_, err = s.repository.CreateTeam("external-id", "external-name")
	s.NotNil(err)
	s.True(mgo.IsDup(err))
}

func (s *TeamRepositoryTestSuite) TestSave(t *testing.T) {
	tt := &models.Team{
		ExternalTeamID:   "team-id",
		ExternalTeamName: "team-name",
		SlackOAuth:       nil,
	}
	err := s.repository.save(tt)
	s.Nil(err)
	team, err := s.repository.FindByExternalID("team-id")
	s.NotNil(team)
}

func (s *TeamRepositoryTestSuite) TestSaveUpdatesExisting(t *testing.T) {

	team, err := s.repository.CreateTeam("external-id", "external-name")
	s.Nil(err)
	s.NotNil(team)

	team.ExternalTeamName = "new-name"
	err = s.repository.save(team)
	s.Nil(err)

	tt, err := s.repository.FindByExternalID("external-id")
	s.Nil(err)
	s.Equal("new-name", tt.ExternalTeamName)
}

type TeamRepositoryTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	repository *TeamRepository
}

func (s *TeamRepositoryTestSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repository = NewTeamRepository(s.session)
}

func (s *TeamRepositoryTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *TeamRepositoryTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *TeamRepositoryTestSuite) TearDown() {}

