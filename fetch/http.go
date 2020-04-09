package fetch

import (
	"bytes"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"gitlab.cs.fau.de/kissen/fed/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const _CONTENT_TYPE = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`
const _HTTP_TIMEOUT = 8 * time.Second

// this file keeps some global state; in particular it contains an
// http client for multiple use
var cache struct {
	sync.Mutex
	Client *http.Client
}

// Issue an HTTP request that GETs the ActivityPub resource at iri.
func Get(iri *url.URL) (body []byte, err error) {
	log.Printf("Get(%v)", iri)

	// build up the request

	var req *http.Request

	if req, err = http.NewRequest("GET", iri.String(), nil); err != nil {
		return nil, errors.Wrap(err, "cannot set up request")
	}

	setActivityPubHeaders(req)

	// GET to the address

	var resp *http.Response

	if resp, err = client().Do(req); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// evaluate result

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v returned status=%v", iri, resp.Status)
	}

	// return the body

	return ioutil.ReadAll(resp.Body)
}

// Issue an HTTP request that POSTs body to the ActiviyPub endpoint at
// iri.
func Post(body []byte, iri *url.URL) (err error) {
	log.Printf("Post(%v)", iri)

	// preapre the io.Reader that contains the request body

	copy := append([]byte(nil), body...)
	respBody := bytes.NewReader(copy)

	// build up the request

	var req *http.Request

	if req, err = http.NewRequest("POST", iri.String(), respBody); err != nil {
		return errors.Wrap(err, "cannot set up request")
	}

	setActivityPubHeaders(req)

	// POST to the address

	var resp *http.Response

	if resp, err = client().Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	// evaluate result

	if !util.IsHTTPSuccess(resp.StatusCode) {
		// read error body if there is one

		body := ""

		if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			body = string(bodyBytes)
		}

		// return error

		return fmt.Errorf(`%v returned status="%v" body="%v"`, iri, resp.Status, body)
	}

	return nil
}

// Return an HTTP client that may be used to interact with ActivityPub
// servers.
func client() *http.Client {
	cache.Lock()
	defer cache.Unlock()

	// only create one client per transport and re-use them; see
	// https://golang.org/pkg/net/http/#Client for details

	if cache.Client == nil {
		cache.Client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				setActivityPubHeaders(req)
				return nil
			},
			Timeout: _HTTP_TIMEOUT,
		}
	}

	return cache.Client
}

// Set Content-Type/Accept and User-Agent headers on req such that it
// may interact with ActivityPub services.
func setActivityPubHeaders(req *http.Request) {
	switch req.Method {
	case "GET":
		req.Header.Set("Accept", _CONTENT_TYPE)
	case "POST":
		req.Header.Set("Content-Type", _CONTENT_TYPE)
	}

	req.Header.Set("User-Agent", "fed/0.x")
}
