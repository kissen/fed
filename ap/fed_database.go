package ap

import (
	"context"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
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
	f.lock.Lock()
	return nil
}

// Unlock makes the lock for the object at the specified id available.
// If an error is returned, the lock must have still been freed.
//
// Used to ensure race conditions in multiple requests do not occur.
func (f *FedDatabase) Unlock(c context.Context, id *url.URL) error {
	f.lock.Unlock()
	return nil
}

// InboxContains returns true if the OrderedCollection at 'inbox'
// contains the specified 'id'.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
	// get id of item to search

	needle, err := parseActivityIdFromIri(c, id)
	if err != nil {
		return false, errors.Wrapf(err, "could not determine activity id from id=%v", id)
	}

	// look up owner of inbox

	username, err := parseInboxOwnerFromIri(c, inbox)
	if err != nil {
		return false, errors.Wrapf(err, "could not determine owner of inbox=%v", inbox)
	}

	user, err := FromContext(c).Storage.FindUser(username)
	if err != nil {
		return false, errors.Wrapf(err, "no user found for username=%v", username)
	}

	// get posts

	posts, err := FromContext(c).Storage.GetPostsFrom(user.Id)
	if err != nil {
		return false, errors.Wrapf(err, "no posts found for username=%v", username)
	}

	// search

	for _, post := range posts {
		if post.Id == needle {
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
	// look up owner of inbox
	username, err := parseInboxOwnerFromIri(c, inboxIRI)
	if err != nil {
		return nil, errors.Wrapf(err, "could not determine owner of inbox=%v", inbox)
	}

	user, err := FromContext(c).Storage.FindUser(username)
	if err != nil {
		return nil, errors.Wrapf(err, "no user found for username=%v", username)
	}

	// get posts

	posts, err := FromContext(c).Storage.GetPostsFrom(user.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "no posts found for username=%v", username)
	}

	// build up go-fed data type

	inbox = streams.NewActivityStreamsOrderedCollectionPage()
	notes := convertPostsToNotes(posts)
	inbox.SetActivityStreamsOrderedItems(notes)

	return inbox, nil
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	return errors.New("not implemented")
}

// Owns returns true if the database has an entry for the IRI and it
// exists in the database.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	return false, errors.New("not implemented")
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	return nil, errors.New("not implemented")
}

// ActorForInbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	return nil, errors.New("not implemented")
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	return nil, errors.New("not implemented")
}

// Exists returns true if the database has an entry for the specified
// id. It may not be owned by this application instance.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	return false, errors.New("not implemented")
}

// Get returns the database entry for the specified id.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
	return nil, errors.New("not implemented")
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
	return errors.New("not implemented")
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Delete(c context.Context, id *url.URL) error {
	return errors.New("not implemented")
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetOutbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	return nil, errors.New("not implemented")
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetOutbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	return errors.New("not implemented")
}

// NewId creates a new IRI id for the provided activity or object. The
// implementation does not need to set the 'id' property and simply
// needs to determine the value.
//
// The go-fed library will handle setting the 'id' property on the
// activity or object provided with the value returned.
func (f *FedDatabase) NewId(c context.Context, t vocab.Type) (id *url.URL, err error) {
	return nil, errors.New("not implemented")
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	return nil, errors.New("not implemented")
}

// Following obtains the Following Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Following(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	return nil, errors.New("not implemented")
}

// Liked obtains the Liked Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Liked(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	return nil, errors.New("not implemented")
}
