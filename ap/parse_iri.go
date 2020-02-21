package ap

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Helper used by other parse* functions.
func parsePayloadFromIri(ctx context.Context, subpath string, iri *url.URL) (string, error) {
	// base path split into segements

	base := path.Join(*FromContext(ctx).BasePath, subpath)
	baseComps := strings.Split(base, "/")

	// iri path split into segments

	iriComps := strings.Split(iri.Path, "/")

	// match number

	if len(baseComps)+1 != len(iriComps) {
		return "", errors.New("bad number of path segements")
	}

	// match strings

	for i := range baseComps {
		expected := baseComps[i]
		actual := iriComps[i]

		if expected != actual {
			return "", fmt.Errorf("bad pattern, expected=%v actual=%v", expected, actual)
		}
	}

	// match last component

	if payload := iriComps[len(iriComps)-1]; isEmpty(payload) {
		return "", errors.New("last path component empty")
	} else {
		return payload, nil
	}
}

// Return the unique actor this IRI is pointing to. Actor IRIs have the form
//
//   */actor/{Actor}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func parseActorFromIri(ctx context.Context, iri *url.URL) (string, error) {
	return parsePayloadFromIri(ctx, "actor", iri)
}

// Return the unique owner of the outbox this IRI is pointing to.
// Outbox IRIs have the form
//
//   */outbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func parseOutboxOwnerFromIri(ctx context.Context, iri *url.URL) (string, error) {
	return parsePayloadFromIri(ctx, "outbox", iri)
}

// Return the unique owner of the inbox this IRI is pointing to.
// Inbox IRIs have the form
//
//   */inbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func parseInboxOwnerFromIri(ctx context.Context, iri *url.URL) (string, error) {
	s, e := parsePayloadFromIri(ctx, "inbox", iri)
	return s, e
}
