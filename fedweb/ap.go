package main

import (
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"net/url"
)

func Fetch(addr string) (vocab.Type, error) {
	iri, err := url.Parse(addr)
	if err != nil {
		return nil, errors.Wrap(err, "bad address")
	}

	body, err := fedutil.Get(iri)
	if err != nil {
		return nil, errors.Wrap(err, "remote server error")
	}

	obj, err := fedutil.BytesToVocab(body)
	if err != nil {
		return nil, errors.Wrap(err, "malformed response")
	}

	return obj, nil
}
