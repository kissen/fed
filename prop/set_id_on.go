package prop

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// Set the JSONLDId field on target to iri.
func SetIdOn(target vocab.Type, iri *url.URL) {
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(iri)
	target.SetJSONLDId(id)
}
