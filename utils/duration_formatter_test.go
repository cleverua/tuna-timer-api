package utils

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func (s *FormatDurationTestSuite) TestFormatting(c *C) {
	d := time.Duration(1 * time.Minute)
	c.Assert(FormatDuration(d), Equals, "00:01h")

	d = time.Duration(5*time.Hour + 25*time.Minute)
	c.Assert(FormatDuration(d), Equals, "05:25h")

	d = time.Duration(500*time.Hour + 25*time.Minute)
	c.Assert(FormatDuration(d), Equals, "500:25h")

	d = time.Duration(0)
	c.Assert(FormatDuration(d), Equals, "00:00h")
}

func TestFormatDuration(t *testing.T) { TestingT(t) }

type FormatDurationTestSuite struct{}

var _ = Suite(&FormatDurationTestSuite{})
