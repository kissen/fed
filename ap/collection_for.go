package ap

import (
	"context"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/db"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/fediri"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"log"
	"net/url"
)

// For whoever owns addr, return the contents of field (e.g. "Following")
// as an Activity Streams collection.
//
// The returned collection contains only IRIs, it is up to an Activity
// Pub client to dereference these IRIs.
func collectionFor(c context.Context, addr *url.URL, field string) (vocab.ActivityStreamsCollection, error) {
	storage := fedcontext.From(c).Storage
	reqIRI := fediri.IRI{addr}

	user, err := retrieveOwner(&reqIRI, storage)
	if err != nil {
		return nil, err
	}

	var collection vocab.ActivityStreamsCollection
	var id fediri.IRI

	switch field {
	case "Following":
		collection = prop.ToCollection(user.Following)
		id = fediri.FollowingIRI(user.Name)
	case "Followers":
		collection = prop.ToCollection(user.Followers)
		id = fediri.FollowersIRI(user.Name)
	case "Liked":
		collection = prop.ToCollection(user.Liked)
		id = fediri.LikedIRI(user.Name)
	default:
		log.Fatalf("bad field=%v", field)
	}

	prop.SetIdOn(collection, id.URL())
	return collection, nil
}

// For whoever owns addr, return the contents of field (e.g. "Following")
// as an Activity Streams ordered collection page.
//
// The returned collection contains only IRIs, it is up to an Activity
// Pub client to dereference these IRIs.
func pageFor(c context.Context, addr *url.URL, field string) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	storage := fedcontext.From(c).Storage
	reqIRI := fediri.IRI{addr}

	user, err := retrieveOwner(&reqIRI, storage)
	if err != nil {
		return nil, err
	}

	var page vocab.ActivityStreamsOrderedCollectionPage
	var id fediri.IRI

	switch field {
	case "Inbox":
		page = prop.ToPage(user.Inbox)
		id = fediri.InboxIRI(user.Name)
	case "Outbox":
		page = prop.ToPage(user.Outbox)
		id = fediri.OutboxIRI(user.Name)
	default:
		log.Fatalf("bad field=%v", field)
	}

	prop.SetIdOn(page, id.URL())
	return page, nil
}

// Return the owner of this IRI.
func retrieveOwner(iri *fediri.IRI, from db.FedStorage) (*db.FedUser, error) {
	// the owner of an IRI, in the easy case, is the first
	// path component; we do not support getting the owner
	// of object IRIs yet

	if username, err := iri.Owner(); err != nil {
		return nil, err
	} else {
		return from.RetrieveUser(username)
	}
}
