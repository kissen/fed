package ap

import (
	"context"
	"net/url"
	"testing"
)

func testSetUpContext() context.Context {
	ctx := context.Background()
	ctx = WithFedContext(ctx)

	FromContext(ctx).Scheme = Just("http")
	FromContext(ctx).Host = Just("example.com")
	FromContext(ctx).BasePath = Just("/fed/")

	return ctx
}

func toUrl(iri string, t *testing.T) *url.URL {
	u, err := url.Parse(iri)
	if err != nil {
		t.Fatalf("iri=%v is not parsable to url.URL", iri)
	}
	return u
}
