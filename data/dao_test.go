package data

import (
	"testing"

	"github.com/pavlo/slack-time/utils"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func TestDao(t *testing.T) { TestingT(t) }

type DaoTestSuite struct {
	env *utils.Environment
	dao *Dao
}

var _ = Suite(&DaoTestSuite{})

func (s *DaoTestSuite) TestFindOrCreateTeamBySlackTeamIdNew(c *C) {
	t := s.dao.FindOrCreateTeamBySlackTeamId("slack-team-id")
	c.Assert(1, Equals, s.count(Team{}))
	c.Assert(t, NotNil)
	c.Assert(t.Id, NotNil)
	c.Assert(t.SlackTeamId, Equals, "slack-team-id")
	c.Assert(t.CreatedAt, NotNil)
}

func (s *DaoTestSuite) TestFindOrCreateTeamBySlackTeamIdExisting(c *C) {
	c.Assert(0, Equals, s.count(Team{}))

	_ = s.env.OrmDB.Create(&Team{ SlackTeamId: "existing-slack-team-id"})
	c.Assert(1, Equals, s.count(Team{}))

	t := s.dao.FindOrCreateTeamBySlackTeamId("existing-slack-team-id")
	c.Assert(1, Equals, s.count(Team{}))

	c.Assert(t, NotNil)
	c.Assert(t.Id, NotNil)
	c.Assert(t.SlackTeamId, Equals, "existing-slack-team-id")
	c.Assert(t.CreatedAt, NotNil)
}

// Helper methods
func (s *DaoTestSuite) count(aType interface{}) int {
	count := 0
	s.env.OrmDB.Model(aType).Count(&count)
	return count
}

// Suite lifecycle and callbacks
func (s *DaoTestSuite) SetUpSuite(c *C) {
	e, err := utils.NewEnvironment(utils.TestEnv)
	if err != nil {
		c.Error(err)
	}

	s.env = e
	err = s.env.MigrateDatabase()
	if err != nil {
		c.Error(err)
	}

	s.env.OrmDB.LogMode(true)
	s.dao = &Dao{db: s.env.OrmDB}
}

func (s *DaoTestSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *DaoTestSuite) SetUpTest(c *C) {
	s.env.OrmDB.LogMode(true)
	s.env.OrmDB.Exec("truncate table teams cascade")
}
