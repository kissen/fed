package ap

import (
	"context"
	"github.com/google/uuid"
	"net/url"
)

// Return a new IRI pointing to some ActicityPub object. Such IRIs
// have the form
//
//   */object/{UUID}
//
// where the asterix is a placeholder for protocol, hostname and
// base path.
func rollObjectIri(ctx context.Context) *url.URL {
	tail := uuid.New().String()
	return constructPrefixedIri(ctx, "object", tail)
}
