package data

import (
	"log"
	"testing"

	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

//EnsureTeamExists
func (s *TeamServiceTestSuite) TestEnsureTeamExistsNewTeam(c *C) {
	cmd := getSlackCustomCommand()
	team, project, user, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, IsNil)

	assertTeam(c, team)
	assertProject(c, project)
	assertUser(c, user)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExists(c *C) {
	cmd := getSlackCustomCommand()
	existingTeam, err := s.repository.createTeam("team-id", "team-domain")
	c.Assert(err, IsNil)

	team, project, user, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, IsNil)

	assertTeam(c, team)
	assertProject(c, project)
	assertUser(c, user)

	c.Assert(existingTeam.ID, Equals, team.ID)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExistsWhenTeamAndUserAndProjectExist(c *C) {
	cmd := getSlackCustomCommand()

	existingTeam, err := s.repository.createTeam("team-id", "team-domain")
	c.Assert(err, IsNil)

	err = s.repository.addProject(existingTeam, "channel-id", "channel-name")
	c.Assert(err, IsNil)

	err = s.repository.addUser(existingTeam, "user-id", "user-name")
	c.Assert(err, IsNil)

	team, project, user, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, IsNil)

	assertTeam(c, team)
	assertProject(c, project)
	assertUser(c, user)

	c.Assert(existingTeam.ID, Equals, team.ID)
}

func getSlackCustomCommand() *models.SlackCustomCommand {
	return &models.SlackCustomCommand{
		ChannelID:   "channel-id",
		ChannelName: "channel-name",
		Command:     "timer",
		ResponseURL: "http://www.cleverua.com",
		TeamDomain:  "team-domain",
		TeamID:      "team-id",
		Text:        "the text of the command",
		Token:       "token",
		UserID:      "user-id",
		UserName:    "user-name",
	}
}

func assertTeam(c *C, team *models.Team) {
	c.Assert(team, NotNil)
	c.Assert(team.ID, NotNil)
	c.Assert(team.ExternalTeamID, Equals, "team-id")
	c.Assert(team.ExternalTeamName, Equals, "team-domain")
	c.Assert(team.CreatedAt, NotNil)
	c.Assert(len(team.Projects), Equals, 1)
	c.Assert(len(team.Users), Equals, 1)

	assertProject(c, team.Projects[0])
	assertUser(c, team.Users[0])
}

func assertProject(c *C, project *models.Project) {
	c.Assert(project, NotNil)
	c.Assert(project.ID, NotNil)
	c.Assert(project.ExternalProjectID, Equals, "channel-id")
	c.Assert(project.ExternalProjectName, Equals, "channel-name")
	c.Assert(project.CreatedAt, NotNil)
}

func assertUser(c *C, user *models.TeamUser) {
	c.Assert(user, NotNil)
	c.Assert(user.ID, NotNil)
	c.Assert(user.ExternalUserID, Equals, "user-id")
	c.Assert(user.ExternalUserName, Equals, "user-name")
	c.Assert(user.CreatedAt, NotNil)
}

// Suite lifecycle and callbacks
func (s *TeamServiceTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.service = NewTeamService(s.session)
	s.repository = NewTeamRepository(s.session)
}

func (s *TeamServiceTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TeamServiceTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestTeamService(t *testing.T) { TestingT(t) }

type TeamServiceTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	service    *TeamService
	repository *TeamRepository
}

var _ = Suite(&TeamServiceTestSuite{})
