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

// Given iterable, which is hopefully actually iterable (as defined by
// Begin), return all objects in the underlying collection.
func FetchAll(iterable interface{}) (vs []vocab.Type, err error) {
	if it, err := Begin(iterable); err != nil {
		return nil, err
	} else {
		return fetchIter(it)
	}
}

func fetchIterEntry(it IterEntry) (vocab.Type, error) {
	if !it.HasAny() {
		return nil, errors.New("no value present")
	} else if it.IsIRI() {
		return Fetch(it.GetIRI())
	} else {
		return it.GetType(), nil
	}
}

func fetchIter(it Iter) (vs []vocab.Type, err error) {
	for ; it != it.End(); it = it.Next() {
		if v, err := fetchIterEntry(it); err != nil {
			return nil, err
		} else {
			vs = append(vs, v)
		}
	}

	return vs, err
}
