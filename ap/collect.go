package ap

import (
	"context"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"net/url"
)

func collectPage(c context.Context, iris []*url.URL) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	items := streams.NewActivityStreamsOrderedItemsProperty()

	for _, iri := range iris {
		if obj, err := fedcontext.From(c).Storage.RetrieveObject(iri); err != nil {
			return nil, errors.Wrapf(err, "missing iri=%v in database", iri)
		} else if err := items.AppendType(obj); err != nil {
			return nil, errors.Wrapf(err, "cannot add iri=%v to items", iri)
		}
	}

	page := streams.NewActivityStreamsOrderedCollectionPage()
	page.SetActivityStreamsOrderedItems(items)

	return page, nil
}

func collectSet(c context.Context, iris []*url.URL) (vocab.ActivityStreamsCollection, error) {
	items := streams.NewActivityStreamsItemsProperty()

	for _, iri := range iris {
		if obj, err := fedcontext.From(c).Storage.RetrieveObject(iri); err != nil {
			return nil, errors.Wrapf(err, "missing iri=%v in database", iri)
		} else if err := items.AppendType(obj); err != nil {
			return nil, errors.Wrapf(err, "cannot add iri=%v to items", iri)
		}
	}

	collection := streams.NewActivityStreamsCollection()
	collection.SetActivityStreamsItems(items)

	return collection, nil
}
