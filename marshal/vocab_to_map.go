package marshal

import "github.com/go-fed/activity/streams/vocab"

// Given an activity streams object, return a key/value representation
// of that object.
func VocabToMap(obj vocab.Type) (map[string]interface{}, error) {
	// really just a wrapper for consistent naming in the codebase
	return obj.Serialize()
}
