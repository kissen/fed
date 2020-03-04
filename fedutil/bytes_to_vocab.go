package fedutil

import (
	"context"
	"encoding/json"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
)

func BytesToVocab(bin []byte) (vocab.Type, error) {
	// convert from []byte -> map

	var mappings map[string]interface{}

	if err := json.Unmarshal(bin, &mappings); err != nil {
		return nil, errors.Wrap(err, "byte unmarshal from object failed")
	}

	// have to set the context; see
	//
	//   https://go-fed.org/tutorial#ActivityStreams-Serialization
	//
	// for details

	mappings["@context"] = "https://www.w3.org/ns/activitystreams"

	// convert from map -> vocab.Type; first create the resolver

	var obj vocab.Type

	resolver, err := streams.NewJSONResolver(
		func(c context.Context, create vocab.ActivityStreamsCreate) error {
			obj = create
			return nil
		},

		func(c context.Context, note vocab.ActivityStreamsNote) error {
			obj = note
			return nil
		},

		func(c context.Context, note vocab.ActivityStreamsPerson) error {
			obj = note
			return nil
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not create type resolver")
	}

	// do actual resolving; Resolve() should populate obj

	if err := resolver.Resolve(context.Background(), mappings); err != nil {
		if streams.IsUnmatchedErr(err) {
			return nil, errors.Wrap(err, "missing callback, specific type not supported")
		} else {
			return nil, err
		}
	}

	return obj, nil
}
