package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/complcache"
	"github.com/pkg/errors"
	"net/url"
	"log"
	"time"
)

const _EXPIRATION = 30 * time.Second
const _FILL = 15 * time.Second
const _GC = 4 * time.Minute

var fetchCache complcache.Cache

func init() {
	var err error

	if fetchCache,err = complcache.New(_EXPIRATION, _FILL, _GC); err != nil {
		log.Panic("cannot create http cache:", err)
	}
}

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

func FetchString(addr string) (vocab.Type, error) {
	if iri, err := url.Parse(addr); err != nil {
		return nil, errors.Wrap(err, "bad address")
	} else {
		return Fetch(iri)
	}
}
