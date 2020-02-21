package ap

import (
	"context"
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/db"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"sync"
)

// Implements the go-fed/activity/pub/Datbase interface (version 1.0)
type FedDatabase struct {
	lock sync.Mutex
}

// Lock takes a lock for the object at the specified id. If an error
// is returned, the lock must not have been taken.
//
// The lock must be able to succeed for an id that does not exist in
// the database. This means acquiring the lock does not guarantee the
// entry exists in the database.
//
// Locks are encouraged to be lightweight and in the Go layer, as some
// processes require tight loops acquiring and releasing locks.
//
// Used to ensure race conditions in multiple requests do not occur.
func (f *FedDatabase) Lock(c context.Context, id *url.URL) error {
	log.Printf("Lock(%v)", id)

	f.lock.Lock()
	return nil
}

// Unlock makes the lock for the object at the specified id available.
// If an error is returned, the lock must have still been freed.
//
// Used to ensure race conditions in multiple requests do not occur.
func (f *FedDatabase) Unlock(c context.Context, id *url.URL) error {
	log.Printf("Unlock(%v)", id)

	f.lock.Unlock()
	return nil
}

// InboxContains returns true if the OrderedCollection at 'inbox'
// contains the specified 'id'.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
	log.Printf("InboxContains(inbox=%v id=%v)\n", inbox, id)

	// get user

	var user *db.FedUser

	if username, err := parseInboxOwnerFromIri(c, inbox); err != nil {
		return false, errors.Wrapf(err, "could not determine owner of inbox=%v", inbox)
	} else if user, err = FromContext(c).Storage.RetrieveUser(username); err != nil {
		return false, errors.Wrapf(err, "no user found for username=%v", username)
	}

	// look for post in inbox

	for _, member := range user.Inbox {
		if urlEq(id, member) {
			return true, nil
		}
	}

	return false, nil
}

// GetInbox returns the first ordered collection page of the outbox at
// the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("GetInbox(%v)\n", inboxIRI)

	inbox = streams.NewActivityStreamsOrderedCollectionPage()
	return inbox, errors.New("not implemented")
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetInbox()")

	return errors.New("not implemented")
}

func (f *FedDatabase) ownsInbox(c context.Context, iri *url.URL) bool {
	if username, err := parseInboxOwnerFromIri(c, iri); err != nil {
		return false
	} else {
		_, err := FromContext(c).Storage.RetrieveUser(username)
		return err != nil
	}
}

func (f *FedDatabase) ownsOutbox(c context.Context, iri *url.URL) bool {
	if username, err := parseOutboxOwnerFromIri(c, iri); err != nil {
		return false
	} else {
		_, err := FromContext(c).Storage.RetrieveUser(username)
		return err != nil
	}
}

func (f *FedDatabase) ownsActivity(c context.Context, iri *url.URL) bool {
	_, err := FromContext(c).Storage.RetrieveObject(iri)
	return err != nil
}

// Owns returns true if the database has an entry for the IRI and it
// exists in the database.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	log.Println("Owns()")

	if f.ownsInbox(c, id) {
		return true, nil
	}

	if f.ownsOutbox(c, id) {
		return true, nil
	}

	if f.ownsActivity(c, id) {
		return true, nil
	}

	return false, nil
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Printf("ActorForOutbox(%v)\n", outboxIRI)

	if username, err := parseOutboxOwnerFromIri(c, outboxIRI); err != nil {
		return nil, err
	} else {
		return constructActorIri(c, username), nil
	}
}

// ActorForInbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Printf("ActorForInbox(%v)\n", inboxIRI)

	if username, err := parseInboxOwnerFromIri(c, inboxIRI); err != nil {
		return nil, err
	} else {
		return constructActorIri(c, username), nil
	}
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	log.Printf("OutboxForInbox(%v)\n", outboxIRI)

	if username, err := parseInboxOwnerFromIri(c, inboxIRI); err != nil {
		return nil, err
	} else {
		return constructOutboxIri(c, username), nil
	}
}

// Exists returns true if the database has an entry for the specified
// id. It may not be owned by this application instance.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	log.Printf("Exists(%v)\n", id)

	if exists, err = f.Owns(c, id); err != nil {
		return exists, errors.Wrap(err, "using Owns() to implement Exists() failed")
	} else {
		return exists, nil
	}
}

// Get returns the database entry for the specified id.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Get(c context.Context, iri *url.URL) (value vocab.Type, err error) {
	log.Printf("Get(%v)\n", iri)

	// try out collections

	if _, err := parseInboxOwnerFromIri(c, iri); err == nil {
		return f.GetInbox(c, iri)
	}

	if _, err := parseOutboxOwnerFromIri(c, iri); err == nil {
		return f.GetOutbox(c, iri)
	}

	// try out actors

	if actor, err := parseActorFromIri(c, iri); err == nil {
		// TODO
		return nil, fmt.Errorf("getting actor=%v not yet implemented", actor)
	}

	// try serving plain documents

	if obj, err := FromContext(c).Storage.RetrieveObject(iri); err != nil {
		return nil, errors.Wrapf(err, "could not fetch iri=%v from storage", iri)
	} else {
		return obj, nil
	}
}

// Create adds a new entry to the database which must be able to be
// keyed by its id.
//
// Note that Activity values received from federated peers may also be
// created in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
//
// Under certain conditions and network activities, Create may be called
// multiple times for the same ActivityStreams object.
func (f *FedDatabase) Create(c context.Context, asType vocab.Type) error {
	log.Println("Create()")

	return errors.New("not implemented")
}

// Update sets an existing entry to the database based on the value's
// id.
//
// Note that Activity values received from federated peers may also be
// updated in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Update(c context.Context, asType vocab.Type) error {
	log.Println("Update()")

	return errors.New("not implemented")
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Delete(c context.Context, id *url.URL) error {
	log.Printf("Delete(%v)\n", id)

	return errors.New("not implemented")
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetOutbox(c context.Context, outboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("GetOutbox(%v)\n", outboxIRI)

	return nil, errors.New("not implemented")
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetOutbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetOutbox()")

	return errors.New("not implemented")
}

// NewId creates a new IRI id for the provided activity or object. The
// implementation does not need to set the 'id' property and simply
// needs to determine the value.
//
// The go-fed library will handle setting the 'id' property on the
// activity or object provided with the value returned.
func (f *FedDatabase) NewId(c context.Context, t vocab.Type) (id *url.URL, err error) {
	log.Println("NewId()")

	return nil, errors.New("not implemented")
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Followers(%v)", actorIRI)

	return nil, errors.New("not implemented")
}

// Following obtains the Following Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Following(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Following(%v)", actorIRI)

	return nil, errors.New("not implemented")
}

// Liked obtains the Liked Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Liked(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Liked(%v)", actorIRI)

	return nil, errors.New("not implemented")
}
