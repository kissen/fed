package prop

import (
	"net/url"
	"testing"
)

func toUrl(t *testing.T, iri string) *url.URL {
	u, err := url.Parse(iri)
	if err != nil {
		t.Fatalf("iri=%v is not parsable to url.URL", iri)
	}
	return u
}
