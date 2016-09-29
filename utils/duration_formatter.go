package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats duration to be like HH:MMh regardless of days, months etc
func FormatDuration(d time.Duration) string {
	if d.Minutes() == 0 {
		return "0:00"
	}

	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	return fmt.Sprintf("%.1d:%.2d", hours, minutes)
}
