package commands

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Hook up gocheck into the "go test" runner.
func TestStatusCommand(t *testing.T) { TestingT(t) }

type TestStatusCommandSuite struct {
	env *utils.Environment
	dao *data.Dao
}

var _ = Suite(&TestStatusCommandSuite{})

func (s *TestStatusCommandSuite) TestSimpleStatusCommandNoTimerFound(c *C) {

	c.Assert(utils.Count(s.env, data.Timer{}), Equals, 0)

	slackCmd := data.SlackCommand{
		ChannelID:   "channelId",
		ChannelName: "ACME",
		Command:     "timer",
		ResponseURL: "http://www.disney.com",
		TeamDomain:  "cleverua.com",
		TeamID:      "teamId",
		Text:        "status",
		Token:       "123e4567-e89b-12d3-a456-426655440000",
		UserID:      "userId",
		UserName:    "pavlo",
	}

	cmd, err := Get(slackCmd)
	c.Assert(err, IsNil)

	cmdType := fmt.Sprintf("%T", cmd)
	c.Assert(cmdType, Equals, "commands.Status")

	result := cmd.Execute(s.env)
	c.Assert(result, NotNil)

	// Asserting team
	team := data.Team{}
	s.env.OrmDB.First(&team)
	assertTeam(c, &team)
	assertTeam(c, result.data["team"].(*data.Team))

	// Team users
	c.Assert(1, Equals, s.env.OrmDB.Model(&team).Association("TeamUsers").Count())
	users := []data.TeamUser{}
	s.env.OrmDB.Model(&team).Association("TeamUsers").Find(&users)
	user := users[0]
	c.Assert("userId", Equals, user.SlackUserID)

	// Project
	c.Assert(1, Equals, s.env.OrmDB.Model(&team).Association("Projects").Count())
	projects := []data.Project{}
	s.env.OrmDB.Model(&team).Association("Projects").Find(&projects)
	project := projects[0]
	c.Assert("channelId", Equals, project.SlackChannelID)
}

func (s *TestStatusCommandSuite) TestSimpleStopCommand(c *C) {

	team := s.dao.FindOrCreateTeamBySlackTeamID("slack-team-id")
	project := s.dao.FindOrCreateProjectBySlackChannelID(team, "slack-channel-id")
	user := s.dao.FindOrCreateTeamUserBySlackUserID(team, "test-user")
	task := s.dao.FindOrCreateTaskByName(team, project, "task-name")
	timer := s.dao.CreateTimer(user, task)

	duration, _ := time.ParseDuration("1h25m")
	startedAt := time.Now().Add(duration * -1)
	timer.StartedAt = startedAt
	s.env.OrmDB.Save(&timer)

	slackCmd := data.SlackCommand{
		ChannelID:   "slack-channel-id",
		ChannelName: "slack-channel-id",
		Command:     "timer",
		ResponseURL: "http://www.disney.com",
		TeamDomain:  "cleverua.com",
		TeamID:      "slack-team-id",
		Text:        "status",
		Token:       "123e4567-e89b-12d3-a456-426655440000",
		UserID:      "test-user",
		UserName:    "test-user",
	}

	cmd, err := Get(slackCmd)
	c.Assert(err, IsNil)

	cmdType := fmt.Sprintf("%T", cmd)
	c.Assert(cmdType, Equals, "commands.Status")

	result := cmd.Execute(s.env)
	c.Assert(result, NotNil)

	verifyTimer := result.data["timer"].(*data.Timer)
	c.Assert(verifyTimer.ID, Equals, timer.ID)
	c.Assert(verifyTimer.Minutes, Equals, 60+25)

	c.Assert(result.data["team"], NotNil)
	c.Assert(result.data["user"], NotNil)
	c.Assert(result.data["project"], NotNil)
	c.Assert(result.data["task"], NotNil)
}

// Suite lifecycle and callbacks
func (s *TestStatusCommandSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")
	e.MigrateDatabase()

	s.env = e
	s.dao = &data.Dao{DB: s.env.OrmDB}
}

func (s *TestStatusCommandSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *TestStatusCommandSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.env)
}
