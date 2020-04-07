package prop

import (
	"net/url"
)

type Appender interface {
	AppendIRI(*url.URL)
}

// Append all elements in iris to target.
func AppendIRIs(target Appender, iris []*url.URL) {
	for _, iri := range iris {
		target.AppendIRI(iri)
	}
}
