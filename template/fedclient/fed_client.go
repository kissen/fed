package fedclient

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"net/url"
	"time"
)

type FedClient interface {
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
}

type fedclient struct {
	iri          *url.URL
	inboxIRI     *url.URL
	outboxIRI    *url.URL
	likedIRI     *url.URL
	followersIRI *url.URL
}

func New(actorAddr string) (_ FedClient, err error) {
	fc := &fedclient{}

	fc.iri, err = url.Parse(actorAddr)
	if err != nil {
		return nil, errors.Wrap(err, "bad actor address")
	}

	obj, err := fetch.Fetch(fc.iri)
	if err != nil {
		return nil, err
	}

	p, ok := obj.(vocab.ActivityStreamsPerson)
	if !ok {
		return nil, fmt.Errorf("%T not a supported actor type", obj)
	}

	fc.inboxIRI, err = getIRI(p.GetActivityStreamsInbox())
	if err != nil {
		return nil, err
	}

	fc.outboxIRI, err = getIRI(p.GetActivityStreamsOutbox())
	if err != nil {
		return nil, err
	}

	fc.likedIRI, err = getIRI(p.GetActivityStreamsLiked())
	if err != nil {
		return nil, err
	}

	fc.followersIRI, err = getIRI(p.GetActivityStreamsFollowers())
	if err != nil {
		return nil, err
	}

	return fc, nil
}

func (fc *fedclient) IRI() *url.URL {
	return fc.iri
}

func (fc *fedclient) Stream() (fetch.Iter, error) {
	in, err := fc.Inbox()
	if err != nil {
		return nil, err
	}

	out, err := fc.Outbox()
	if err != nil {
		return nil, err
	}

	return fetch.Begins(in, out)
}

func (fc *fedclient) Inbox() (fetch.Iter, error) {
	return fc.fetchCollection(fc.inboxIRI)
}

func (fc *fedclient) InboxIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Outbox() (fetch.Iter, error) {
	return fc.fetchCollection(fc.outboxIRI)
}

func (fc *fedclient) OutboxIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Liked() (fetch.Iter, error) {
	return fc.fetchCollection(fc.likedIRI)
}

func (fc *fedclient) LikedIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Followers() (fetch.Iter, error) {
	return fc.fetchCollection(fc.followersIRI)
}

func (fc *fedclient) FollowersIRI() *url.URL {
	return fc.followersIRI
}

func (fc *fedclient) Create(event vocab.Type) error {
	// check whether event is a supported type

	note, ok := event.(vocab.ActivityStreamsNote)
	if !ok {
		return fmt.Errorf("event of type=%T cannot be wrapped in Create", event)
	}

	// get the published date for the Create; if none is available,
	// use the current time

	publishedDate, err := prop.Published(event)
	if err != nil {
		publishedDate = time.Now()
	}

	// build up the Create

	create := streams.NewActivityStreamsCreate()

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsNote(note)
	create.SetActivityStreamsObject(object)

	published := streams.NewActivityStreamsPublishedProperty()
	published.Set(publishedDate)
	create.SetActivityStreamsPublished(published)

	audience := streams.NewActivityStreamsAudienceProperty()
	audience.AppendIRI(fc.followersIRI)
	create.SetActivityStreamsAudience(audience)

	// send it out

	return fetch.Submit(create, fc.outboxIRI)
}

func (fc *fedclient) fetchCollection(target *url.URL) (fetch.Iter, error) {
	collection, err := fetch.Fetch(target)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch collection")
	}

	return fetch.Begin(collection)
}

func getIRI(ie fetch.IterEntry) (*url.URL, error) {
	if !ie.HasAny() || !ie.IsIRI() {
		return nil, errors.New("not an IRI")
	} else {
		return ie.GetIRI(), nil
	}
}
