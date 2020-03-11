package fedutil

import (
	"github.com/go-fed/activity/streams/vocab"
	"log"
)

func Type(obj vocab.Type) string {
	if mappings, err := obj.Serialize(); err != nil {
		log.Println("cannot serialize:", err)
		return "(unknown)"
	} else if typename, ok := mappings["type"]; !ok {
		log.Println("missing type property")
		return "(unknown)"
	} else if s, ok := typename.(string); !ok {
		log.Println("bad type for type property")
		return "(unknown)"
	} else {
		return s
	}
}
