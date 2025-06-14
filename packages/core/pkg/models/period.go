package models

import "time"

// Period represents a time period with start and end dates
type Period struct {
	Start time.Time
	End   time.Time
}
