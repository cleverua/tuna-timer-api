package data

import (
	"log"
	"testing"

	"errors"

	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

func (s *TeamServiceTestSuite) TestEnsureTeamNoTeamExist(c *C) {
	cmd := getSlackCustomCommand()
	_, _, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, NotNil)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExists(c *C) {
	cmd := getSlackCustomCommand()
	existingTeam, err := s.repository.createTeam("team-id", "team-domain")
	c.Assert(err, IsNil)

	team, project, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, IsNil)

	assertTeam(c, team)
	assertProject(c, project)

	c.Assert(existingTeam.ID, Equals, team.ID)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExistsWhenTeamAndUserAndProjectExist(c *C) {
	cmd := getSlackCustomCommand()

	existingTeam, err := s.repository.createTeam("team-id", "team-domain")
	c.Assert(err, IsNil)

	err = s.repository.addProject(existingTeam, "channel-id", "channel-name")
	c.Assert(err, IsNil)

	team, project, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, IsNil)

	assertTeam(c, team)
	assertProject(c, project)

	c.Assert(existingTeam.ID, Equals, team.ID)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExistsFailureOnFindTeam(c *C) {
	modifiedRepository := &testTeamRepositoryImpl{
		findByExternalIDSuccess: false,
		createTeamSuccess:       true,
		addProjectSuccess:       true,
		addUserSuccess:          true,
		repository:              s.repository,
	}

	s.service.repository = modifiedRepository
	defer func() {
		s.service.repository = s.repository
	}()

	cmd := getSlackCustomCommand()

	// - FindTeam failure case
	_, _, err := s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, NotNil)

	// - Create team failure case
	modifiedRepository.findByExternalIDSuccess = true
	modifiedRepository.createTeamSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, NotNil)

	// - Add project failure case
	modifiedRepository.createTeamSuccess = true
	modifiedRepository.addProjectSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, NotNil)

	// - Add user failure case
	modifiedRepository.addProjectSuccess = true
	modifiedRepository.addUserSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	c.Assert(err, NotNil)
}

func (r *TeamServiceTestSuite) TestCreateOrUpdateWithSlackOAuthResponseNew(c *C) {
	oauthResponse := &slack.OAuthResponse{
		TeamID:      "ext-id",
		TeamName:    "ext-name",
		AccessToken: "access-token",
		Scope:       "scope",
	}

	err := r.service.CreateOrUpdateWithSlackOAuthResponse(oauthResponse)
	c.Assert(err, IsNil)

	team, err := r.repository.FindByExternalID("ext-id")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	c.Assert(team.ExternalTeamName, Equals, "ext-name")

	details := team.SlackOAuth
	c.Assert(details, NotNil)
	c.Assert(details.AccessToken, Equals, "access-token")
	c.Assert(details.Scope, Equals, "scope")
}

func (r *TeamServiceTestSuite) TestCreateOrUpdateWithSlackOAuthResponseExisting(c *C) {
	_, err := r.repository.createTeam("ext-id", "ext-name")
	c.Assert(err, IsNil)

	oauthResponse := &slack.OAuthResponse{
		TeamID:      "ext-id",
		TeamName:    "ext-name-changed",
		AccessToken: "access-token",
		Scope:       "scope",
	}

	err = r.service.CreateOrUpdateWithSlackOAuthResponse(oauthResponse)
	c.Assert(err, IsNil)

	team, err := r.repository.FindByExternalID("ext-id")
	c.Assert(err, IsNil)
	c.Assert(team, NotNil)

	c.Assert(team.ExternalTeamName, Equals, "ext-name-changed")

	details := team.SlackOAuth
	c.Assert(details, NotNil)
	c.Assert(details.AccessToken, Equals, "access-token")
	c.Assert(details.Scope, Equals, "scope")
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

	assertProject(c, team.Projects[0])
}

func assertProject(c *C, project *models.Project) {
	c.Assert(project, NotNil)
	c.Assert(project.ID, NotNil)
	c.Assert(project.ExternalProjectID, Equals, "channel-id")
	c.Assert(project.ExternalProjectName, Equals, "channel-name")
	c.Assert(project.CreatedAt, NotNil)
}

// testTeamRepositoryImpl allows is a TeamRepositoryInterface that is able to simulate returned errors
type testTeamRepositoryImpl struct {
	repository              TeamRepositoryInterface
	findByExternalIDSuccess bool
	createTeamSuccess       bool
	addProjectSuccess       bool
	addUserSuccess          bool
}

func (r *testTeamRepositoryImpl) FindByExternalID(externalTeamID string) (*models.Team, error) {
	if !r.findByExternalIDSuccess {
		return nil, errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.FindByExternalID(externalTeamID)
}

func (r *testTeamRepositoryImpl) createTeam(externalID, externalName string) (*models.Team, error) {
	if !r.createTeamSuccess {
		return nil, errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.createTeam(externalID, externalName)
}

func (r *testTeamRepositoryImpl) addProject(team *models.Team, externalProjectID, externalProjectName string) error {
	if !r.addProjectSuccess {
		return errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.addProject(team, externalProjectID, externalProjectName)
}

func (r *testTeamRepositoryImpl) save(team *models.Team) error {
	return nil
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
	repository TeamRepositoryInterface
}

var _ = Suite(&TeamServiceTestSuite{})
