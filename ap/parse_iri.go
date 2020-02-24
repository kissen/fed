package ap

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/db"
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Helper used by other parse* functions.
func parsePayloadFromIri(c context.Context, subpath string, iri *url.URL) (string, error) {
	// base path split into segements

	base := path.Join(*FromContext(c).BasePath, subpath)
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
func parseActorFromIri(c context.Context, iri *url.URL) (string, error) {
	return parsePayloadFromIri(c, "actor", iri)
}

// Return the unique owner of the outbox this IRI is pointing to.
// Outbox IRIs have the form
//
//   */outbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func parseOutboxOwnerFromIri(c context.Context, iri *url.URL) (string, error) {
	return parsePayloadFromIri(c, "outbox", iri)
}

// Return the unique owner of the inbox this IRI is pointing to.
// Inbox IRIs have the form
//
//   */inbox/{Owner}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func parseInboxOwnerFromIri(c context.Context, iri *url.URL) (string, error) {
	s, e := parsePayloadFromIri(c, "inbox", iri)
	return s, e
}

type ParseFunc func(c context.Context, iri *url.URL) (string, error)

// Given a parse function (e.g. parseInboxOwnerFromIri) return the user from iri.
func parseUserFrom(c context.Context, parse ParseFunc, iri *url.URL) (*db.FedUser, error) {
	var user *db.FedUser

	if username, err := parse(c, iri); err != nil {
		return nil, errors.Wrapf(err, "could not determine owner of=%v", iri)
	} else if user, err = FromContext(c).Storage.RetrieveUser(username); err != nil {
		return nil, errors.Wrapf(err, "no user found for username=%v", username)
	} else {
		return user, nil
	}
}
