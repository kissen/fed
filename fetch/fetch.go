package fetch

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/complcache"
	"github.com/kissen/fed/errors"
	"github.com/kissen/fed/marshal"
	"log"
	"net/url"
	"time"
)

const _EXPIRATION = 2 * time.Second
const _FILL = 10 * time.Second
const _GC = 1 * time.Minute

// Contains cached versions of vocab.Type objects identified by their
// IRI.
var fetchCache complcache.Cache

// Create the cache on init.
func init() {
	var err error

	if fetchCache, err = complcache.New(_EXPIRATION, _FILL, _GC); err != nil {
		log.Panicf("cannot create cache: %v", err)
	}
}

// Fetch the resource at iri.
//
// If a cached version of the resource at iri is available, that
// one is returned. If no cached version is available, creates an
// HTTP GET request that gets the resource.
func Fetch(iri *url.URL) (vocab.Type, error) {
	creator := func() (interface{}, error) {
		if raw, err := Get(iri); err != nil {
			return nil, err
		} else if obj, err := marshal.BytesToVocab(raw); err != nil {
			return nil, err
		} else {
			return obj, nil
		}
	}

	if obj, err := fetchCache.GetOrCreate(iri.String(), creator); err != nil {
		return nil, err
	} else {
		return obj.(vocab.Type), nil
	}
}

// Fetch the resource at it.
//
// If a cached version of the resource at iri is available, that
// one is returned. If no cached version is available, creates an
// HTTP GET request that gets the resource.
func FetchIter(it IterEntry) (vocab.Type, error) {
	if !it.HasAny() {
		return nil, errors.New("no value present")
	} else if it.IsIRI() {
		return Fetch(it.GetIRI())
	} else {
		return it.GetType(), nil
	}
}

// Fetch all objects from iterator it.
//
// If cached versions of the resources are available, these ones
// are returned. For each object for which no cached version is
// available, creates an HTTP GET request that gets the resource.
func FetchIters(it Iter) (vs []vocab.Type, err error) {
	for ; it != it.End(); it = it.Next() {
		if v, err := FetchIter(it); err != nil {
			return nil, err
		} else {
			vs = append(vs, v)
		}
	}

	return vs, err
}
