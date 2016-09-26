package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
)

// Team User
func (s *TeamRepositoryTestSuite) TestAddUser(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	err = s.repository.addUser(team, "external-user-id", "external-user-name")
	c.Assert(err, IsNil)

	reloadedTeam, _ := s.repository.findByExternalID("external-id")
	c.Assert(len(reloadedTeam.Users), Equals, 1)
	testUser := reloadedTeam.Users[0]

	c.Assert(testUser.ID, NotNil)
	c.Assert(testUser.CreatedAt, NotNil)
	c.Assert(testUser.ExternalUserID, Equals, "external-user-id")
	c.Assert(testUser.ExternalUserName, Equals, "external-user-name")
}

func (s *TeamRepositoryTestSuite) TestAddUserExists(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	err = s.repository.addUser(team, "external-user-id", "external-user-name")
	c.Assert(err, IsNil)

	err = s.repository.addUser(team, "external-user-id", "external-use-name")
	c.Assert(err, IsNil)

	reloadedTeam, _ := s.repository.findByExternalID("external-id")
	c.Assert(len(reloadedTeam.Users), Equals, 1)
}

// Add Project
func (s *TeamRepositoryTestSuite) TestAddProject(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	c.Assert(err, IsNil)

	reloadedTeam, _ := s.repository.findByExternalID("external-id")
	c.Assert(len(reloadedTeam.Projects), Equals, 1)
	testProject := reloadedTeam.Projects[0]

	c.Assert(testProject.ID, NotNil)
	c.Assert(testProject.CreatedAt, NotNil)
	c.Assert(testProject.ExternalProjectID, Equals, "external-project-id")
	c.Assert(testProject.ExternalProjectName, Equals, "external-project-name")
}

func (s *TeamRepositoryTestSuite) TestAddProjectExists(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	c.Assert(err, IsNil)

	err = s.repository.addProject(team, "external-project-id", "external-project-name")
	c.Assert(err, IsNil)

	reloadedTeam, _ := s.repository.findByExternalID("external-id")
	c.Assert(len(reloadedTeam.Projects), Equals, 1)
}

// Find By External ID
func (s *TeamRepositoryTestSuite) TestFindByExternalId(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	resultTeam, err := s.repository.findByExternalID("external-id")
	c.Assert(err, IsNil)
	c.Assert(resultTeam, NotNil)
	c.Assert(resultTeam.ID, Equals, team.ID)
}

func (s *TeamRepositoryTestSuite) TestFindByExternalIdNotExist(c *C) {
	resultTeam, err := s.repository.findByExternalID("external-id")
	c.Assert(err, IsNil)
	c.Assert(resultTeam, IsNil)
}

// CREATE TEAM
func (s *TeamRepositoryTestSuite) TestCreateTeam(c *C) {
	team, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)
	c.Assert(team.ID, NotNil)
	c.Assert(team.ExternalTeamID, Equals, "external-id")
	c.Assert(team.ExternalTeamName, Equals, "external-name")
	c.Assert(team.CreatedAt, NotNil)
	c.Assert(len(team.Projects), Equals, 0)
	c.Assert(len(team.Users), Equals, 0)
}

func (s *TeamRepositoryTestSuite) TestCreateTeamWhenAlreadyExists(c *C) {
	_, err := s.repository.createTeam("external-id", "external-name")
	c.Assert(err, IsNil)

	_, err = s.repository.createTeam("external-id", "external-name")
	c.Assert(err, NotNil)
	c.Assert(mgo.IsDup(err), Equals, true)
}

// Suite lifecycle and callbacks
func (s *TeamRepositoryTestSuite) SetUpSuite(c *C) {
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

func (s *TeamRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TeamRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestTeamRepository(t *testing.T) { TestingT(t) }

type TeamRepositoryTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	repository *TeamRepository
}

var _ = Suite(&TeamRepositoryTestSuite{})
