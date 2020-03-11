package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

type Iterator interface {
	// Returns wheter any value is set.
	HasAny() bool

	// Returns whether the underlying value is an IRI and not
	// a full object.
	IsIRI() bool

	// Return the underlying IRI. Only call this if IsIRI
	// returns true.
	GetIRI() *url.URL

	// Return the underlying object. Only call this if IsIRI
	// returns false and HasAny returns true.
	GetType() vocab.Type
}
