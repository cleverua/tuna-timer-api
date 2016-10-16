package utils

import (
	//. "gopkg.in/check.v1"
	"testing"
	"gopkg.in/tylerb/is.v1"
)

// Used this to verify results: http://everytimezone.com
func TestWhichTimezoneIsMidnightAt(t *testing.T) {

	s := is.New(t)
	// let it be midnight in Greenwich first
	utcHour := 0
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), 0)

	// Rio (-3)
	utcHour = 3
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), -3*60*60)

	// San Francisco (-7)
	utcHour = 7
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), -7*60*60)

	// Honolulu (-10)
	utcHour = 10
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), -10*60*60)

	// Oakland (+13)
	utcHour = 11
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), 13*60*60)

	// Sydney (+10)
	utcHour = 14
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), 10*60*60)

	// Vienna (+2)
	utcHour = 22
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, 0), 2*60*60)

	// Mumbai (+5:30)
	utcHour = 18
	utcMinute := 30
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, utcMinute), 19800)

	utcHour = 21
	utcMinute = 30
	s.Equal(WhichTimezoneIsMidnightAt(utcHour, utcMinute), 10800+30*60)

}

