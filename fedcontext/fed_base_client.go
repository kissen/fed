package fedcontext

import (
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"net/url"
)

// Implements FedClient.
type fedbaseclient struct {
	// IRI pointing to the actor that "owns" this client.
	iri *url.URL

	// IRI pointing to the inbox of the owner.
	inboxIRI *url.URL

	// IRI pointing to the outbox of the owner.
	outboxIRI *url.URL

	// IRI pointing to the liked collection of the owner.
	likedIRI *url.URL

	// IRI pointing to the followers collection of the owner.
	followersIRI *url.URL

	// Function that gets invoked on Create calls.
	create func(vocab.Type) error

	// Function that gets invoked on Like calls.
	like func(*url.URL) error
}

func (fc *fedbaseclient) fill(actorAddr string) error {
	var err error

	fc.iri, err = url.Parse(actorAddr)
	if err != nil {
		return errors.Wrap(err, "bad actor address")
	}

	obj, err := fetch.Fetch(fc.iri)
	if err != nil {
		return err
	}

	p, ok := obj.(vocab.ActivityStreamsPerson)
	if !ok {
		return fmt.Errorf("%T not a supported actor type", obj)
	}

	fc.inboxIRI, err = getIRI(p.GetActivityStreamsInbox())
	if err != nil {
		return err
	}

	fc.outboxIRI, err = getIRI(p.GetActivityStreamsOutbox())
	if err != nil {
		return err
	}

	fc.likedIRI, err = getIRI(p.GetActivityStreamsLiked())
	if err != nil {
		return err
	}

	fc.followersIRI, err = getIRI(p.GetActivityStreamsFollowers())
	if err != nil {
		return err
	}

	return nil
}

func (fc *fedbaseclient) IRI() *url.URL {
	return fc.iri
}

func (fc *fedbaseclient) Stream() (fetch.Iter, error) {
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

func (fc *fedbaseclient) Inbox() (fetch.Iter, error) {
	return fc.fetchCollection(fc.inboxIRI)
}

func (fc *fedbaseclient) InboxIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedbaseclient) Outbox() (fetch.Iter, error) {
	return fc.fetchCollection(fc.outboxIRI)
}

func (fc *fedbaseclient) OutboxIRI() *url.URL {
	return fc.outboxIRI
}

func (fc *fedbaseclient) Liked() (fetch.Iter, error) {
	return fc.fetchCollection(fc.likedIRI)
}

func (fc *fedbaseclient) LikedIRI() *url.URL {
	return fc.inboxIRI
}

func (fc *fedbaseclient) Followers() (fetch.Iter, error) {
	return fc.fetchCollection(fc.followersIRI)
}

func (fc *fedbaseclient) FollowersIRI() *url.URL {
	return fc.followersIRI
}

func (fc *fedbaseclient) Create(event vocab.Type) error {
	return fc.create(event)
}

func (fc *fedbaseclient) Like(iri *url.URL) error {
	return fc.like(iri)
}

func (fc *fedbaseclient) fetchCollection(target *url.URL) (fetch.Iter, error) {
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
