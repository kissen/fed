package main

import "log"

// Panic with err if it is not nil. If err is nil,
// Must does nothing.
func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
