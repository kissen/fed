package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"net/url"
)

var fetchCache VocabCache

func init() {
	fetchCache = NewVocabCache()
}

func Fetch(iri *url.URL) (vocab.Type, error) {
	return fetchCache.Get(iri)
}

func FetchString(addr string) (vocab.Type, error) {
	if iri, err := url.Parse(addr); err != nil {
		return nil, errors.Wrap(err, "bad address")
	} else {
		return Fetch(iri)
	}
}
