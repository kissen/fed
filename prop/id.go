package prop

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

func Id(object vocab.Type) *url.URL {
	if object == nil {
		panic("argument is nil")
	} else if property := object.GetJSONLDId(); property.IsIRI() {
		return property.GetIRI()
	} else if property.IsXMLSchemaAnyURI() {
		return property.Get()
	} else {
		panic("property is neither IRI nor XMLSchemaAnyURI")
	}
}
