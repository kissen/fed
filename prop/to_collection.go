package prop

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Return a new collection that contains all elements in iris.
func ToCollection(iris []*url.URL) vocab.ActivityStreamsCollection {
	items := streams.NewActivityStreamsItemsProperty()
	AppendIRIs(items, iris)

	collection := streams.NewActivityStreamsCollection()
	collection.SetActivityStreamsItems(items)

	return collection
}
