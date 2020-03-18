package fedclient

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"net/url"
)

type FedClient interface {
	// Return the clients actor IRI.
	IRI() *url.URL

	// Return an iterator of the users inbox.
	Inbox() (fedutil.Iter, error)

	// Return the IRI of the users inbox.
	InboxIRI() *url.URL

	// Return an iterator of the users outbox.
	Outbox() (fedutil.Iter, error)

	// Return the IRI of the users outbox.
	OutboxIRI() *url.URL

	// Return an iterator of the objects this user liked.
	Liked() (fedutil.Iter, error)

	// Return the IRI of the users liked collection.
	LikedIRI() *url.URL

	// Wrap event into an Create activity and submit
	// it to the users outbox.
	Create(event vocab.Type) error
}

type fedclient struct {
	iri       *url.URL
	inboxIRI  *url.URL
	outboxIRI *url.URL
	likedIRI  *url.URL
}

func New(actorAddr string) (_ FedClient, err error) {
	fc := &fedclient{}

	fc.iri, err = url.Parse(actorAddr)
	if err != nil {
		return nil, errors.Wrap(err, "bad actor address")
	}

	obj, err := fedutil.Fetch(fc.iri)
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

	return fc, nil
}

func (fc *fedclient) IRI() *url.URL {
	return fc.iri
}

func (fc *fedclient) Inbox() (fedutil.Iter, error) {
	return fc.fetchCollection(fc.inboxIRI)
}

func (fc *fedclient) InboxIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Outbox() (fedutil.Iter, error) {
	return fc.fetchCollection(fc.outboxIRI)
}

func (fc *fedclient) OutboxIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Liked() (fedutil.Iter, error) {
	return fc.fetchCollection(fc.likedIRI)
}

func (fc *fedclient) LikedIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedclient) Create(event vocab.Type) error {
	note, ok := event.(vocab.ActivityStreamsNote)
	if !ok {
		return fmt.Errorf("event of type=%T cannot be wrapped in Create", event)
	}

	return fedutil.Submit(fc.outboxIRI, note)
}

func (fc *fedclient) fetchCollection(target *url.URL) (fedutil.Iter, error) {
	collection, err := fedutil.Fetch(target)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch collection")
	}

	return fedutil.Begin(collection)
}

func getIRI(ie fedutil.IterEntry) (*url.URL, error) {
	if !ie.HasAny() || !ie.IsIRI() {
		return nil, errors.New("not an IRI")
	} else {
		return ie.GetIRI(), nil
	}
}
