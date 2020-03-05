package fedutil

import "github.com/go-fed/activity/streams/vocab"

func VocabToMap(obj vocab.Type) (map[string]interface{}, error) {
	return obj.Serialize()
}
