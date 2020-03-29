package fedcontext

import (
	"fmt"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"gitlab.cs.fau.de/kissen/fed/prop"
	"gitlab.cs.fau.de/kissen/fed/util"
	"time"
)

func NewRemoteClient(actoraddr, token string) (FedClient, error) {
	bc := &fedbaseclient{}

	if err := bc.fill(actoraddr); err != nil {
		return nil, err
	}

	bc.create = func(event vocab.Type) error {
		create, err := createCreate(bc, event)
		if err != nil {
			return err
		}

		// append token to target so the request is authenticated
		target := bc.outboxIRI
		target = util.WithParam(target, "token", token)
		return fetch.Submit(create, target)
	}

	return bc, nil
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
