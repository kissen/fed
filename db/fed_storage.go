package db

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

type FedStorage interface {
	// Lifetime management.
	Open() error
	Close() error

	// User management.
	RetrieveUser(username string) (*FedUser, error)
	StoreUser(user *FedUser) error

	// Reading and writing objects. Objects are the base type
	// for all subtypes, e.g. Actor, Activity, Link or Collection.
	// (See https://www.w3.org/TR/activitystreams-core/#object)
	RetrieveObject(iri *url.URL) (vocab.Type, error)
	StoreObject(iri *url.URL, obj vocab.Type) error
	DeleteObject(iri *url.URL) error
}
