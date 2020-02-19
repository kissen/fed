package ap

import (
	"context"
	"log"
	"errors"
	"net/url"
)

// Implements the go-fed/activity/pub/FedTransport interface (version 1.0)
type FedTransport struct {
	Context   context.Context
	Target    *url.URL
	UserAgent string
}

// Dereference fetches the ActivityStreams object located at this IRI
// with a GET request.
func (f *FedTransport) Dereference(c context.Context, iri *url.URL) ([]byte, error) {
	log.Printf("Dereference(%v)\n", iri)

	return nil, errors.New("not implemented")
}

// Deliver sends an ActivityStreams object.
func (f *FedTransport) Deliver(c context.Context, b []byte, to *url.URL) error {
	log.Printf("Deliver(%v)\n", to)

	return errors.New("not implemented")
}

// BatchDeliver sends an ActivityStreams object to multiple recipients.
func (f *FedTransport) BatchDeliver(c context.Context, b []byte, recipients []*url.URL) error {
	log.Printf("BatchDeliver(%v)\n", recipients)

	return errors.New("not implemented")
}
