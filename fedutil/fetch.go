package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/complcache"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"time"
)

const _EXPIRATION = 2 * time.Second
const _FILL = 10 * time.Second
const _GC = 1 * time.Minute

var fetchCache complcache.Cache

func init() {
	var err error

	if fetchCache, err = complcache.New(_EXPIRATION, _FILL, _GC); err != nil {
		log.Panic("cannot create http cache:", err)
	}
}

// Fetch the resource at iri. If a cached version of the resource at iri
// is available, that one is returned instead.
func Fetch(iri *url.URL) (vocab.Type, error) {
	creator := func() (interface{}, error) {
		if raw, err := Get(iri); err != nil {
			return nil, err
		} else if obj, err := BytesToVocab(raw); err != nil {
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

// Fetch the resource at addr. If a cached version of the resource at addr
// is available, that one is returned instead.
func FetchString(addr string) (vocab.Type, error) {
	if iri, err := url.Parse(addr); err != nil {
		return nil, errors.Wrap(err, "bad address")
	} else {
		return Fetch(iri)
	}
}

// Fetch the items part of collection. If an item is already a vocab.Type,
// simply return that item as-is as part of the returned slice. If it is
// an IRI, go out to the network to retrieve it first. If a cached version
// of a remote object is available, that one is returned instead as part
// of the slice.
func FetchOrGet(collection vocab.ActivityStreamsOrderedCollectionPage) (items []vocab.Type, err error) {
	iprop := collection.GetActivityStreamsOrderedItems()
	if iprop == nil {
		return nil, errors.New("items property is nil")
	}

	for it := iprop.Begin(); it != iprop.End(); it = it.Next() {
		if item, err := FetchIter(it); err != nil {
			return nil, err
		} else {
			items = append(items, item)
		}
	}

	return items, nil
}

func FetchIter(it IterEntry) (vocab.Type, error) {
	if !it.HasAny() {
		return nil, errors.New("no value present")
	} else if it.IsIRI() {
		return Fetch(it.GetIRI())
	} else {
		return it.GetType(), nil
	}
}
