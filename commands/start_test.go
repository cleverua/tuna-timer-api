package commands

import (
	"fmt"
	"testing"
	. "gopkg.in/check.v1"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/utils"
)

// Hook up gocheck into the "go test" runner.
func TestStartCommand(t *testing.T) { TestingT(t) }

type TestStartCommandSuite struct {
	env *utils.Environment
	dao *data.Dao
}

var _ = Suite(&TestStartCommandSuite{})

func (s *TestStartCommandSuite) TestSimpleStartCommand(c *C) {
	slackCmd := data.SlackCommand{
		ChannelID: "channelId",
		ChannelName: "ACME",
		Command: "timer",
		ResponseURL: "http://www.disney.com",
		TeamDomain: "cleverua.com",
		TeamID: "teamId",
		Text: "start Convert the logotype to PNG",
		Token: "123e4567-e89b-12d3-a456-426655440000",
		UserID: "userId",
		UserName: "pavlo",
	}

	cmd, err := Get(slackCmd)
	c.Assert(err, IsNil)

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

	// Tasks
	c.Assert(1, Equals, s.env.OrmDB.Model(&project).Association("Tasks").Count())
	tasks := []data.Task{}
	s.env.OrmDB.Model(&project).Association("Tasks").Find(&tasks)
	task := tasks[0]
	c.Assert("Convert the logotype to PNG", Equals, task.Name)
	c.Assert(0, Equals, task.TotalMinutes)

	// Timers
	c.Assert(1, Equals, s.env.OrmDB.Model(&task).Association("Timers").Count())
	timers := []data.Timer{}
	s.env.OrmDB.Model(&task).Association("Timers").Find(&timers)
	timer := timers[0]
	c.Assert(user.ID, Equals, timer.TeamUserID)
	c.Assert(timer.StartedAt, NotNil)
	c.Assert(timer.FinishedAt, IsNil)
}

// let some code duplicate stay here...
func (s *TestStartCommandSuite) TestGetSimpleStartCommand(c *C) {
	slackCmd := data.SlackCommand{
		Text: "start Convert the logotype to PNG",
	}
	cmd, err := Get(slackCmd)
	c.Assert(err, IsNil)

	commandType := fmt.Sprintf("%T", cmd)
	c.Assert(commandType, Equals, "commands.Start")

	start := cmd.(Start)
	c.Assert(start.rawCommand, Equals, "Convert the logotype to PNG")
	c.Assert(start.slackCommand, NotNil)
}

func (s *TestStartCommandSuite) TestGetSimpleStartCommandWithUnicodeArgument(c *C) {
	slackCmd := data.SlackCommand{
		Text: "start Сконвертировать логотип в PNG",
	}

	cmd, err := Get(slackCmd)
	c.Assert(err, IsNil)

	commandType := fmt.Sprintf("%T", cmd)
	c.Assert(commandType, Equals, "commands.Start")

	start := cmd.(Start)
	c.Assert(start.rawCommand, Equals, "Сконвертировать логотип в PNG")
	c.Assert(start.slackCommand, NotNil)
}

// Suite lifecycle and callbacks
func (s *TestStartCommandSuite) SetUpSuite(c *C) {
	e, err := utils.NewEnvironment(utils.TestEnv)
	if err != nil {
		c.Error(err)
	}

	s.env = e
	s.dao = &data.Dao{DB: s.env.OrmDB}
	s.env.OrmDB.LogMode(true)
}

func (s *TestStartCommandSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *TestStartCommandSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.env)
}

func assertTeam(c *C, team *data.Team) {
	c.Assert(team, NotNil)
	c.Assert("teamId", Equals, team.SlackTeamID)
	c.Assert(team.ID, NotNil)
	c.Assert(team.CreatedAt, NotNil)
}