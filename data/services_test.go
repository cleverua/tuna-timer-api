package data

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/jinzhu/gorm"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
)

// ========================================================================
// FindOrCreateTeamUser tests
// ========================================================================
func (s *DataServiceTestSuite) FindOrCreateTeamUser(c *C) {

	c.Assert(0, Equals, utils.Count(s.db, models.TeamUser{}))

	team := &models.Team{}
	s.db.FirstOrCreate(team, models.Team{SlackTeamID: "a-team-id"})

	teamUser, err := s.service.FindOrCreateTeamUser(s.db, team, "user-id", "pavlo")
	c.Assert(err, IsNil)
	c.Assert(teamUser, NotNil)
	c.Assert(teamUser.ID, NotNil)
	c.Assert(teamUser.Name, Equals, "pavlo")
	c.Assert(teamUser.SlackUserID, Equals, "user-id")
	c.Assert(teamUser.TeamID, Equals, "a-team-id")
}

func TestDataService(t *testing.T) { TestingT(t) }

type DataServiceTestSuite struct {
	env     *utils.Environment
	db      *gorm.DB
	service *DataService
}

var _ = Suite(&DataServiceTestSuite{})

func (s *DataServiceTestSuite) SetUpSuite(c *C) {
	e, conn := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	e.MigrateDatabase(conn.DB())

	s.env = e
	s.db = conn
	s.service = CreateDataService()
}

func (s *DataServiceTestSuite) TearDownSuite(c *C) {
	// s.env.ReleaseResources()
}

func (s *DataServiceTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.db)
}
