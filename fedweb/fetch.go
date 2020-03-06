package main

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"net/url"
	"sync"
)

var cache struct {
	sync.Mutex
	instance fedutil.VocabCache
}

func getCache() fedutil.VocabCache {
	cache.Lock()
	defer cache.Unlock()

	if cache.instance == nil {
		cache.instance = fedutil.NewVocabCache()
	}

	return cache.instance
}

func FetchIRI(iri *url.URL) (vocab.Type, error) {
	return getCache().Get(iri)
}

func Fetch(addr string) (vocab.Type, error) {
	iri, err := url.Parse(addr)
	if err != nil {
		return nil, errors.Wrap(err, "bad address")
	}

	return FetchIRI(iri)
}
