package db

import (
	"context"
	"encoding/json"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
)

func vocabToBytes(obj vocab.Type) ([]byte, error) {
	// convert from vocab.Type -> map

	mappings, err := obj.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "serialize from object failed")
	}

	// convert from map -> []byte

	if bytes, err := json.Marshal(mappings); err != nil {
		return nil, errors.Wrap(err, "byte marshal from object failed")
	} else {
		return bytes, nil
	}
}

func bytesToVocab(bin []byte) (vocab.Type, error) {
	// convert from []byte -> map

	var mappings map[string]interface{}

	if err := json.Unmarshal(bin, &mappings); err != nil {
		return nil, errors.Wrap(err, "byte unmarshal from object failed")
	}

	// convert from map -> vocab.Type

	var obj vocab.Type

	resolver, err := streams.NewJSONResolver(
		func(c context.Context, create vocab.ActivityStreamsCreate) error {
			obj = create
			return nil
		},

		func(c context.Context, create vocab.ActivityStreamsNote) error {
			obj = create
			return nil
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not create type resolver")
	}

	resolver.Resolve(context.Background(), mappings) // populates obj
	return obj, nil
}
