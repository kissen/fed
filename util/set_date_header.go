package util

import (
	"net/http"
	"time"
)

// Set the Date header on headersr h to the current time
// in ANSIC format.
func SetDateHeader(h http.Header) {
	h["Date"] = []string{
		time.Now().UTC().Format(time.ANSIC),
	}
}
