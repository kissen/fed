package help

import "github.com/go-fed/activity/streams/vocab"

func Name(object vocab.Type) string {
	if object == nil {
		panic("argument is nil")
	} else if mappings, err := object.Serialize(); err != nil {
		panic("serialize failed")
	} else if value, ok := mappings["name"]; !ok {
		panic("missing name property")
	} else if str, ok := value.(string); !ok {
		panic("name property has wrong type")
	} else {
		return str
	}
}
