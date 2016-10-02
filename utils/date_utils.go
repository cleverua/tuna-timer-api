package utils

// WhichTimezoneIsMidnightAt - based up on UTC hour it returns an offset (in seconds) for a timezone where midnight (00:00) is
func WhichTimezoneIsMidnightAt(utcHour, utcMinute int) int {

	// we're handling just one corner case here for Mumbai TZ which is +5:30
	if utcHour == 18 && utcMinute == 30 {
		return (5*60 + 30) * 60
	}

	midnightAtOffset := utcHour * -1 * 60 * 60
	if utcHour >= 11 {
		midnightAtOffset = (24 - utcHour) * 60 * 60
	}
	return midnightAtOffset
}
