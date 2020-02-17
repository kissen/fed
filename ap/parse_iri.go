package ap

import (
	"context"
	"errors"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/db"
	"net/url"
	"path"
	"strconv"
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

	payload := iriComps[len(iriComps)-1]

	if payload == "" {
		return "", errors.New("last path component empty")
	}

	return payload, nil
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
	return parsePayloadFromIri(ctx, "inbox", iri)
}

// Return the unique id of the activity this IRI is pointing to.
// Activity IRIs have the form
//
//   */activity/{ActivityId}
//
// where the asterix is a placeholder for protocol, hostname and
// base path. An ActivityId is a non-negative 64 bit integer
// in base 10 representation.
func parseActivityIdFromIri(ctx context.Context, iri *url.URL) (db.FedId, error) {
	if activityIdStr, err := parsePayloadFromIri(ctx, "activity", iri); err == nil {
		id, err := strconv.ParseUint(activityIdStr, 10, 64)
		return db.FedId(id), err
	} else {
		return 0, err
	}
}
