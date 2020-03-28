package db

import (
	"crypto/rand"
	"fmt"
	"log"
)

// Return a random string from a secure source that is sufficently long
// for use as session tokens.
func random() string {
	nbytes := 16
	b := make([]byte, nbytes)

	if _, err := rand.Read(b); err != nil {
		log.Fatal("could not generate random string:", err)
	}

	return fmt.Sprintf("%x", b)
}
