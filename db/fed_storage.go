package db

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Defines operations on a database required by fed.
type Storer interface {
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

// Represents a connection to some database that takes care of storing
// all fed related data that isn't configuration, that is user meta
// data, active session tokens and the actual activities and objects.
//
// FedStorage provides the Begin method which returns a transaction.
// A transaction combines any number of operations (as defined by Storer)
// and allows us to apply them atomically.
//
// If you just want to run a single operation, you can also call the Storer
// methods directly.
type FedStorage interface {
	// Open the underlying connection to the database. Before using
	// any of the other methods of a FedStorage,  call Open first.
	Open() error

	// Close the connection to the underlying database.
	Close() error

	// Start a new transaction. Remember to call Rollback or Commit!
	Begin() (Tx, error)

	// FedStorage implements Storer. Calling its methods creates
	// a single-operation transaction and automatically commits.
	Storer
}

// A single transaction. Create one if you do operations on the database
// you may want to revert. You absolutly musn't forget to call either
// Commit or Rollback, otherwise fed will deadlock.
type Tx interface {
	// Commit all changes made within this transaction. You can
	// call this method as often as you want, only the first call
	// to Commit or Rollback for a given instance will be applied.
	Commit() error

	// Undo all changes made by this transaction. You can call
	// this method as often as you want, only the first call to Commit
	// or Rollback for a given instance will be applied.
	Rollback() error

	// Tx implements Storer. Calling its methods is only applied
	// after a successful call to Commit. If you wish to undo
	// the changes, call Rollback instead.
	Storer
}
