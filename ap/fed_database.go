package ap

import (
	"context"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/fediri"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"log"
	"net/http"
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

	inboxIri := fediri.IRI{inbox}

	if user, err := retrieveOwner(&inboxIri, fedcontext.From(c).Storage); err != nil {
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

	iri := fediri.IRI{inboxIRI}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, err
	} else if page, err := collectPage(c, user.Inbox); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(page, iri.URL())
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

	id := prop.Id(inbox)
	iri := fediri.IRI{id}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return err
	} else if slice, err := f.addToStorage(c, inbox); err != nil {
		return err
	} else {
		user.Inbox = slice
		return fedcontext.From(c).Storage.StoreUser(user)
	}
}

// Owns returns true if the database has an entry for the IRI and it
// exists in the database.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	log.Printf("Owns(%v)", id)

	// first, check if it is an retrievable object

	if _, err := fedcontext.From(c).Storage.RetrieveObject(id); err == nil {
		return true, nil
	}

	// it isn't; check if it's a users collection

	iri := fediri.IRI{id}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err == nil {
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

	iri := fediri.IRI{outboxIRI}

	if username, err := iri.OutboxOwner(); err != nil {
		return nil, err
	} else {
		return fediri.ActorIRI(username).URL(), nil
	}
}

// ActorForInbox fetches the actor's IRI for the given inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Printf("ActorForInbox(%v)", inboxIRI)

	iri := fediri.IRI{inboxIRI}

	if username, err := iri.InboxOwner(); err != nil {
		return nil, err
	} else {
		return fediri.ActorIRI(username).URL(), nil
	}
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	log.Printf("OutboxForInbox(%v)", outboxIRI)

	iri := fediri.IRI{inboxIRI}

	if username, err := iri.InboxOwner(); err != nil {
		return nil, err
	} else {
		return fediri.InboxIRI(username).URL(), nil
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

	iri := fediri.IRI{addr}

	// try out collections

	if _, err := iri.InboxOwner(); err == nil {
		return f.GetInbox(c, iri.URL())
	}

	if _, err := iri.OutboxOwner(); err == nil {
		return f.GetOutbox(c, iri.URL())
	}

	if _, err := iri.FollowingOwner(); err == nil {
		return f.Following(c, iri.URL())
	}

	if _, err := iri.FollowersOwner(); err == nil {
		return f.Following(c, iri.URL())
	}

	if _, err := iri.LikedOwner(); err == nil {
		return f.Liked(c, iri.URL())
	}

	// try out actors

	if _, err := iri.Actor(); err == nil {
		return f.getActor(c, iri.URL())
	}

	// try serving plain documents

	if obj, err := fedcontext.From(c).Storage.RetrieveObject(iri.URL()); err != nil {
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

	target := prop.Id(asType)
	return fedcontext.From(c).Storage.StoreObject(target, asType)
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

	id := prop.Id(asType)
	iri := fediri.IRI{id}

	log.Printf("\n\nid=%v\n\n", id)

	// try out collections

	if _, err := iri.InboxOwner(); err == nil {
		return errors.NewWith(http.StatusNotImplemented, "update of inbox not supported")
	}

	if _, err := iri.OutboxOwner(); err == nil {
		return errors.NewWith(http.StatusNotImplemented, "update of owner not supported")
	}

	if _, err := iri.FollowingOwner(); err == nil {
		return errors.NewWith(http.StatusNotImplemented, "update of following not supported")
	}

	if _, err := iri.FollowersOwner(); err == nil {
		return errors.NewWith(http.StatusNotImplemented, "update of followers not supported")
	}

	if _, err := iri.LikedOwner(); err == nil {
		if liked, ok := asType.(vocab.ActivityStreamsCollection); !ok {
			return errors.NewfWith(http.StatusInternalServerError, "bad runtime type %T for liked collection", asType)
		} else {
			return f.updateLiked(c, iri, liked)
		}
	}

	// try out actors

	if _, err := iri.Actor(); err == nil {
		if person, ok := asType.(vocab.ActivityStreamsPerson); !ok {
			return errors.NewfWith(http.StatusInternalServerError, "bad runtime type %T for actor", asType)
		} else {
			return f.updatePerson(c, iri, person)
		}
	}

	// try storage as a last resort

	return fedcontext.From(c).Storage.StoreObject(id, asType)
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Delete(c context.Context, id *url.URL) error {
	log.Printf("Delete(%v)", id)

	return fedcontext.From(c).Storage.DeleteObject(id)
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) GetOutbox(c context.Context, outboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("GetOutbox(%v)", outboxIRI)

	iri := fediri.IRI{outboxIRI}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, errors.Wrap(err, "no such outbox")
	} else if page, err := collectPage(c, user.Outbox); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(page, iri.URL())
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

	id := prop.Id(outbox)
	iri := fediri.IRI{id}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return err
	} else if slice, err := f.addToStorage(c, outbox); err != nil {
		return err
	} else {
		user.Outbox = slice
		return fedcontext.From(c).Storage.StoreUser(user)
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

	return fediri.RollObjectIRI().URL(), nil
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (f *FedDatabase) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	log.Printf("Followers(%v)", actorIRI)

	iri := fediri.IRI{actorIRI}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Followers); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(set, fediri.FollowersIRI(user.Name).URL())
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

	iri := fediri.IRI{actorIRI}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Following); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(set, fediri.FollowingIRI(user.Name).URL())
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

	iri := fediri.IRI{actorIRI}

	if user, err := retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, err
	} else if set, err := collectSet(c, user.Liked); err != nil {
		return nil, errors.Wrap(err, "collect failed")
	} else {
		prop.SetIdOn(set, fediri.LikedIRI(user.Name).URL())
		return set, nil
	}
}

