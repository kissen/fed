package ap

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/fetch"
	"gitlab.cs.fau.de/kissen/fed/marshal"
	"log"
	"net/url"
)

// Implements the go-fed/activity/pub/FedTransport interface (version 1.0)
type FedTransport struct {
	Context   context.Context
	Target    *url.URL // TODO: use for authentication
	UserAgent string
}

// Dereference fetches the ActivityStreams object located at this IRI
// with a GET request.
func (f *FedTransport) Dereference(c context.Context, iri *url.URL) ([]byte, error) {
	log.Printf("Dereference(%v)", iri)

	if bytes, err := f.dereferenceFromStorage(c, iri); err == nil {
		return bytes, nil
	} else if bytes, err := fetch.Get(iri); err == nil {
		return bytes, nil
	} else {
		return nil, errors.Wrap(err, "cannot dereference")
	}
}

// Deliver sends an ActivityStreams object.
func (f *FedTransport) Deliver(c context.Context, b []byte, to *url.URL) (err error) {
	log.Printf("Deliver(%v)", to)
	return fetch.Post(b, to)
}

// BatchDeliver sends an ActivityStreams object to multiple recipients.
func (f *FedTransport) BatchDeliver(c context.Context, b []byte, recipients []*url.URL) error {
	log.Printf("BatchDeliver(%v)", recipients)

	// XXX: slow and wrong (quits halfway through on errors)

	for _, recipent := range recipients {
		if err := f.Deliver(c, b, recipent); err != nil {
			return err
		}
	}

	return nil
}

func (f *FedTransport) dereferenceFromStorage(c context.Context, iri *url.URL) ([]byte, error) {
	if obj, err := fedcontext.From(c).Storage.RetrieveObject(iri); err != nil {
		return nil, err
	} else if bytes, err := marshal.VocabToBytes(obj); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}
