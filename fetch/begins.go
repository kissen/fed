package fetch

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/kissen/fed/errors"
	"github.com/kissen/fed/prop"
	"log"
	"net/url"
	"time"
)

// Implements the iter.Iter and iter.IterEntry interfaces.
// This struct combines multiple iterators into one.
type begins struct {
	iterators []Iter
	values    []vocab.Type
}

// Singleton of begins that is returned on cals to begins.End.
var _END Iter = &begins{}

// Given iterables, return an iterator that combines all of
// these iterators and returns the individual items in chronological
// order.
//
// To sort the items in chronological order, IterEntrys that contain
// IRIs are dereferenced. We cannot know the published date of just
// an IRI after all.
func Begins(iterables ...interface{}) (Iter, error) {
	// convert iterables into Iter instances

	var its []Iter

	for i, iterable := range iterables {
		it, err := Begin(iterable)
		if err != nil {
			return nil, errors.Wrapf(err, "iterables[%v] bad", i)
		}

		its = append(its, it)
	}

	// retrieve head for all its; this checks whether they are any
	// good

	bs := &begins{}

	for i, it := range its {
		// forward to the end or the first iterator that actually
		// has a value
		for it != it.End() && !it.HasAny() {
			it = it.Next()
		}

		// skip empty its
		if it == it.End() {
			continue
		}

		// fetch value for head
		v, err := FetchIter(it)
		if err != nil {
			return nil, errors.Wrapf(err, "its[%v] bad", i)
		}

		// add to returned struct
		bs.iterators = append(bs.iterators, it)
		bs.values = append(bs.values, v)
	}

	if len(bs.iterators) == 0 {
		return _END, nil
	}

	return bs, nil
}

func (b *begins) HasAny() bool {
	return b.head() != -1
}

func (b *begins) IsIRI() bool {
	// we never have IRIs because we don't know the
	// published date for IRIs
	return false
}

func (b *begins) GetIRI() *url.URL {
	// IsIRI is always false -> GetIRI never returns
	// an actually usable value
	return nil
}

func (b *begins) GetType() vocab.Type {
	if i := b.head(); i == -1 {
		return nil
	} else {
		return b.values[i]
	}
}

func (b *begins) Next() Iter {
	// check if we are at the end

	if len(b.iterators) == 0 {
		return _END
	}

	// get the index of the current head; we'll remove
	// that head in the returned iterator

	i := b.head()
	if i == -1 {
		log.Println("this should not happen")
		return _END
	}

	// copy the lists; update the current head

	its := append([]Iter(nil), b.iterators...)
	vs := append([]vocab.Type(nil), b.values...)
	its[i] = its[i].Next()

	// keep moving to the next until we are either (1) at the
	// end or (2) we have an iterator with a value

	for its[i] != its[i].End() && !its[i].HasAny() {
		its[i] = its[i].Next()
	}

	// if its[i] is at its end, remove it; otherwise update
	// the value

	removeI := its[i] == its[i].End()

	if !removeI {
		if v, err := FetchIter(its[i]); err != nil {
			log.Println(err)
			removeI = true
		} else {
			vs[i] = v
		}
	}

	if removeI {
		lastIdx := len(its) - 1

		its[i] = its[lastIdx]
		its = its[:lastIdx]

		vs[i] = vs[lastIdx]
		vs = vs[:lastIdx]
	}

	// to make sure that checking for == End works, return the
	// End singleton if this iterator is now empty

	if len(its) == 0 {
		return _END
	}

	// we have at least one active iterator; how nice

	return &begins{
		iterators: its,
		values:    vs,
	}
}

func (b *begins) End() Iter {
	return _END
}

// Return the index in the iterators/values slice that is
// the current head. Returns -1 if these slices are empty.
func (b *begins) head() int {
	// empty -> no slice

	if len(b.iterators) == 0 {
		return -1
	}

	// find newest

	var newestIdx int
	var newestTime time.Time

	for i, v := range b.values {
		if iTime, err := prop.Published(v); err == nil {
			if iTime.After(newestTime) {
				newestIdx = i
			}
		}
	}

	return newestIdx
}
