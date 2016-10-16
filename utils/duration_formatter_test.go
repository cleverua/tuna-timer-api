package utils

import (
	"testing"
	"time"
	"gopkg.in/tylerb/is.v1"
)

func TestDurationFormatter(t *testing.T) {
	s := is.New(t)

	d := time.Duration(1 * time.Minute)
	s.Equal(FormatDuration(d), "0:01")

	d = time.Duration(5*time.Hour + 25*time.Minute)
	s.Equal(FormatDuration(d), "5:25")

	d = time.Duration(500*time.Hour + 25*time.Minute)
	s.Equal(FormatDuration(d), "500:25")

	d = time.Duration(0)
	s.Equal(FormatDuration(d), "0:00")
}
