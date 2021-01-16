package fetch

import (
	"errors"
	"github.com/kissen/fed/prop"
	"net/url"
)

// Given ie, return the IRI to that entry.
//
// This will not involve any network request, the
// information is sourced from ie alone.
func IRI(ie IterEntry) (iri *url.URL, err error) {
	if !ie.HasAny() {
		return nil, errors.New("ie is empty")
	}

	// try out iri

	if ie.IsIRI() {
		return ie.GetIRI(), nil
	}

	// try out object

	obj := ie.GetType()
	return prop.Id(obj), nil
}

// Iterate over it and return IRIs to each element
// in that iterator.
//
// This will not involve any network request, the
// information is sourced from ie alone.
func IRIs(it Iter) (iris []*url.URL, err error) {
	for ; it != it.End(); it = it.Next() {
		iri, err := IRI(it)
		if err != nil {
			return nil, err
		}

		iris = append(iris, iri)
	}

	return iris, nil
}
