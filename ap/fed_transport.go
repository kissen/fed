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
	"sync"
	"time"
)

// The content type to use for PUT (Content Type) and GET (Accept)
const _CONTENT_TYPE = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

// Timeout for our HTTP requests. I guess eventually we will implement
// retries and the like; at that point the timeouts will become more complex
const _HTTP_TIMEOUT = 8 * time.Second

// Implements the go-fed/activity/pub/FedTransport interface (version 1.0)
type FedTransport struct {
	Context   context.Context
	Target    *url.URL // TODO: use for authentication
	UserAgent string   // TODO: use in http requests

	// http client used for http requests; use client() to get a usable
	// client
	cachedClient     *http.Client
	cachedClientLock sync.Mutex
}

// Dereference fetches the ActivityStreams object located at this IRI
// with a GET request.
func (f *FedTransport) Dereference(c context.Context, iri *url.URL) ([]byte, error) {
	log.Printf("Dereference(%v)", iri)

	if bytes, err := f.dereferenceFromStorage(c, iri); err == nil {
		return bytes, nil
	} else if bytes, err := f.dereferenceFromRemote(c, iri); err == nil {
		return bytes, nil
	} else {
		return nil, errors.Wrap(err, "cannot dereference")
	}
}

// Deliver sends an ActivityStreams object.
func (f *FedTransport) Deliver(c context.Context, b []byte, to *url.URL) (err error) {
	log.Printf("Deliver(%v)", to)

	// preapre the io.Reader that contains the request body

	copy := f.copyOf(b)
	body := bytes.NewReader(copy)

	// build up the request

	var req *http.Request

	if req, err = http.NewRequest("PUT", to.String(), body); err != nil {
		return errors.Wrap(err, "cannot set up request")
	}

	req.Header.Set("Content-Type", _CONTENT_TYPE)

	// POST to the address

	var resp *http.Response

	if resp, err = f.client().Do(req); err != nil {
		return errors.Wrap(err, "cannot do request")
	}

	defer resp.Body.Close()

	// evaluate result

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("to=%v replied status=%v", to, resp.Status)
	}

	return nil
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
		return nil, err
	}

	defer resp.Body.Close()

	// evaluate result

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iri=%v replied status=%v", iri, resp.Status)
	}

	// read into buffer

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)

	// put into byte slice

	return buf.Bytes(), nil
}

// Return a copy of src.
func (f *FedTransport) copyOf(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// Return an HTTP client usable for requests. It comes with a timeout so we don't
// get stuck waiting forever.
func (f *FedTransport) client() *http.Client {
	f.cachedClientLock.Lock()
	defer f.cachedClientLock.Unlock()

	// only create one client per transport and re-use them; see
	// https://golang.org/pkg/net/http/#Client for details

	if f.cachedClient == nil {
		f.cachedClient = &http.Client{
			Timeout: _HTTP_TIMEOUT,
		}
	}

	return f.cachedClient
}
