package db

import (
	"net/url"
	"testing"
)

func toUrl(t *testing.T, iri string) *url.URL {
	var u *url.URL
	var err error

	if u, err = url.Parse(iri); err != nil {
		t.Fatalf("iri=%v is not parsable to url.URL", iri)
	}

	return u
}
