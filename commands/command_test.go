package commands

import (
	"testing"
	"time"

	"github.com/pavlo/slack-time/data"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func TestCommand(t *testing.T) { TestingT(t) }

type TestCommandSuite struct{}

var _ = Suite(&TestCommandSuite{})

func (s *TestCommandSuite) TestUnknownCommand(c *C) {

	slackCmd := data.SlackCommand{
		ChannelID:   "channelId",
		ChannelName: "ACME",
		Command:     "timer",
		ResponseURL: "http://www.disney.com",
		TeamDomain:  "cleverua.com",
		TeamID:      "teamId",
		Text:        "unknown",
		Token:       "123e4567-e89b-12d3-a456-426655440000",
		UserID:      "userId",
		UserName:    "pavlo",
	}

	_, err := Get(slackCmd)
	c.Assert(err, NotNil)
}

func (s *TestCommandSuite) TestMarkTimerAsFinished(c *C) {
	now := time.Now()
	duration, _ := time.ParseDuration("1h25m")

	startedAt := now.Add(duration * -1)
	timer := &data.Timer{StartedAt: startedAt}
	task := &data.Task{TotalMinutes: 3}

	MarkTimerAsFinished(task, timer)
	c.Assert(timer.FinishedAt, NotNil)
	c.Assert(timer.Minutes, Equals, 60+25)
	c.Assert(task.TotalMinutes, Equals, 3+60+25)
}
