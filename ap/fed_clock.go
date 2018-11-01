package ap

import (
	"time"
	"log"
)

// Implements the go-fed/activity/pub/Clock interface
type FedClock struct{}

// Return current UTC time
func (f *FedClock) Now() time.Time {
	t := time.Now()
	return t.UTC()
}
