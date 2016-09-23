package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats duration to be like HH:MMh regardless of days, months etc
func FormatDuration(d time.Duration) string {
	if d.Minutes() == 0 {
		return "00:01h"
	}

	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	return fmt.Sprintf("%.2d:%.2dh", hours, minutes)
}
