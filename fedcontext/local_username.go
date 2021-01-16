package fedcontext

import "github.com/kissen/fed/fediri"

// If fc is a client for a user on this instance, return the
// username of that user.
func LocalUsername(fc FedClient) (username string, ok bool) {
	iri := fediri.IRI{fc.IRI()}

	if username, err := iri.Actor(); err != nil {
		return "", false
	} else {
		return username, true
	}
}
