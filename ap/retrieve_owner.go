package ap

import (
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fediri"
)

// Return the owner (username) of this IRI.
func retrieveOwner(iri *fediri.IRI, from db.FedStorage) (*db.FedUser, error) {
	// the owner of an IRI, in the easy case, is the first
	// path component; we do not support getting the owner
	// of object IRIs yet

	if username, err := iri.Owner(); err != nil {
		return nil, err
	} else {
		return from.RetrieveUser(username)
	}
}
