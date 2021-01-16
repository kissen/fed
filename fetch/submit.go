package fetch

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/fed/errors"
	"github.com/kissen/fed/marshal"
	"net/url"
)

// Submit activity pub payload object to iri. This should involve
// an HTTP POST request to iri.
func Submit(object vocab.Type, iri *url.URL) error {
	bs, err := marshal.VocabToBytes(object)
	if err != nil {
		return errors.Wrap(err, "serialization failed before submitting")
	}

	return Post(bs, iri)
}
