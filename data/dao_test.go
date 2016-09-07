package data

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/pavlo/slack-time/utils"
)

// Hook up gocheck into the "go test" runner.
func TestDao(t *testing.T) { TestingT(t) }

type DaoTestSuite struct {
	env *utils.Environment
	dao *Dao
}

var _ = Suite(&DaoTestSuite{})

// ========================================================================
// CreateTimer tests
// ========================================================================
func (s *DaoTestSuite) TestCreateTimer(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "slack-channel-id")
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "test-user")
	task := s.dao.FindOrCreateTaskByName(team, project, "task-name")
	timer := s.dao.CreateTimer(user, task)

	c.Assert(timer, NotNil)
	c.Assert(timer.DeletedAt, IsNil)
	c.Assert(timer.FinishedAt, IsNil)
	c.Assert(timer.Minutes, Equals, 0)
	c.Assert(timer.TaskID, Equals, task.ID)
	c.Assert(timer.TeamUserID, Equals, user.ID)
}

// ========================================================================
// FindNotFinishedTimerForUser tests
// ========================================================================
func (s *DaoTestSuite) TestFindNotFinishedTimerForUserNotExists(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "test-user")
	c.Assert(s.dao.FindNotFinishedTimerForUser(user), IsNil)
}

func (s *DaoTestSuite) TestFindNotFinishedTimerForUserExists(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "slack-channel-id")
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "test-user")
	task := s.dao.FindOrCreateTaskByName(team, project, "task-name")
	_ = s.env.OrmDB.Create(&Timer{TeamUserID: user.ID, TaskID: task.ID, StartedAt: time.Now()})

	timer := s.dao.FindNotFinishedTimerForUser(user)
	c.Assert(timer, NotNil)

	c.Assert(timer.DeletedAt, IsNil)
	c.Assert(timer.FinishedAt, IsNil)
	c.Assert(timer.Minutes, Equals, 0)
	c.Assert(timer.TaskID, Equals, task.ID)
	c.Assert(timer.TeamUserID, Equals, user.ID)
}

func (s *DaoTestSuite) TestFindNotFinishedTimerForUserExistsButDeleter(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "slack-channel-id")
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "test-user")
	task := s.dao.FindOrCreateTaskByName(team, project, "task-name")
	now := time.Now()
	_ = s.env.OrmDB.Create(&Timer{TeamUserID: user.ID, TaskID: task.ID, StartedAt: time.Now(), DeletedAt: &now})

	c.Assert(s.dao.FindNotFinishedTimerForUser(user), IsNil)
}

// ========================================================================
// FindOrCreateTaskByName tests
// ========================================================================
func (s *DaoTestSuite) TestFindOrCreateTaskByNameNew(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "test-channel")

	c.Assert(0, Equals, utils.Count(s.env, Task{}))
	t := s.dao.FindOrCreateTaskByName(team, project, "my task")
	c.Assert(1, Equals, utils.Count(s.env, Task{}))

	c.Assert(t, NotNil)
	c.Assert(t.ID, NotNil)

	c.Assert(t.Hash, NotNil)
	c.Assert(8, Equals, len(*t.Hash))
	c.Assert(t.Name, Equals, "my task")
	c.Assert(t.ProjectID, Equals, project.ID)
	c.Assert(t.TeamID, Equals, team.ID)
}

func (s *DaoTestSuite) TestFindOrCreateTaskByNameExisting(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "test-channel")

	_ = s.env.OrmDB.Create(&Task{ProjectID: project.ID, TeamID: team.ID, Name: "my task"})
	c.Assert(1, Equals, utils.Count(s.env, Task{}))

	t := s.dao.FindOrCreateTaskByName(team, project, "my task")
	c.Assert(1, Equals, utils.Count(s.env, Task{}))

	c.Assert(t, NotNil)
	c.Assert(t.ID, NotNil)
	c.Assert(t.Name, Equals, "my task")
	c.Assert(t.ProjectID, Equals, project.ID)
	c.Assert(t.TeamID, Equals, team.ID)
}

