package ap

import (
	"context"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"net/url"
	"testing"
)

func testSetUpContext() context.Context {
	ctx := context.Background()
	ctx = WithFedContext(ctx)

	FromContext(ctx).Scheme = Just("http")
	FromContext(ctx).Host = Just("example.com")
	FromContext(ctx).BasePath = Just("/fed/")

	// empty storage; feel free to overwrite with an actual
	// implementation
	FromContext(ctx).Storage = &db.FedEmptyStorage{}

	return ctx
}

func toUrl(iri string, t *testing.T) *url.URL {
	u, err := url.Parse(iri)
	if err != nil {
		t.Fatalf("iri=%v is not parsable to url.URL", iri)
	}
	return u
}
