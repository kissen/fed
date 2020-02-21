package db

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

type FedStorage interface {
	// User management.
	RetrieveUser(username string) (*FedUser, error)
	StoreUser(username string) (*FedUser, error)

	// Reading and writing objects. Objects are the base type
	// for all subtypes, e.g. Actor, Activity, Link or Collection.
	// (See https://www.w3.org/TR/activitystreams-core/#object)
	RetrieveObject(iri *url.URL) (vocab.Type, error)
	StoreObject(obj vocab.Type) (*url.URL, error)
	StoreObjectAt(iri *url.URL, obj vocab.Type) error
}
