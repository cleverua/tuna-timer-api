package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"log"
	"testing"
	"gopkg.in/tylerb/is.v1"
	"gopkg.in/mgo.v2"
	"github.com/pavlo/gosuite"
	"github.com/nlopes/slack"
	"errors"
)

func TestTeamService(t *testing.T) {
	gosuite.Run(t, &TeamServiceTestSuite{Is: is.New(t)})
}

func (s *TeamServiceTestSuite) TestEnsureTeamNoTeamExist(t *testing.T) {
	cmd := getSlackCustomCommand()
	_, _, err := s.service.EnsureTeamSetUp(cmd)
	s.NotNil(err)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExists(t *testing.T) {
	cmd := getSlackCustomCommand()
	existingTeam, err := s.repository.CreateTeam("team-id", "team-domain")
	s.Nil(err)

	team, project, err := s.service.EnsureTeamSetUp(cmd)
	s.Nil(err)

	s.assertTeam(team)
	s.assertProject(project)

	if err != nil {
		return
	}

	s.Equal(existingTeam.ID, team.ID)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExistsWhenTeamAndUserAndProjectExist(t *testing.T) {
	cmd := getSlackCustomCommand()

	existingTeam, err := s.repository.CreateTeam("team-id", "team-domain")
	//c.Assert(err, IsNil)
	s.Nil(err)

	err = s.repository.AddProject(existingTeam, "channel-id", "channel-name")
	s.Nil(err)
	//c.Assert(err, IsNil)

	team, project, err := s.service.EnsureTeamSetUp(cmd)
	s.Nil(err)
	//c.Assert(err, IsNil)

	s.assertTeam(team)
	s.assertProject(project)

	//c.Assert(existingTeam.ID, Equals, team.ID)
	s.Equal(team.ID, existingTeam.ID)
}

func (s *TeamServiceTestSuite) TestEnsureTeamExistsFailureOnFindTeam(t *testing.T) {
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
	s.NotNil(err)
	//c.Assert(err, NotNil)

	// - Create team failure case
	modifiedRepository.findByExternalIDSuccess = true
	modifiedRepository.createTeamSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	s.NotNil(err)
	//c.Assert(err, NotNil)

	// - Add project failure case
	modifiedRepository.createTeamSuccess = true
	modifiedRepository.addProjectSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	//c.Assert(err, NotNil)
	s.NotNil(err)

	// - Add user failure case
	modifiedRepository.addProjectSuccess = true
	modifiedRepository.addUserSuccess = false
	_, _, err = s.service.EnsureTeamSetUp(cmd)
	//c.Assert(err, NotNil)
	s.NotNil(err)
}

func (s *TeamServiceTestSuite) TestCreateOrUpdateWithSlackOAuthResponseNew(t *testing.T) {
	oauthResponse := &slack.OAuthResponse{
		TeamID:      "ext-id",
		TeamName:    "ext-name",
		AccessToken: "access-token",
		Scope:       "scope",
	}

	err := s.service.CreateOrUpdateWithSlackOAuthResponse(oauthResponse)
	s.Nil(err)
	//c.Assert(err, IsNil)

	team, err := s.repository.FindByExternalID("ext-id")
	//c.Assert(err, IsNil)
	s.Nil(err)
	//c.Assert(team, NotNil)
	s.NotNil(team)

	//c.Assert(team.ExternalTeamName, Equals, "ext-name")
	s.Equal(team.ExternalTeamName, "ext-name")

	details := team.SlackOAuth
	s.NotNil(details)

	//c.Assert(details, NotNil)
	//c.Assert(details.AccessToken, Equals, "access-token")
	s.Equal(details.AccessToken, "access-token")
	//c.Assert(details.Scope, Equals, "scope")
	s.Equal(details.Scope, "scope")
}

func (s *TeamServiceTestSuite) TestCreateOrUpdateWithSlackOAuthResponseExisting(t *testing.T) {
	_, err := s.repository.CreateTeam("ext-id", "ext-name")
	s.Nil(err)
	//c.Assert(err, IsNil)

	oauthResponse := &slack.OAuthResponse{
		TeamID:      "ext-id",
		TeamName:    "ext-name-changed",
		AccessToken: "access-token",
		Scope:       "scope",
	}

	err = s.service.CreateOrUpdateWithSlackOAuthResponse(oauthResponse)
	s.Nil(err)
	//c.Assert(err, IsNil)

	team, err := s.repository.FindByExternalID("ext-id")
	s.Nil(err)
	//c.Assert(err, IsNil)
	s.NotNil(team)
	//c.Assert(team, NotNil)

	//c.Assert(team.ExternalTeamName, Equals, "ext-name-changed")
	s.Equal(team.ExternalTeamName, "ext-name-changed")

	details := team.SlackOAuth
	//c.Assert(details, NotNil)
	s.NotNil(details)

	//c.Assert(details.AccessToken, Equals, "access-token")
	s.Equal(details.AccessToken, "access-token")

	//c.Assert(details.Scope, Equals, "scope")
	s.Equal(details.Scope, "scope")
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

func (s *TeamServiceTestSuite) assertTeam(team *models.Team) {
	//c.Assert(team, NotNil)
	s.NotNil(team)
	//c.Assert(team.ID, NotNil)
	s.NotNil(team.ID)

	//c.Assert(team.ExternalTeamID, Equals, "team-id")
	s.Equal("team-id", team.ExternalTeamID)

	//c.Assert(team.ExternalTeamName, Equals, "team-domain")
	s.Equal("team-domain", team.ExternalTeamName)

	//c.Assert(team.CreatedAt, NotNil)
	s.NotNil(team.CreatedAt) // todo change to type checking

	//c.Assert(len(team.Projects), Equals, 1)
	s.Equal(1, len(team.Projects))

	s.assertProject(team.Projects[0])
}

func (s *TeamServiceTestSuite) assertProject(project *models.Project) {
	s.NotNil(project)
	s.NotNil(project.ID)
	s.Equal("channel-id", project.ExternalProjectID)
	s.Equal("channel-name", project.ExternalProjectName)
	s.NotNil(project.CreatedAt) //todo - check type rather
}


// testTeamRepositoryImpl allows is a TeamRepositoryInterface that is able to simulate returned errors
type testTeamRepositoryImpl struct {
	repository              TeamRepositoryInterface
	findByExternalIDSuccess bool
	findByIDSuccess		bool
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

func (r *testTeamRepositoryImpl) FindByID(teamID string) (*models.Team, error) {
	if !r.findByIDSuccess {
		return nil, errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.FindByID(teamID)
}

func (r *testTeamRepositoryImpl) CreateTeam(externalID, externalName string) (*models.Team, error) {
	if !r.createTeamSuccess {
		return nil, errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.CreateTeam(externalID, externalName)
}

func (r *testTeamRepositoryImpl) AddProject(team *models.Team, externalProjectID, externalProjectName string) error {
	if !r.addProjectSuccess {
		return errors.New("TestTeamRepositoryImpl error")
	}
	return r.repository.AddProject(team, externalProjectID, externalProjectName)
}

func (r *testTeamRepositoryImpl) save(team *models.Team) error {
	return nil
}

type TeamServiceTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	service    *TeamService
	repository TeamRepositoryInterface
}

func (s *TeamServiceTestSuite) SetUpSuite() {
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

func (s *TeamServiceTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *TeamServiceTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *TeamServiceTestSuite) TearDown() {}
