package fedcontext

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/fed/fetch"
	"net/url"
)

// A client that can be used to issue command from the point of
// view of an actor on the ActivityPub fediverse.
type FedClient interface {
	// The name of this user.
	Username() string

	// Return the clients actor IRI.
	IRI() *url.URL

	// Return the stream iterator. Stream is a combination of both
	// inbox and outbox.
	Stream() (fetch.Iter, error)

	// Return an iterator of the users inbox.
	Inbox() (fetch.Iter, error)

	// Return the IRI of the users inbox.
	InboxIRI() *url.URL

	// Return an iterator of the users outbox.
	Outbox() (fetch.Iter, error)

	// Return the IRI of the users outbox.
	OutboxIRI() *url.URL

	// Return an iterator of the objects this user liked.
	Liked() (fetch.Iter, error)

	// Return the IRI of the users liked collection.
	LikedIRI() *url.URL

	// Return an iterator of actors that follow the user.
	Followers() (fetch.Iter, error)

	// Return the IRI to the collection of actors that follow this user.
	FollowersIRI() *url.URL

	// Wrap event into an Create activity and submit
	// it to the users outbox.
	Create(event vocab.Type) error

	// Like the object at iri.
	Like(iri *url.URL) error
}
