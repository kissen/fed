package ap

import (
	"time"
)

// Implements the go-fed/activity/pub/Clock interface (version 1.0)
type FedClock struct{}

// Return current server time
func (f *FedClock) Now() time.Time {
	return time.Now().UTC()
}
