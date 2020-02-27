package ap

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Implements the go-fed/activity/pub/FedTransport interface (version 1.0)
type FedTransport struct {
	Context   context.Context
	Target    *url.URL // TODO: use for authentication
	UserAgent string   // TODO: use in http requests
}

// Dereference fetches the ActivityStreams object located at this IRI
// with a GET request.
func (f *FedTransport) Dereference(c context.Context, iri *url.URL) ([]byte, error) {
	log.Printf("Dereference(%v)\n", iri)

	if bytes, err := f.dereferenceFromStorage(c, iri); err == nil {
		return bytes, nil
	} else if bytes, err := f.dereferenceFromRemote(c, iri); err == nil {
		return bytes, nil
	} else {
		return nil, fmt.Errorf("cannot dereference iri=%v", iri)
	}
}

// Deliver sends an ActivityStreams object.
func (f *FedTransport) Deliver(c context.Context, b []byte, to *url.URL) error {
	log.Printf("Deliver(%v)\n", to)

	// create a copy of the raw bytes as they will get consumed by us

	copy := f.copyOf(b)

	// POST to the address

	resp, err := http.Post(to.String(), "application/ld+json", bytes.NewBuffer(copy))

	if err != nil {
		return errors.Wrapf(err, "POST to=%v failed", to)
	}

	defer resp.Body.Close()

	// evaluate result

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server to=%v rejected with status=%v", to, resp.Status)
	}

	return nil
}

// BatchDeliver sends an ActivityStreams object to multiple recipients.
func (f *FedTransport) BatchDeliver(c context.Context, b []byte, recipients []*url.URL) error {
	log.Printf("BatchDeliver(%v)\n", recipients)

	// XXX: slow and wrong (quits halfway through)

	for _, recipent := range recipients {
		if err := f.Deliver(c, b, recipent); err != nil {
			return err
		}
	}

	return nil
}

func (f *FedTransport) dereferenceFromStorage(c context.Context, iri *url.URL) ([]byte, error) {
	if obj, err := FromContext(c).Storage.RetrieveObject(iri); err != nil {
		return nil, err
	} else if bytes, err := db.VocabToBytes(obj); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func (f *FedTransport) dereferenceFromRemote(c context.Context, iri *url.URL) ([]byte, error) {
	// GET the object

	resp, err := http.Get(iri.String())

	if err != nil {
		return nil, errors.Wrapf(err, "GET to iri=%v failed", iri)
	}

	defer resp.Body.Close()

	// evaluate result

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server at iri=%v rejected with status=%v", iri, resp.Status)
	}

	// read into buffer

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)

	// put into byte slice

	return buf.Bytes(), nil
}

func (f *FedTransport) copyOf(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
