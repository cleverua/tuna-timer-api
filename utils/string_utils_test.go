package utils

import (
	"github.com/pavlo/slack-time/models"
	. "gopkg.in/check.v1"
	"testing"
)

func (s *FormatDurationTestSuite) TestNormalizeSlackCustomCommand(c *C) {
	cmd := NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "start Add MongoDB service to docker-compose.yml",
	})
	c.Assert(cmd.Text, Equals, "Add MongoDB service to docker-compose.yml")
	c.Assert(cmd.SubCommand, Equals, "start")

	cmd = NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "      start         Add MongoDB service to docker-compose.yml",
	})
	c.Assert(cmd.Text, Equals, "Add MongoDB service to docker-compose.yml")
	c.Assert(cmd.SubCommand, Equals, "start")
}

func TestStringUtils(t *testing.T) { TestingT(t) }

type StringUtilsTestSuite struct{}

var _ = Suite(&StringUtilsTestSuite{})
