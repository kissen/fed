package fedutil

import (
	"encoding/json"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
)

func VocabToBytes(obj vocab.Type) ([]byte, error) {
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
