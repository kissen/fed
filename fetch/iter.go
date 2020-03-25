package fetch

import (
	"errors"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// A single entry of an iterator. Many of the go-fed iterator types
// implement this interface.
//
// It either contains a vocab.Type instance,
// an IRI pointing to some object that can be represented as vocab.Type
// or nothing (in which case HasAny returns false).
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
// item and determine whether we have reached the end. To iterate over
// some Iter, you would write code like this
//
//   for it := getIter(); it != it.End(); it = it.Next() {
//       ...
//   }
type Iter interface {
	IterEntry

	Next() Iter
	End() Iter
}

// Return a generic iterator over iterable.
//
// If we don't know how to iterate over iterable or it is
// in fact not possible, an error is returned.
func Begin(iterable interface{}) (Iter, error) {
	if iterable == nil {
		return nil, errors.New("nil argument")
	}

	if it, ok := iterable.(Iter); ok {
		return it, nil
	}

	switch v := iterable.(type) {
	case vocab.ActivityStreamsOrderedCollectionPage:
		items := v.GetActivityStreamsOrderedItems()
		return begin(items)

	default:
		return begin(iterable)
	}
}
