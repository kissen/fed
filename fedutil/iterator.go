package fedutil

import (
	"errors"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// A single entry of an iterator. Many of the go-fed iterator types
// implement this interface.
type IterEntry interface {
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

// A pimped version of IterEntry that also allows us to get the next
// item and determine whether we have reached the end.
//
// This is what I think an iterator should look like.
type Iter interface {
	IterEntry

	Next() Iter
	End() Iter
}

// Return a generic iterator if iterable is actually iterable or some
// error if it is no. I wish I knew of a better way short of writing
// my own lib.
func Begin(iterable interface{}) (Iter, error) {
	if iterable == nil {
		return nil, errors.New("nil argument")
	}

	switch v := iterable.(type) {
	case vocab.ActivityStreamsOrderedCollectionPage:
		items := v.GetActivityStreamsOrderedItems()
		return begin(items)

	default:
		return begin(iterable)
	}
}
