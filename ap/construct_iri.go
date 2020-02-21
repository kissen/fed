package ap

import (
	"context"
	"net/url"
	"path"
)

// Helper used by other construct*Iri functions.
func constructPrefixedIri(ctx context.Context, prefix, payload string) *url.URL {
	if isEmpty(prefix) {
		panic("cannot construct IRI with empty prefix string")
	}

	if isEmpty(payload) {
		panic("cannot construct IRI with empty payload string")
	}

	var iri url.URL

	iri.Scheme = *FromContext(ctx).Scheme
	iri.Host = *FromContext(ctx).Host
	iri.Path = path.Join(*FromContext(ctx).BasePath, prefix, payload)

	return &iri
}

// Return an IRI pointing to an actor. Actor IRIs have the form
//
//   */actor/{Actor}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func constructActorIri(ctx context.Context, actor string) *url.URL {
	return constructPrefixedIri(ctx, "actor", actor)
}

// Return an IRI pointing to an outbox. Outbox IRIs have the form
//
//   */outbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func constructOutboxIri(ctx context.Context, owner string) *url.URL {
	return constructPrefixedIri(ctx, "outbox", owner)
}

// Return an IRI pointing to an inbox. Inbox IRIs have the form
//
//   */inbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func constructInboxIri(ctx context.Context, owner string) *url.URL {
	return constructPrefixedIri(ctx, "inbox", owner)
}
