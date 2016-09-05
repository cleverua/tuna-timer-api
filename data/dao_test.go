package data

import (
	"testing"
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
// FindOrCreateTeamBySlackTeamID tests
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
	e, err := utils.NewEnvironment(utils.TestEnv)
	if err != nil {
		c.Error(err)
	}

	s.env = e
	s.dao = &Dao{DB: s.env.OrmDB}
	s.env.OrmDB.LogMode(true)
}

func (s *DaoTestSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *DaoTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.env)
}