// Return the ActivityStreams representation of the actor at actorIRI.
func (f *FedDatabase) getActor(c context.Context, actorIRI *url.URL) (actor vocab.ActivityStreamsPerson, err error) {
	// look up the user

	var user *db.FedUser

	iri := fediri.IRI{actorIRI}

	if user, err = retrieveOwner(&iri, fedcontext.From(c).Storage); err != nil {
		return nil, errors.Wrap(err, "not an actor")
	}

	// build up the actor object

	actor = streams.NewActivityStreamsPerson()
	prop.SetIdOn(actor, iri.URL())

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(user.Name)
	actor.SetActivityStreamsName(name)

	inbox := streams.NewActivityStreamsInboxProperty()
	inbox.SetIRI(fediri.InboxIRI(user.Name).URL())
	actor.SetActivityStreamsInbox(inbox)

	outbox := streams.NewActivityStreamsOutboxProperty()
	outbox.SetIRI(fediri.OutboxIRI(user.Name).URL())
	actor.SetActivityStreamsOutbox(outbox)

	followers := streams.NewActivityStreamsFollowersProperty()
	followers.SetIRI(fediri.FollowersIRI(user.Name).URL())
	actor.SetActivityStreamsFollowers(followers)

	following := streams.NewActivityStreamsFollowingProperty()
	following.SetIRI(fediri.FollowingIRI(user.Name).URL())
	actor.SetActivityStreamsFollowing(following)

	liked := streams.NewActivityStreamsLikedProperty()
	liked.SetIRI(fediri.LikedIRI(user.Name).URL())
	actor.SetActivityStreamsLiked(liked)

	return actor, nil
}

func (f *FedDatabase) updateLiked(c context.Context, actoriri fediri.IRI, liked vocab.ActivityStreamsCollection) error {
	// XXX: racy: need transactions

	storage := fedcontext.From(c).Storage

	username, err := actoriri.LikedOwner()
	if err != nil {
		return err
	}

	user, err := storage.RetrieveUser(username)
	if err != nil {
		return err
	}

	user.Liked, err = f.iris(liked)
	if err != nil {
		return errors.Wrap(err, "bad liked collection")
	}

	return storage.StoreUser(user)
}

// Update actor which should represent a user on our instance. In particular, update
// the liked, follows (and so on) collections.
func (f *FedDatabase) updatePerson(c context.Context, actoriri fediri.IRI, actor vocab.ActivityStreamsPerson) error {
	storage := fedcontext.From(c).Storage

	// fetch metadata from the database

	username, err := actoriri.Actor()
	if err != nil {
		return err
	}

	user, err := storage.RetrieveUser(username)
	if err != nil {
		return err
	}

	// for each supported collection, fetch the iris and update
	// the user meta data accordingly

	user.Followers, err = f.iris(actor.GetActivityStreamsFollowers)
	if err != nil {
		return errors.Wrap(err, "bad followers collection")
	}

	user.Following, err = f.iris(actor.GetActivityStreamsFollowing)
	if err != nil {
		return errors.Wrap(err, "bad following collection")
	}

	user.Liked, err = f.iris(actor.GetActivityStreamsLiked)
	if err != nil {
		return errors.Wrap(err, "bad liked collection")
	}

	// update
	// XXX: racy; use transactions here (see also /NEXTUP.md)

	if err := storage.StoreUser(user); err != nil {
		return errors.Wrap(err, "overwriting user failed")
	}

	return nil
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
		} else if err := fedcontext.From(c).Storage.StoreObject(prop.Id(obj), obj); err != nil {
			// XXX: here we quit while having modified the database; maybe need
			// to think about transaction for the FedStorage interface to easily
			// roll back such changes
			return nil, errors.Wrapf(err, "cannot store iri=%v", prop.Id(obj))
		} else {
			// object was successfully added to database
			colIRIs = append(colIRIs, prop.Id(obj))
		}
	}

	return
}

func (f *FedDatabase) iris(from interface{}) ([]*url.URL, error) {
	it, err := fetch.Begin(from)
	if err != nil {
		return nil, err
	}

	return fetch.IRIs(it)
}
