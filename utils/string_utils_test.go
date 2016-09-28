package utils

import (
	"crypto/tls"
	"github.com/tuna-timer/tuna-timer-api/models"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
)

func (s *StringUtilsTestSuite) TestNormalizeSlackCustomCommand(c *C) {
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

	cmd = NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "status",
	})
	c.Assert(cmd.Text, Equals, "")
	c.Assert(cmd.SubCommand, Equals, "status")
}

func (s *StringUtilsTestSuite) TestGetSelfURLFromRequest(c *C) {
	r := &http.Request{
		Host: "subdomain.domain.com",
	}

	c.Assert(GetSelfURLFromRequest(r), Equals, "http://subdomain.domain.com")

	r = &http.Request{
		Host: "subdomain.domain.com",
		TLS:  &tls.ConnectionState{},
	}

	c.Assert(GetSelfURLFromRequest(r), Equals, "https://subdomain.domain.com")
}

func TestStringUtils(t *testing.T) { TestingT(t) }

type StringUtilsTestSuite struct{}

var _ = Suite(&StringUtilsTestSuite{})
