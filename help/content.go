package help

import "github.com/go-fed/activity/streams/vocab"

func Content(object vocab.Type) string {
	if mappings, err := object.Serialize(); err != nil {
		panic("serialize failed")
	} else if value, ok := mappings["content"]; !ok {
		panic("missing content property")
	} else if str, ok := value.(string); !ok {
		panic("content property has wrong type")
	} else {
		return str
	}
}
