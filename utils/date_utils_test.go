package utils

import (
	. "gopkg.in/check.v1"
	"testing"
)

// Used this to verify results: http://everytimezone.com
func (s *StringUtilsTestSuite) TestWhichTimezoneIsMidnightAt(c *C) {

	// let it be midnight in Greenwich first
	utcHour := 0
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, 0)

	// Rio (-3)
	utcHour = 3
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, -3*60*60)

	// San Francisco (-7)
	utcHour = 7
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, -7*60*60)

	// Honolulu (-10)
	utcHour = 10
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, -10*60*60)

	// Oakland (+13)
	utcHour = 11
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, 13*60*60)

	// Sydney (+10)
	utcHour = 14
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, 10*60*60)

	// Vienna (+2)
	utcHour = 22
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, 0), Equals, 2*60*60)

	// Mumbai (+5:30)
	utcHour = 18
	utcMinute := 30
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, utcMinute), Equals, 19800)

	utcHour = 21
	utcMinute = 30
	c.Assert(WhichTimezoneIsMidnightAt(utcHour, utcMinute), Equals, 10800+30*60)

}

func TestDateUtils(t *testing.T) { TestingT(t) }

type DateUtilsTestSuite struct{}

var _ = Suite(&DateUtilsTestSuite{})
