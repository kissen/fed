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
