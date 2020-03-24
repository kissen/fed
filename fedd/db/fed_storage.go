package db

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Represents a connection to some database that takes care of storing
// all fed related data that isn't configuration, that is user meta
// data, active session tokens and the actual activities and objects.
type FedStorage interface {
	// Open the underlying connection to the database. Before using
	// any of the other methods of a FedStorage,  call Open first.
	Open() error

	// Close the connection to the underlying database.
	Close() error

	// Retrieve the metadata for a user with the given username.
	// If no such user exists, an error is returned.
	RetrieveUser(username string) (*FedUser, error)

	// Write metadata for user. If a user with matching user.Name
	// already exists, it is overwritten.
	StoreUser(user *FedUser) error

	// Retreive metadta for given code. If no such code is recorded
	// or if it is expired, an error is returned.
	RetrieveCode(code string) (*FedOAuthCode, error)

	// Write metadata for code. If a code with matching code.Code
	// already exists, it is overwritten.
	StoreCode(code *FedOAuthCode) error

	// Retreive metadta for given token. If no such token is recorded
	// or if it is expired, an error is returned.
	RetrieveToken(token string) (*FedOAuthToken, error)

	// Write metadata for token. If a token with matching token.Token
	// already exists, it is overwritten.
	StoreToken(token *FedOAuthToken) error

	// Retrieve the object at iri.
	RetrieveObject(iri *url.URL) (vocab.Type, error)

	// Write the object at iri. If there already is an object
	// stored at the given iri, it is overwritten.
	StoreObject(iri *url.URL, obj vocab.Type) error

	// Delete the object at iri.
	DeleteObject(iri *url.URL) error
}
