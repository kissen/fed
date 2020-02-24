package help

import (
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

func Id(object vocab.Type) *url.URL {
	property := object.GetJSONLDId()

	if property.IsIRI() {
		return property.GetIRI()
	} else if property.IsXMLSchemaAnyURI() {
		return property.Get()
	} else {
		panic("property is neither IRI nor XMLSchemaAnyURI")
	}
}
