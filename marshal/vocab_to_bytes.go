package marshal

import (
	"encoding/json"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/errors"
)

// Given an activity streams object, return a marshaled binary
// representation of that object that may be sent over the wire.
func VocabToBytes(obj vocab.Type) ([]byte, error) {
	// convert from vocab.Type -> map

	mappings, err := obj.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "serialize from object failed")
	}

	// have to set the context; see
	//
	//   https://go-fed.org/tutorial#ActivityStreams-Serialization
	//
	// for details

	mappings["@context"] = "https://www.w3.org/ns/activitystreams"

	// convert from map -> []byte

	if bytes, err := json.Marshal(mappings); err != nil {
		return nil, errors.Wrap(err, "byte marshal from object failed")
	} else {
		return bytes, nil
	}
}