// ========================================================================
// FindOrCreateTeamBySlackTeamID tests
// ========================================================================
func (s *DaoTestSuite) TestFindOrCreateTeamBySlackTeamIDNew(c *C) {
	c.Assert(0, Equals, utils.Count(s.env, Team{}))
	t := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	c.Assert(1, Equals, utils.Count(s.env, Team{}))

	c.Assert(t, NotNil)
	c.Assert(t.ID, NotNil)
	c.Assert(t.SlackTeamID, Equals, "slack-team-id")
	c.Assert(t.CreatedAt, NotNil)
}

func (s *DaoTestSuite) TestFindOrCreateTeamBySlackTeamIDExisting(c *C) {
	c.Assert(0, Equals, utils.Count(s.env, Team{}))

	_ = s.env.OrmDB.Create(&Team{SlackTeamID: "existing-slack-team-id"})
	c.Assert(1, Equals, utils.Count(s.env, Team{}))

	t := s.dao.FindOrCreateTeamBySlackTeamID("existing-slack-team-id")
	c.Assert(1, Equals, utils.Count(s.env, Team{}))

	c.Assert(t, NotNil)
	c.Assert(t.ID, NotNil)
	c.Assert(t.SlackTeamID, Equals, "existing-slack-team-id")
	c.Assert(t.CreatedAt, NotNil)
}

// ========================================================================
// FindOrCreateTeamUserBySlackUserID tests
// ========================================================================
func (s *DaoTestSuite) TestFindOrCreateTeamUserBySlackUserIDNew(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")

	c.Assert(0, Equals, utils.Count(s.env, TeamUser{}))
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "U2147483697")
	c.Assert(1, Equals, utils.Count(s.env, TeamUser{}))

	c.Assert(user, NotNil)
	c.Assert(user.ID, NotNil)
	c.Assert(user.SlackUserID, Equals, "U2147483697")

	verifyTeam := &Team{}
	s.env.OrmDB.Model(user).Related(verifyTeam)

	c.Assert(verifyTeam.ID, Equals, team.ID)
}

func (s *DaoTestSuite) TestFindOrCreateTeamUserBySlackUserIdExisting(c *C) {

	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	s.env.OrmDB.Model(&team).Association("TeamUsers").Append(&TeamUser{SlackUserID: "U2147483697"})
	c.Assert(1, Equals, utils.Count(s.env, TeamUser{}))

	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "U2147483697")
	c.Assert(1, Equals, utils.Count(s.env, TeamUser{}))
	c.Assert(user.ID, NotNil)
	c.Assert(user.SlackUserID, Equals, "U2147483697")
}

// ========================================================================
// FindOrCreateProjectBySlackChannelId tests
// ========================================================================
func (s *DaoTestSuite) TestFindOrCreateProjectBySlackChannelIDNew(c *C) {
	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")

	c.Assert(0, Equals, utils.Count(s.env, Project{}))
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "Slack-Time")
	c.Assert(1, Equals, utils.Count(s.env, Project{}))

	c.Assert(project, NotNil)
	c.Assert(project.ID, NotNil)
	c.Assert(project.SlackChannelID, Equals, "Slack-Time")

	verifyTeam := &Team{}
	s.env.OrmDB.Model(project).Related(verifyTeam)

	c.Assert(verifyTeam.ID, Equals, team.ID)
}

func (s *DaoTestSuite) TestFindOrCreateProjectBySlackChannelIDExisting(c *C) {

	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")

	s.env.OrmDB.Model(&team).Association("Projects").Append(&Project{SlackChannelID: "Slack-Time"})
	c.Assert(1, Equals, utils.Count(s.env, Project{}))

	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "Slack-Time")
	c.Assert(1, Equals, utils.Count(s.env, Project{}))
	c.Assert(project.ID, NotNil)
	c.Assert(project.SlackChannelID, Equals, "Slack-Time")
}

// Suite lifecycle and callbacks
func (s *DaoTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	e.MigrateDatabase()

	s.env = e
	s.dao = &Dao{DB: s.env.OrmDB}
}

func (s *DaoTestSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *DaoTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.env)
}
