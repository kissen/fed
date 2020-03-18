package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"net/url"
)

func Submit(to *url.URL, object vocab.Type) error {
	bs, err := VocabToBytes(object)
	if err != nil {
		return errors.Wrap(err, "serialization failed before submitting")
	}

	return Post(bs, to)
}
