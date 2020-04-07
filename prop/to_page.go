package prop

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Return a new ordered collection page that contains all
// elements in iris.
func ToPage(iris []*url.URL) vocab.ActivityStreamsOrderedCollectionPage {
	items := streams.NewActivityStreamsOrderedItemsProperty()
	AppendIRIs(items, iris)

	page := streams.NewActivityStreamsOrderedCollectionPage()
	page.SetActivityStreamsOrderedItems(items)

	return page
}

// Return a new collection that contains all elements in iris.
func ToCollection(iris []*url.URL) vocab.ActivityStreamsCollection {
	items := streams.NewActivityStreamsItemsProperty()
	AppendIRIs(items, iris)

	collection := streams.NewActivityStreamsCollection()
	collection.SetActivityStreamsItems(items)

	return collection
}
