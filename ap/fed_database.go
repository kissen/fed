package ap

import (
	"context"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
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
	log.Printf("InboxContains(inbox=%v id=%v)", inbox, id)

	inboxIri := IRI{c, inbox}

	if user, err := inboxIri.RetrieveOwner(); err != nil {
		return false, err
	} else {
		return urlIn(id, user.Inbox), nil
	}
}

// GetInbox returns the first ordered collection page of the outbox at
// the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetInbox(c context.Context, inboxIRI *url.URL) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	log.Printf("GetInbox(%v)", inboxIRI)

	iri := IRI{c, inboxIRI}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, err
	} else if page, err := collectPage(c, user.Inbox); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		fedutil.SetIdOn(page, iri.URL())
		return page, nil
	}
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetInbox()")

	id := fedutil.Id(inbox)
	iri := IRI{c, id}

	if user, err := iri.RetrieveOwner(); err != nil {
		return err
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
	log.Printf("Owns(%v)", id)

	// first, check if it is an retrievable object

	if _, err := FromContext(c).Storage.RetrieveObject(id); err == nil {
		return true, nil
	}

	// it isn't; check if it's a users collection

	iri := IRI{c, id}

	if user, err := iri.RetrieveOwner(); err == nil {
		if urlInAny(id, user.Collections()) {
			return true, nil
		}
	}

	// not owned by us

	return false, nil
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Printf("ActorForOutbox(%v)", outboxIRI)

	iri := IRI{c, outboxIRI}

	if username, err := iri.OutboxOwner(); err != nil {
		return nil, err
	} else {
		return ActorIRI(c, username).URL(), nil
	}
}

// ActorForInbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Printf("ActorForInbox(%v)", inboxIRI)

	iri := IRI{c, inboxIRI}

	if username, err := iri.InboxOwner(); err != nil {
		return nil, err
	} else {
		return ActorIRI(c, username).URL(), nil
	}
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	log.Printf("OutboxForInbox(%v)", outboxIRI)

	iri := IRI{c, inboxIRI}

	if username, err := iri.InboxOwner(); err != nil {
		return nil, err
	} else {
		return InboxIRI(c, username).URL(), nil
	}
}

// Exists returns true if the database has an entry for the specified
// id. It may not be owned by this application instance.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	log.Printf("Exists(%v)", id)

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
func (f *FedDatabase) Get(c context.Context, addr *url.URL) (value vocab.Type, err error) {
	log.Printf("Get(%v)", addr)

	iri := IRI{c, addr}

	// try out collections

	if _, err := iri.InboxOwner(); err == nil {
		return f.GetInbox(c, iri.URL())
	}

	if _, err := iri.OutboxOwner(); err == nil {
		return f.GetOutbox(c, iri.URL())
	}

	// try out actors

	if _, err := iri.Actor(); err == nil {
		return f.getActor(c, iri.URL())
	}

	// try serving plain documents

	if obj, err := FromContext(c).Storage.RetrieveObject(iri.URL()); err != nil {
		return nil, errors.Wrapf(err, "cannot retrieve addr=%v", iri)
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

	target := fedutil.Id(asType)
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

	target := fedutil.Id(asType)
	return FromContext(c).Storage.StoreObject(target, asType)
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Delete(c context.Context, id *url.URL) error {
	log.Printf("Delete(%v)", id)

	return FromContext(c).Storage.DeleteObject(id)
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetOutbox(c context.Context, outboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("GetOutbox(%v)", outboxIRI)

	iri := IRI{c, outboxIRI}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, errors.Wrap(err, "no such outbox")
	} else if page, err := collectPage(c, user.Outbox); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		fedutil.SetIdOn(page, iri.URL())
		return page, nil
	}
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) SetOutbox(c context.Context, outbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("SetOutbox()")

	id := fedutil.Id(outbox)
	iri := IRI{c, id}

	if user, err := iri.RetrieveOwner(); err != nil {
		return err
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

	return RollObjectIRI(c).URL(), nil
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Followers(%v)", actorIRI)

	iri := IRI{c, actorIRI}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Followers); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		fedutil.SetIdOn(set, iri.URL())
		return set, nil
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

	iri := IRI{c, actorIRI}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Following); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		fedutil.SetIdOn(set, iri.URL())
		return set, nil
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

	iri := IRI{c, actorIRI}

	if user, err := iri.RetrieveOwner(); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Liked); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		fedutil.SetIdOn(set, iri.URL())
		return set, nil
	}
}

func (f *FedDatabase) getActor(c context.Context, actorIRI *url.URL) (actor vocab.ActivityStreamsPerson, err error) {
	// look up the user

	var user *db.FedUser

	iri := IRI{c, actorIRI}

	if user, err = iri.RetrieveOwner(); err != nil {
		return nil, errors.Wrap(err, "not an actor")
	}

	// build up the actor object

	actor = streams.NewActivityStreamsPerson()
	fedutil.SetIdOn(actor, iri.URL())

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(user.Name)
	actor.SetActivityStreamsName(name)

	inbox := streams.NewActivityStreamsInboxProperty()
	inbox.SetIRI(InboxIRI(c, user.Name).URL())
	actor.SetActivityStreamsInbox(inbox)

	outbox := streams.NewActivityStreamsOutboxProperty()
	outbox.SetIRI(OutboxIRI(c, user.Name).URL())
	actor.SetActivityStreamsOutbox(outbox)

	followers := streams.NewActivityStreamsFollowersProperty()
	followers.SetIRI(FollowersIRI(c, user.Name).URL())
	actor.SetActivityStreamsFollowers(followers)

	following := streams.NewActivityStreamsFollowingProperty()
	following.SetIRI(FollowingIRI(c, user.Name).URL())
	actor.SetActivityStreamsFollowing(following)

	liked := streams.NewActivityStreamsLikedProperty()
	liked.SetIRI(LikedIRI(c, user.Name).URL())
	actor.SetActivityStreamsLiked(liked)

	return actor, nil
}

// Ensure that all objects in collection are part of our storage. Returns a
// list of all IRIs of all the objects in collection.
func (f *FedDatabase) addToStorage(c context.Context, collection vocab.ActivityStreamsOrderedCollectionPage) (colIRIs []*url.URL, err error) {
	items := collection.GetActivityStreamsOrderedItems()

	for it := items.Begin(); it != items.End(); it = it.Next() {
		if it.IsIRI() {
			// IRIs are only links and do not need to be stored in our database;
			// we only have to add it to the results slice s.t. it is recorded
			// e.g. in a users outbox
			colIRIs = append(colIRIs, it.GetIRI())
		} else if obj := it.GetType(); obj == nil {
			// items that are not IRIs really should be full objects; if they are
			// not something is probably wrong
			panic("obj is nil")
		} else if err := FromContext(c).Storage.StoreObject(fedutil.Id(obj), obj); err != nil {
			// XXX: here we quit while having modified the database; maybe need
			// to think about transaction for the FedStorage interface to easily
			// roll back such changes
			return nil, errors.Wrapf(err, "cannot store iri=%v", fedutil.Id(obj))
		} else {
			// object was successfully added to database
			colIRIs = append(colIRIs, fedutil.Id(obj))
		}
	}

	return
}
