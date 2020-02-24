package ap

import (
	"context"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/help"
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

	var user *db.FedUser

	if user, err = parseUserFrom(c, parseInboxOwnerFromIri, inbox); err != nil {
		return false, errors.Wrap(err, "no such inbox")
	}

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
func (f *FedDatabase) GetInbox(c context.Context, inboxIRI *url.URL) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	log.Printf("GetInbox(%v)\n", inboxIRI)

	if user, err := parseUserFrom(c, parseInboxOwnerFromIri, inboxIRI); err != nil {
		return nil, errors.Wrap(err, "no such inbox")
	} else {
		return collectPage(c, user.Inbox)
	}
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetInbox()")

	if user, err := parseUserFrom(c, parseInboxOwnerFromIri, help.Id(inbox)); err != nil {
		return errors.Wrap(err, "unknown user")
	} else if slice, err := f.addToStorage(c, inbox); err != nil {
		return err
	} else {
		user.Inbox = slice
		return FromContext(c).Storage.StoreUser(user)
	}
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

	// TODO: custom impl

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

	target := help.Id(asType)
	return FromContext(c).Storage.StoreObject(target, asType)
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

	target := help.Id(asType)
	return FromContext(c).Storage.StoreObject(target, asType)
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Delete(c context.Context, id *url.URL) error {
	log.Printf("Delete(%v)\n", id)

	return FromContext(c).Storage.DeleteObject(id)
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetOutbox(c context.Context, outboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("GetOutbox(%v)\n", outboxIRI)

	if user, err := parseUserFrom(c, parseOutboxOwnerFromIri, outboxIRI); err != nil {
		return nil, errors.Wrap(err, "no such outbox")
	} else {
		return collectPage(c, user.Outbox)
	}
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetOutbox(c context.Context, outbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetOutbox()")

	if user, err := parseUserFrom(c, parseInboxOwnerFromIri, help.Id(outbox)); err != nil {
		return errors.Wrap(err, "unknown user")
	} else if slice, err := f.addToStorage(c, outbox); err != nil {
		return err
	} else {
		user.Outbox = slice
		return FromContext(c).Storage.StoreUser(user)
	}
}

// NewId creates a new IRI id for the provided activity or object. The
// implementation does not need to set the 'id' property and simply
// needs to determine the value.
//
// The go-fed library will handle setting the 'id' property on the
// activity or object provided with the value returned.
func (f *FedDatabase) NewId(c context.Context, t vocab.Type) (id *url.URL, err error) {
	log.Println("NewId()")

	return rollObjectIri(c), nil
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Followers(%v)", actorIRI)

	if user, err := parseUserFrom(c, parseActorFromIri, actorIRI); err != nil {
		return nil, err
	} else {
		return collectSet(c, user.Followers)
	}
}

// Following obtains the Following Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Following(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Following(%v)", actorIRI)

	if user, err := parseUserFrom(c, parseActorFromIri, actorIRI); err != nil {
		return nil, err
	} else {
		return collectSet(c, user.Following)
	}
}

// Liked obtains the Liked Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Liked(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Liked(%v)", actorIRI)

	if user, err := parseUserFrom(c, parseActorFromIri, actorIRI); err != nil {
		return nil, err
	} else {
		return collectSet(c, user.Liked)
	}
}

func (f *FedDatabase) ownsInbox(c context.Context, iri *url.URL) bool {
	_, err := parseUserFrom(c, parseInboxOwnerFromIri, iri)
	return err != nil
}

func (f *FedDatabase) ownsOutbox(c context.Context, iri *url.URL) bool {
	_, err := parseUserFrom(c, parseOutboxOwnerFromIri, iri)
	return err != nil
}

func (f *FedDatabase) ownsActivity(c context.Context, iri *url.URL) bool {
	_, err := FromContext(c).Storage.RetrieveObject(iri)
	return err != nil
}

// Ensure that all objects in collection are part of our storage. Returns a
// list of all IRIs of all the objects in collection.
func (f *FedDatabase) addToStorage(c context.Context, collection vocab.ActivityStreamsOrderedCollectionPage) ([]*url.URL, error) {
	items := collection.GetActivityStreamsOrderedItems()
	var result []*url.URL

	for it := items.Begin(); it != nil; it = it.Next() {
		obj := it.GetActivityStreamsObject()
		id := help.Id(obj)

		if err := FromContext(c).Storage.StoreObject(id, obj); err != nil {
			return nil, errors.Wrap(err, "at least one entry was not understood")
		}

		result = append(result, id)
	}

	return result, nil
}
