package fedcontext

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/fed/fetch"
	"github.com/kissen/fed/prop"
	"github.com/kissen/fed/util"
	"net/url"
	"time"
)

func NewRemoteClient(actoraddr, token string) (FedClient, error) {
	bc := &fedbaseclient{}

	if err := bc.fill(actoraddr); err != nil {
		return nil, err
	}

	bc.create = func(event vocab.Type) error {
		if create, err := createCreate(bc, event); err != nil {
			return err
		} else {
			target := bc.OutboxIRI()
			return submitWithToken(create, target, token)
		}

	}

	bc.like = func(iri *url.URL) error {
		like := createLike(bc, iri)
		target := bc.OutboxIRI()
		return submitWithToken(like, target, token)
	}

	return bc, nil
}

// Submit obj to target with OAuth token authorization.
func submitWithToken(obj vocab.Type, target *url.URL, token string) error {
	target = util.WithParam(target, "token", token)
	return fetch.Submit(obj, target)
}

// Wrap event into a Create activity.
func createCreate(fc FedClient, event vocab.Type) (vocab.ActivityStreamsCreate, error) {
	// check whether event is a supported type
	note, ok := event.(vocab.ActivityStreamsNote)
	if !ok {
		return nil, fmt.Errorf("event of type=%T cannot be wrapped in Create", event)
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
	audience.AppendIRI(fc.FollowersIRI())
	create.SetActivityStreamsAudience(audience)

	return create, nil
}

// Create a Like activity for the object at iri.
func createLike(fc FedClient, iri *url.URL) vocab.ActivityStreamsLike {
	like := streams.NewActivityStreamsLike()
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(iri)
	like.SetActivityStreamsObject(object)
	return like
}
